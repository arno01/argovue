package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Event struct {
}

func sendMsg(w http.ResponseWriter, flusher http.Flusher, m *Event) {
	jsonMsg, err := json.Marshal(m)
	if err != nil {
		log.Errorf("Can't encode event:%s", m)
		return
	}
	w.Write([]byte(fmt.Sprintf("data: %s\n\n", jsonMsg)))
	flusher.Flush()
}

func (a *App) handleEvents(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Transfer-Encoding", "identity")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()
	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		if session, err := a.Store().Get(r, "auth-session"); err == nil {
			log.Debugf("Events: close connection, session id:%s", session.ID)
			a.onLogout(session.ID)
		}
	}()

	for {
		m, open := <-a.events
		if !open {
			break
		}
		sendMsg(w, flusher, m)
	}
}
