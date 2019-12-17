package app

import (
	"encoding/gob"
	"fmt"
	"kubevue/args"
	"kubevue/auth"
	"kubevue/crd"

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
		switch msg.Action {
		case "add":
			m := crd.Parse(msg.Content)
			if broker := a.brokers[m.Name]; broker == nil {
				log.Infof("adding object %s/%s to roster", m.Namespace, m.Name)
				a.newBroker(m.Group, m.Version, m.Name, m.Namespace).PassMessages()
			} else {
				log.Infof("skip adding object %s/%s to roster", m.Namespace, m.Name)
			}
		case "delete":
			m := crd.Parse(msg.Content)
			if cb := a.getBroker(m.Name, m.Namespace); cb != nil {
				log.Infof("deleting object %s/%s from roster", m.Namespace, m.Name)
				cb.Stop()
				a.deleteBroker(m.Name, m.Namespace)
			} else {
				log.Infof("skip deleting object %s/%s to roster", m.Namespace, m.Name)
			}
		case "update":
			m := crd.Parse(msg.Content)
			log.Infof("skip updating object %s/%s", m.Namespace, m.Name)
		default:
		}
	}
}

// Serve ui and api endpoints
func (a *App) Serve() {
	bindAddr := fmt.Sprintf("%s:%d", a.Args().BindAddr(), a.Args().Port())
	log.Infof("Serving %s, static folder:%s", bindAddr, a.Args().Dir())
	r := mux.NewRouter()
	r.PathPrefix("/ui/").Handler(http.StripPrefix("/ui/", http.FileServer(http.Dir(a.Args().Dir()))))
	r.HandleFunc("/watch/{namespace}/{objects}", a.Watch)
	r.HandleFunc("/objects", a.Objects)
	r.HandleFunc("/auth", a.AuthInitiate)
	r.HandleFunc("/callback", a.AuthCallback)
	r.HandleFunc("/profile", a.Profile)
	r.HandleFunc("/logout", a.Logout)
	srv := &http.Server{
		Handler: r,
		Addr:    bindAddr,
	}
	srv.SetKeepAlivesEnabled(true)
	log.Fatal(srv.ListenAndServe())
}
