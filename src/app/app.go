package app

import (
	"argovue/args"
	"argovue/auth"
	"encoding/gob"
	"fmt"
	"regexp"
	"sync"

	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type BrokerMap map[string]map[string]*CrdBroker

// App defines application dependencies
type App struct {
	args    *args.Args
	auth    *auth.Auth
	store   *sessions.FilesystemStore
	wg      sync.WaitGroup
	brokers BrokerMap
	subset  BrokerMap
	events  chan *Event
}

// Args returns application argumnets
func (a *App) Args() *args.Args {
	return a.args
}

// Auth returns authenticator
func (a *App) Auth() *auth.Auth {
	return a.auth
}

// Store returns session store
func (a *App) Store() *sessions.FilesystemStore {
	return a.store
}

// New creates an application instance
func New() *App {
	a := new(App)
	a.args = args.New().LogLevel()
	a.store = sessions.NewFilesystemStore("", []byte("session-secret"))
	a.brokers = make(BrokerMap)
	a.subset = make(BrokerMap)
	a.events = make(chan *Event)
	gob.Register(map[string]interface{}{})
	go a.Serve()
	a.auth = auth.New(a.Args().OIDC())
	return a
}

var bypassAuth []*regexp.Regexp = []*regexp.Regexp{
	regexp.MustCompile("^/profile$"),
	regexp.MustCompile("^/auth$"),
	regexp.MustCompile("^/logout$"),
	regexp.MustCompile("^/dex/.*"),
	regexp.MustCompile("^/callback.*$"),
	regexp.MustCompile("^/ui/.*$"),
}

func (a *App) checkAuth(w http.ResponseWriter, r *http.Request) map[string]interface{} {
	session, _ := a.Store().Get(r, "auth-session")
	profileRef := session.Values["profile"]
	if profileRef == nil {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return nil
	}
	profile, ok := profileRef.(map[string]interface{})
	if !ok {
		http.Error(w, "Profile is not a map", http.StatusInternalServerError)
		return nil
	}
	log.Debugf("HTTP: '%s' from:%s %s", profile["name"], r.RemoteAddr, r.RequestURI)
	return profile
}

func (a *App) authMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		for _, re := range bypassAuth {
			if re.MatchString(r.RequestURI) {
				log.Debugf("HTTP: no-auth from:%s %v", r.RemoteAddr, r.RequestURI)
				next.ServeHTTP(w, r)
				return
			}
		}
		if a.checkAuth(w, r) != nil {
			if _, err := a.Store().Get(r, "auth-session"); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			next.ServeHTTP(w, r)
		}
	})
}

// Serve ui and api endpoints
func (a *App) Serve() {
	a.wg.Add(1)
	defer a.wg.Done()
	bindAddr := fmt.Sprintf("%s:%d", a.Args().BindAddr(), a.Args().Port())
	log.Infof("HTTP: at %s static:%s start", bindAddr, a.Args().Dir())
	r := mux.NewRouter()
	r.PathPrefix("/ui/").Handler(http.StripPrefix("/ui/", http.FileServer(http.Dir(a.Args().Dir()))))
	r.HandleFunc("/events", a.handleEvents)
	r.HandleFunc("/objects", a.Objects)
	r.HandleFunc("/auth", a.AuthInitiate)
	r.HandleFunc("/callback", a.AuthCallback)
	r.HandleFunc("/profile", a.Profile)
	r.HandleFunc("/logout", a.Logout)
	r.HandleFunc("/watch/{kind}", a.watchKind)

	r.HandleFunc("/proxy/{namespace}/{name}/{port}/{rest:.*}", a.proxyService)
	r.HandleFunc("/proxy/{namespace}/{name}/{port}", a.proxyService)
	r.HandleFunc("/dex/{rest:.*}", a.proxyDex)

	r.HandleFunc("/k8s/{kind}/{namespace}/{name}", a.watchObject)

	r.HandleFunc("/catalogue/{namespace}/{name}", a.watchCatalogue)
	r.HandleFunc("/catalogue/{namespace}/{name}/instances", a.watchCatalogueInstances)
	r.HandleFunc("/catalogue/{namespace}/{name}/instance/{instance}", a.watchCatalogueInstance)
	r.HandleFunc("/catalogue/{namespace}/{name}/{action}", a.commandCatalogue).Methods("POST")
	r.HandleFunc("/catalogue/{namespace}/{name}/instance/{instance}/action/{action}", a.controlCatalogueInstance).Methods("POST")

	r.HandleFunc("/workflow/{namespace}/{name}", a.watchWorkflow)
	r.HandleFunc("/workflow/{namespace}/{name}/services", a.watchWorkflowServices)
	r.HandleFunc("/workflow/{namespace}/{name}/service/{service}/action/{action}", a.controlWorkflowService).Methods("POST")
	r.HandleFunc("/workflow/{namespace}/{name}/pod/{pod}", a.watchWorkflowPods)
	r.HandleFunc("/workflow/{namespace}/{name}/pod/{pod}/container/{container}/logs", a.watchWorkflowPodLogs)
	r.HandleFunc("/workflow/{namespace}/{name}/action/{action}", a.commandWorkflow).Methods("POST")

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
