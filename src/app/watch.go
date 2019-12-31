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
	kind := vars["kind"]
	log.Debugf("SSE: start kind %s", kind)
	broker := a.getBroker(session.ID, kind)
	if broker == nil {
		http.Error(w, "Objects not found", http.StatusNotFound)
		return
	}
	a.watchBroker(broker, w, r)
	log.Debugf("SSE: stop kind %s", kind)
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
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)
		flusher.Flush()
		cb.broker.Serve(w, flusher)
	} else {
		log.Errorf("watchBroker: not found:%s", mux.Vars(r))
		http.Error(w, "Objects not found", http.StatusNotFound)
	}
}
