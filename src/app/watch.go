package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func (a *App) onLogin(sessionId string, profile map[string]interface{}) {
	log.Debugf("profile:%s", profile)
	groups := profile["groups"].([]interface{})
	wfBroker := a.newBroker(sessionId, "workflows")
	svcBroker := a.newBroker(sessionId, "services")
	if len(groups) > 0 {
		strGroups := []string{}
		for _, group := range groups {
			if strGroup, ok := group.(string); ok {
				strGroups = append(strGroups, strGroup)
			}
		}
		selector := fmt.Sprintf("oidc.argovue.io/group in (%s)", strings.Join(strGroups, ","))
		wfBroker.AddCrd("argoproj.io", "v1alpha1", "workflows", selector)
		svcBroker.AddCrd("", "v1", "services", selector)
	}
	if sub, ok := profile["sub"].(string); ok {
		selector := fmt.Sprintf("oidc.argovue.io/id in (%s)", sub)
		wfBroker.AddCrd("argoproj.io", "v1alpha1", "workflows", selector)
		svcBroker.AddCrd("", "v1", "services", selector)
	}
	wfBroker.PassMessages()
	svcBroker.PassMessages()
}

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
