package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Objects returns list of known objects
func (a *App) Objects(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(a.GetObjects())
}

// Watch writes events to SSE stream
func (a *App) Watch(w http.ResponseWriter, r *http.Request) {

	name := mux.Vars(r)["objects"]
	namespace := mux.Vars(r)["namespace"]

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	if cb := a.getBroker(name, namespace); cb != nil {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Transfer-Encoding", "identity")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Connection", "keep-alive")
		w.WriteHeader(http.StatusOK)
		flusher.Flush()
		cb.broker.Serve(w, flusher)
		log.Debugf("Closing SSE connection")
	} else {
		log.Errorf("Can't subscribe to %s/%s", namespace, name)
		http.Error(w, "Objects not found", http.StatusNotFound)
	}
}
