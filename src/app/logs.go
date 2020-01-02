package app

import (
	"argovue/kube"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func (a *App) streamPodLogs(w http.ResponseWriter, r *http.Request, name, namespace, container string) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	stream, err := kube.GetPodLogs(name, namespace, container)
	if err != nil {
		log.Errorf("Error getting pod logs %s/%s/%s, error:%s", namespace, name, container, err)
		http.Error(w, "Error getting logs", http.StatusInternalServerError)
		return
	}

	defer stream.Close()
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	for {
		n, err := io.CopyN(w, stream, 1024)
		if err != nil {
			log.Errorf("Log stream error:%s", err)
			break
		}
		if n < 1024 {
			break
		}
		flusher.Flush()
	}
}

func (a *App) streamLogs(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	namespace := v["namespace"]
	pod := v["pod"]
	container := v["container"]
	a.streamPodLogs(w, r, pod, namespace, container)
}
