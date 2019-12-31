package app

import (
	"argovue/crd"
	"argovue/kube"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

func (a *App) watchObject(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	kind, name, namespace := v["kind"], v["name"], v["namespace"]
	session, err := a.Store().Get(r, "auth-session")
	if err != nil {
		log.Debugf("Can't get session")
		http.Error(w, "Can't get session", http.StatusInternalServerError)
		return
	}
	log.Debugf("SSE: start %s/%s/%s", kind, namespace, name)
	obj, err := kube.GetByKind(kind, name, namespace)
	if err != nil {
		log.Debugf("Can't find object %s/%s/%s", kind, namespace, name)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if !authHTTP(obj, a.Store(), r) {
		log.Debugf("Not authorized to access object %s/%s/%s", kind, namespace, name)
		http.Error(w, "Not authorized", http.StatusForbidden)
		return
	}
	crd, err := crd.GetByKind(kind, namespace, name)
	if err != nil {
		log.Debugf("Can't create watcher %s/%s/%s", kind, namespace, name)
		http.Error(w, "Watcher not found", http.StatusNotFound)
		return
	}
	cb := a.maybeNewSubsetBroker(session.ID, crd)
	a.watchBroker(cb, w, r)
	log.Debugf("SSE: stop %s/%s/%s", kind, namespace, name)
}
