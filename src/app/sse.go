package app

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// SSE serves SSE stream
func (a *App) SSE(rw http.ResponseWriter, req *http.Request) {
	log.Debugf("Serving SSE connection")

	flusher, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Transfer-Encoding", "identity")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Connection", "keep-alive")

	rw.WriteHeader(http.StatusOK)
	flusher.Flush()
	a.Broker().Serve(rw, flusher)
	log.Debugf("Closing SSE connection")
}
