package app

import (
	"argovue/crd"
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

func (a *App) authWorkflowPod(sessionId, name, namespace string) bool {
	if broker := a.getBroker(sessionId, "workflows"); broker != nil {
		if wf := broker.Broker().Find(name, namespace); wf != nil {
			return true
		} else {
			log.Debugf("authWorkflowPod: no workflow %s/%s", namespace, name)
			return false
		}
	}
	log.Debugf("authWorkflowPod: no workflow broker")
	return false
}

func (a *App) watchWorkflowPods(w http.ResponseWriter, r *http.Request) {
	session, _ := a.Store().Get(r, "auth-session")
	v := mux.Vars(r)
	name, namespace, pod := v["name"], v["namespace"], v["pod"]
	log.Debugf("SSE: start workflow/%s/%s/%s", namespace, name, pod)
	if !a.authWorkflowPod(session.ID, name, namespace) {
		log.Debugf("SSE: access denied, no workflow %s/%s", namespace, name)
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}
	crd := crd.New("", "v1", "pods").
		SetLabelSelector("workflows.argoproj.io/workflow=" + name).
		SetFieldSelector("metadata.name=" + pod)
	cb := a.maybeNewSubsetBroker(session.ID, crd)
	a.watchBroker(cb, w, r)
	log.Debugf("SSE: stop workflow/%s/%s/%s", namespace, name, pod)
}

func (a *App) watchWorkflow(w http.ResponseWriter, r *http.Request) {
	session, _ := a.Store().Get(r, "auth-session")
	v := mux.Vars(r)
	name, namespace := v["name"], v["namespace"]
	log.Debugf("SSE: start workflow/%s/%s", namespace, name)
	if !a.authWorkflowPod(session.ID, name, namespace) {
		log.Debugf("SSE: access denied, no workflow %s/%s", namespace, name)
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}
	crd := crd.New("argoproj.io", "v1alpha1", "workflows").SetFieldSelector("metadata.name=" + name)
	cb := a.maybeNewSubsetBroker(session.ID, crd)
	a.watchBroker(cb, w, r)
	log.Debugf("SSE: stop workflow/%s/%s", namespace, name)
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
