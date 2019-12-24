package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func (a *App) Objects(w http.ResponseWriter, r *http.Request) {
	session, _ := a.Store().Get(r, "auth-session")
	var re []string
	for name, _ := range a.brokers[session.ID] {
		re = append(re, fmt.Sprintf("%s", name))
	}
	json.NewEncoder(w).Encode(re)
}

func (a *App) watchKind(w http.ResponseWriter, r *http.Request) {
	session, _ := a.Store().Get(r, "auth-session")
	vars := mux.Vars(r)
	a.watchBroker(a.getBroker(session.ID, vars["kind"]), w, r)
}

func (a *App) watchName(w http.ResponseWriter, r *http.Request) {
	session, _ := a.Store().Get(r, "auth-session")
	vars := mux.Vars(r)
	if broker := a.getSubsetBroker(session.ID, vars["namespace"], vars["kind"], vars["name"]); broker != nil {
		a.watchBroker(broker, w, r)
		return
	}
	log.Errorf("watchName: not found %s/%s/%s", vars["namespace"], vars["kind"], vars["name"])
	http.Error(w, "Objects not found", http.StatusNotFound)
}

func (a *App) watchNameSubset(w http.ResponseWriter, r *http.Request) {
	session, _ := a.Store().Get(r, "auth-session")
	vars := mux.Vars(r)
	if vars["kind"] == "catalogue" && vars["subset"] == "instances" {
		labelSelector := fmt.Sprintf("%s=%s", "argovue.io/service", vars["name"])
		if broker := a.getSubsetBroker(session.ID, vars["namespace"], vars["kind"], vars["name"], labelSelector); broker != nil {
			a.watchBroker(broker, w, r)
			return
		}
	}
	log.Errorf("watchNameSubset: not found %s/%s/%s/%s", vars["namespace"], vars["kind"], vars["name"], vars["subset"])
	http.Error(w, "Objects not found", http.StatusNotFound)
}

func (a *App) watchBroker(cb *CrdBroker, w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	if cb != nil {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Transfer-Encoding", "identity")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)
		flusher.Flush()
		cb.broker.Serve(w, flusher)
	} else {
		log.Errorf("watchBroker: not found:%s", mux.Vars(r))
		http.Error(w, "Objects not found", http.StatusNotFound)
	}
}
