package app

import (
	"argovue/crd"
	"argovue/profile"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

func (a *App) watchPodLogs(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace, container, kind := v["name"], v["namespace"], v["container"], "pod"
	if err := authObj(kind, name, namespace, p); err != nil {
		return err
	}
	a.streamPodLogs(w, r, name, namespace, container)
	return nil
}

func (a *App) watchObject(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	kind, name, namespace := v["kind"], v["name"], v["namespace"]
	if err := authObj(kind, name, namespace, p); err != nil {
		return err
	}
	crd, err := crd.GetByKind(kind, namespace, name)
	if err != nil {
		return makeError(http.StatusInternalServerError, "Can't create watcher %s/%s/%s", kind, namespace, name)
	}
	cb := a.maybeNewSubsetBroker(sid, crd)
	a.watchBroker(cb, w, r)
	return nil
}

func (a *App) watchKind(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	vars := mux.Vars(r)
	kind := vars["kind"]
	broker := a.getBroker(sid, kind)
	if broker == nil {
		return makeError(http.StatusNotFound, "Can't find broker by kind:%s", kind)
	}
	a.watchBroker(broker, w, r)
	return nil
}

func (a *App) watchBroker(cb *CrdBroker, w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	if cb == nil {
		log.Errorf("watchBroker: not found:%s", mux.Vars(r))
		http.Error(w, "Objects not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Transfer-Encoding", "identity")
	w.Header().Set("Access-Control-Allow-Origin", a.Args().UIRootDomain())
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()
	cb.broker.Serve(w, flusher)
}
