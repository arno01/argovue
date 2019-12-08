package app

import (
	"encoding/gob"
	"fmt"
	"kubevue/args"
	"kubevue/auth"
	"kubevue/broker"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// App defines application dependencies
type App struct {
	args      *args.Args
	clientset *kubernetes.Clientset
	broker    *broker.Broker
	informer  informers.GenericInformer
	auth      *auth.Auth
	store     *sessions.FilesystemStore
	stop      chan struct{}
}

// Args returns application argumnets
func (a *App) Args() *args.Args {
	return a.args
}

// Broker returns event broker
func (a *App) Broker() *broker.Broker {
	return a.broker
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

// New creates an application instance
func New() *App {
	a := new(App)
	a.broker = broker.New()
	a.args = args.New().LogLevel()
	a.auth = auth.New(a.Args().OIDC())
	a.stop = make(chan struct{})
	a.store = sessions.NewFilesystemStore("", []byte("session-secret"))
	gob.Register(map[string]interface{}{})
	return a
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
	r.HandleFunc("/sse", a.SSE)
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
