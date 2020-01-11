package app

import (
	"argovue/kube"
	"bufio"
	"fmt"
	"net/http"

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
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Transfer-Encoding", "identity")
	w.Header().Set("Access-Control-Allow-Origin", a.Args().UIRootDomain())
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	reader := bufio.NewReader(stream)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Errorf("Log stream read error:%s", err)
			break
		}
		_, err = w.Write([]byte(fmt.Sprintf("data: %s\n\n", str)))
		if err != nil {
			log.Errorf("Log stream write error:%s", err)
			break
		}
		flusher.Flush()
	}
	<-w.(http.CloseNotifier).CloseNotify()
	log.Debugf("Logs: close connection,")
}
