package app

import (
	"encoding/gob"
	"fmt"
	"kubevue/args"
	"kubevue/auth"
	"kubevue/crd"
	"regexp"

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
	brokers BrokerMap
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

// Objects return list of known objects
func (a *App) GetObjects() (re []string) {
	for namespace, _ := range a.brokers {
		for name, _ := range a.brokers[namespace] {
			re = append(re, fmt.Sprintf("%s/%s", namespace, name))
		}
	}
	return
}

// New creates an application instance
func New() *App {
	a := new(App)
	a.args = args.New().LogLevel()
	a.auth = auth.New(a.Args().OIDC())
	a.store = sessions.NewFilesystemStore("", []byte("session-secret"))
	a.brokers = make(BrokerMap)
	gob.Register(map[string]interface{}{})
	go a.watchObjects(a.newBroker("kubevue.io", "v1", "objects", a.args.Namespace()))
	return a
}

func (a *App) watchObjects(cb *CrdBroker) {
	for msg := range cb.crd.Notifier() {
		cb.broker.Notifier <- msg
		m := crd.Parse(msg.Content)
		switch msg.Action {
		case "add":
			if broker := a.brokers[m.Name]; broker == nil {
				log.Infof("adding object %s/%s to roster", m.Namespace, m.Name)
				a.newBroker(m.Group, m.Version, m.Name, m.Namespace).PassMessages()
			} else {
				log.Infof("skip adding object %s/%s to roster", m.Namespace, m.Name)
			}
		case "delete":
			if cb := a.getBroker(m.Name, m.Namespace); cb != nil {
				log.Infof("deleting object %s/%s from roster", m.Namespace, m.Name)
				cb.Stop()
				a.deleteBroker(m.Name, m.Namespace)
			} else {
				log.Infof("skip deleting object %s/%s to roster", m.Namespace, m.Name)
			}
		case "update":
			log.Infof("skip updating object %s/%s", m.Namespace, m.Name)
		default:
		}
	}
}

var bypassAuth []*regexp.Regexp = []*regexp.Regexp{
	regexp.MustCompile("^/profile$"),
	regexp.MustCompile("^/auth$"),
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
	log.Debugf("Request:%v user:%s remote:%s", r.RequestURI, profile["name"], r.RemoteAddr)
	return profile
}

func (a *App) authMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		for _, re := range bypassAuth {
			if re.MatchString(r.RequestURI) {
				log.Debugf("Request No-Auth: %v", r.RequestURI)
				next.ServeHTTP(w, r)
				return
			}
		}
		if a.checkAuth(w, r) != nil {
			next.ServeHTTP(w, r)
		}
	})
}

// Serve ui and api endpoints
func (a *App) Serve() {
	bindAddr := fmt.Sprintf("%s:%d", a.Args().BindAddr(), a.Args().Port())
	log.Infof("Serving %s, static folder:%s", bindAddr, a.Args().Dir())
	r := mux.NewRouter()
	r.PathPrefix("/ui/").Handler(http.StripPrefix("/ui/", http.FileServer(http.Dir(a.Args().Dir()))))
	r.HandleFunc("/watch/{namespace}/{kind}", a.Watch)
	r.HandleFunc("/watch/{namespace}/{kind}/{name}", a.Watch)
	r.HandleFunc("/proxy/{namespace}/{name}", a.ProxyService)
	r.HandleFunc("/objects", a.Objects)
	r.HandleFunc("/auth", a.AuthInitiate)
	r.HandleFunc("/callback", a.AuthCallback)
	r.HandleFunc("/profile", a.Profile)
	r.HandleFunc("/logout", a.Logout)
	r.Use(a.authMiddleWare)
	srv := &http.Server{
		Handler: r,
		Addr:    bindAddr,
	}
	srv.SetKeepAlivesEnabled(true)
	log.Fatal(srv.ListenAndServe())
}
