package app

import (
	"encoding/gob"
	"fmt"
	"kubevue/args"
	"kubevue/auth"
	"kubevue/crd"

	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// App defines application dependencies
type App struct {
	args      *args.Args
	clientset *kubernetes.Clientset
	auth      *auth.Auth
	store     *sessions.FilesystemStore
	brokers   map[string]*CrdBroker
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

// ClientSet returns K8S Client Set
func (a *App) ClientSet() *kubernetes.Clientset {
	return a.clientset
}

// Objects return list of known objects
func (a *App) GetObjects() (re []string) {
	for name, _ := range a.brokers {
		re = append(re, name)
	}
	return
}

// New creates an application instance
func New() *App {
	a := new(App)
	a.args = args.New().LogLevel()
	a.auth = auth.New(a.Args().OIDC())
	a.store = sessions.NewFilesystemStore("", []byte("session-secret"))
	a.brokers = make(map[string]*CrdBroker)
	a.brokers["objects"] = NewBroker("kubevue.io", "v1", "objects", a.args.Namespace())
	gob.Register(map[string]interface{}{})
	go a.watchObjects()
	return a
}

func (a *App) watchObjects() {
	cb := a.brokers["objects"]
	for msg := range cb.crd.Notifier() {
		cb.broker.Notifier <- msg
		switch msg.Action {
		case "add":
			m := crd.Parse(msg.Content)
			if broker := a.brokers[m.Name]; broker == nil {
				log.Infof("adding object:%s to roster", m.Name)
				a.brokers[m.Name] = NewBroker(m.Group, m.Version, m.Name, m.Namespace)
				a.brokers[m.Name].PassMessages()
			} else {
				log.Infof("skip adding object:%s to roster", m.Name)
			}
		case "delete":
			m := crd.Parse(msg.Content)
			if cb := a.brokers[m.Name]; cb != nil {
				log.Infof("deleting object:%s from roster", m.Name)
				cb.Stop()
				delete(a.brokers, m.Name)
			} else {
				log.Infof("skip deleting object:%s to roster", m.Name)
			}
		case "update":
			m := crd.Parse(msg.Content)
			log.Infof("skip updating object:%s", m.Name)
		default:
		}
	}
}

// Connect to k8s in-cluster api
func (a *App) Connect() *App {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Errorf("In-cluster config error:%s", err)
		os.Exit(1)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Errorf("Clientset error:%s", err)
		os.Exit(1)
	}
	a.clientset = clientset
	return a
}

// Serve ui and api endpoints
func (a *App) Serve() {
	bindAddr := fmt.Sprintf("%s:%d", a.Args().BindAddr(), a.Args().Port())
	log.Infof("Serving %s, static folder:%s", bindAddr, a.Args().Dir())
	r := mux.NewRouter()
	r.PathPrefix("/ui/").Handler(http.StripPrefix("/ui/", http.FileServer(http.Dir(a.Args().Dir()))))
	r.HandleFunc("/watch/{objects}", a.Watch)
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
