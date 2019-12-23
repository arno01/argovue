package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Objects return list of known objects
func (a *App) GetObjects(sessionId string) (re []string) {
	for name, _ := range a.brokers[sessionId] {
		re = append(re, fmt.Sprintf("%s", name))
	}
	return
}

// Objects returns list of known objects
func (a *App) Objects(w http.ResponseWriter, r *http.Request) {
	session, err := a.Store().Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(a.GetObjects(session.ID))
}

// Watch writes events to SSE stream
func (a *App) Watch(w http.ResponseWriter, r *http.Request) {
	session, err := a.Store().Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	kind := mux.Vars(r)["kind"]
	name := mux.Vars(r)["name"]

	if kind == "workflows" && len(name) > 0 {
		cb := a.getBroker(session.ID, "pods")
		cb.AddCrd("", "v1", "pods", fmt.Sprintf("workflows.argoproj.io/workflow=%s", name))
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	if cb := a.getBroker(session.ID, kind); cb != nil {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Transfer-Encoding", "identity")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)
		flusher.Flush()
		cb.broker.Serve(w, name, flusher)
		log.Debugf("SSE: %s/%s close", kind, name)
	} else {
		log.Errorf("Can't subscribe to %s/%s", kind, name)
		http.Error(w, "Objects not found", http.StatusNotFound)
	}
}
