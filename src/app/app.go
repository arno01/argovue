package app

import (
	"argovue/args"
	"argovue/auth"
	"argovue/crd"
	"argovue/kube"
	"argovue/profile"
	"encoding/gob"
	"fmt"
	"regexp"
	"sync"

	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type BrokerMap map[string]map[string]*CrdBroker

type App struct {
	args    *args.Args
	auth    *auth.Auth
	store   *sessions.FilesystemStore
	wg      sync.WaitGroup
	brokers BrokerMap
	config  *crd.Crd
	groups  map[string]string
	subset  BrokerMap
	events  chan *Event
	ver     map[string]interface{}
}

func (a *App) Args() *args.Args {
	return a.args
}

func (a *App) Auth() *auth.Auth {
	return a.auth
}

func (a *App) Store() *sessions.FilesystemStore {
	return a.store
}

func New() *App {
	a := new(App)
	a.args = args.New().LogLevel()
	a.store = sessions.NewFilesystemStore("/tmp", []byte(a.Args().SessionKey()))
	a.brokers = make(BrokerMap)
	a.subset = make(BrokerMap)
	a.events = make(chan *Event)
	gob.Register(map[string]interface{}{})
	gob.Register(profile.Profile{})
	go a.Serve()
	a.auth = auth.New(a.Args().OIDC())
	a.config = crd.New("argovue.io", "v1", "appconfigs").
		SetFieldSelector(fmt.Sprintf("metadata.namespace=%s,metadata.name=%s-config", a.Args().Namespace(), a.Args().Release())).
		Watch()
	go a.ListenForConfig()
	a.ver = make(map[string]interface{})
	return a
}

func (a *App) SetVersion(version, commit, builddate string) *App {
	a.ver["version"] = version
	a.ver["commit"] = commit
	a.ver["builddate"] = builddate
	if client, err := kube.GetClient(); err == nil {
		if ver, err := client.ServerVersion(); err == nil {
			a.ver["kubernetes"] = ver
		}
	}
	return a
}

var bypassAuth []*regexp.Regexp = []*regexp.Regexp{
	regexp.MustCompile("^/profile$"),
	regexp.MustCompile("^/auth"),
	regexp.MustCompile("^/logout$"),
	regexp.MustCompile("^/dex/.*"),
	regexp.MustCompile("^/callback.*$"),
	regexp.MustCompile("^/ui/.*$"),
}

func makeError(code int, format string, args ...interface{}) *appError {
	return &appError{Error: fmt.Sprintf(format, args...), Code: code}
}

type appError struct {
	Error string
	Code  int
}

type appHandler func(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError
type httpHandler func(w http.ResponseWriter, r *http.Request) *appError

func (fn httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debugf("HTTP: start %s", r.RequestURI)
	if err := fn(w, r); err != nil {
		log.Debug(err.Error)
		http.Error(w, err.Error, err.Code)
	}
	log.Debugf("HTTP: stop %s", r.RequestURI)
}

func (a *App) appHandler(fn appHandler) httpHandler {
	return func(w http.ResponseWriter, r *http.Request) *appError {
		session, err := a.Store().Get(r, "auth-session")
		if err != nil {
			return makeError(http.StatusInternalServerError, "Can't get session, error:%s", err)
		}
		pf := session.Values["profile"].(profile.Profile)
		return fn(session.ID, &pf, w, r)
	}
}

func (a *App) authMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", a.Args().UIRootDomain())
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		for _, re := range bypassAuth {
			if re.MatchString(r.RequestURI) {
				log.Debugf("HTTP: no-auth from:%s %v", r.RemoteAddr, r.RequestURI)
				next.ServeHTTP(w, r)
				return
			}
		}
		session, err := a.Store().Get(r, "auth-session")
		if err != nil {
			log.Debugf("Can't get session, error:%s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		p, ok := session.Values["profile"].(profile.Profile)
		if !ok {
			log.Debugf("Not authorized, request:%s", r.RequestURI)
			http.Redirect(w, r, "/auth?redirect="+url.PathEscape(r.RequestURI), http.StatusFound)
			return
		}
		log.Debugf("HTTP: %s from:%s %s", p.Id, r.RemoteAddr, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func (a *App) Serve() {
	a.wg.Add(1)
	defer a.wg.Done()
	bindAddr := fmt.Sprintf("%s:%d", a.Args().BindAddr(), a.Args().Port())
	log.Infof("HTTP: at %s static:%s start", bindAddr, a.Args().Dir())
	r := mux.NewRouter()
	r.PathPrefix("/ui/").Handler(http.StripPrefix("/ui/", http.FileServer(http.Dir(a.Args().Dir()))))
	r.HandleFunc("/events", a.handleEvents)
	r.HandleFunc("/profile", a.Profile)
	r.HandleFunc("/objects", a.Objects)
	r.HandleFunc("/version", a.Version)
	r.HandleFunc("/auth", a.AuthInitiate)
	r.Handle("/callback", httpHandler(a.AuthCallback))
	r.HandleFunc("/logout", a.Logout)
	r.Handle("/watch/{kind}", a.appHandler(a.watchKind))

	r.Handle("/proxy/{namespace}/{name}/{port}/{rest:.*}", a.appHandler(a.proxyService))
	r.Handle("/proxy/{namespace}/{name}/{port}", a.appHandler(a.proxyService))
	r.Handle("/dex/{rest:.*}", a.appHandler(a.proxyDex))

	r.Handle("/k8s/{kind}/{namespace}/{name}", a.appHandler(a.watchObject))
	r.Handle("/k8s/pod/{namespace}/{name}/container/{container}/logs", a.appHandler(a.watchPodLogs))

	r.Handle("/catalogue/{namespace}/{name}", a.appHandler(a.watchCatalogue))
	r.Handle("/catalogue/{namespace}/{name}/instances", a.appHandler(a.watchCatalogueInstances))
	r.Handle("/catalogue/{namespace}/{name}/resources", a.appHandler(a.watchCatalogueResources))
	r.Handle("/catalogue/{namespace}/{name}/instance/{instance}", a.appHandler(a.watchCatalogueInstance))
	r.Handle("/catalogue/{namespace}/{name}/instance/{instance}/resources", a.appHandler(a.watchCatalogueInstanceResources))
	r.Handle("/catalogue/{namespace}/{name}/{action}", a.appHandler(a.controlCatalogue)).Methods("POST", "OPTIONS")
	r.Handle("/catalogue/{namespace}/{name}/instance/{instance}/action/{action}", a.appHandler(a.controlCatalogueInstance)).Methods("POST", "OPTIONS")

	r.HandleFunc("/workflow/{namespace}/{name}", a.watchWorkflow)
	r.HandleFunc("/workflow/{namespace}/{name}/services", a.watchWorkflowServices)
	r.HandleFunc("/workflow/{namespace}/{name}/mounts", a.watchWorkflowMounts)
	r.HandleFunc("/workflow/{namespace}/{name}/service/{service}/action/{action}", a.controlWorkflowService).Methods("POST", "OPTIONS")
	r.HandleFunc("/workflow/{namespace}/{name}/pod/{pod}", a.watchWorkflowPods)
	r.HandleFunc("/workflow/{namespace}/{name}/pod/{pod}/container/{container}/logs", a.watchWorkflowPodLogs)
	r.HandleFunc("/workflow/{namespace}/{name}/action/{action}", a.commandWorkflow).Methods("POST", "OPTIONS")

	r.Use(a.authMiddleWare)
	srv := &http.Server{
		Handler: r,
		Addr:    bindAddr,
	}
	srv.SetKeepAlivesEnabled(true)
	log.Fatal(srv.ListenAndServe())
}

func (a *App) Wait() {
	a.wg.Wait()
}
