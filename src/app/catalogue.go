package app

import (
	"argovue/crd"
	"argovue/kube"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

func (a *App) authCatalogue(sessionId, name, namespace string, w http.ResponseWriter) bool {
	broker := a.getBroker(sessionId, "catalogue")
	if broker == nil {
		log.Debugf("authCatalogue: no catalogue broker")
		http.Error(w, "Access denied", http.StatusForbidden)
		return false
	}
	wf := broker.Broker().Find(name, namespace)
	if wf == nil {
		log.Debugf("authCatalogue: no catalogue %s/%s in broker", namespace, name)
		http.Error(w, "Access denied", http.StatusForbidden)
		return false
	}
	return true
}

func (a *App) watchCatalogue(w http.ResponseWriter, r *http.Request) {
	session, _ := a.Store().Get(r, "auth-session")
	v := mux.Vars(r)
	name, namespace := v["name"], v["namespace"]
	log.Debugf("SSE: start catalogue/%s/%s", namespace, name)
	if !a.authCatalogue(session.ID, name, namespace, w) {
		return
	}
	crd := crd.Catalogue(name)
	cb := a.maybeNewSubsetBroker(session.ID, crd)
	a.watchBroker(cb, w, r)
	log.Debugf("SSE: stop catalogue/%s/%s", namespace, name)
}

func (a *App) watchCatalogueInstances(w http.ResponseWriter, r *http.Request) {
	session, _ := a.Store().Get(r, "auth-session")
	v := mux.Vars(r)
	name, namespace := v["name"], v["namespace"]
	log.Debugf("SSE: start catalogue/%s/%s instances", namespace, name)
	if !a.authCatalogue(session.ID, name, namespace, w) {
		return
	}
	crd := crd.CatalogueInstances(name)
	cb := a.maybeNewSubsetBroker(session.ID, crd)
	a.watchBroker(cb, w, r)
	log.Debugf("SSE: stop catalogue/%s/%s instances", namespace, name)
}

func (a *App) watchCatalogueInstance(w http.ResponseWriter, r *http.Request) {
	session, _ := a.Store().Get(r, "auth-session")
	v := mux.Vars(r)
	name, namespace, instance := v["name"], v["namespace"], v["instance"]
	log.Debugf("SSE: start catalogue/%s/%s/%s", namespace, name, instance)
	if !a.authCatalogue(session.ID, name, namespace, w) {
		return
	}
	crd := crd.CatalogueInstance(name, instance)
	cb := a.maybeNewSubsetBroker(session.ID, crd)
	a.watchBroker(cb, w, r)
	log.Debugf("SSE: stop catalogue/%s/%s/%s", namespace, name, instance)
}

func (a *App) commandCatalogue(w http.ResponseWriter, r *http.Request) {
	session, err := a.Store().Get(r, "auth-session")
	if err != nil {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	name := mux.Vars(r)["name"]
	namespace := mux.Vars(r)["namespace"]
	action := mux.Vars(r)["action"]
	kubeClient, _ := kube.GetClient()
	service := crd.Parse(a.getBroker(session.ID, "catalogue").Broker().Find(name, namespace))
	switch action {
	case "deploy":
		profile := session.Values["profile"].(map[string]interface{})
		service.Deploy(kubeClient, profile["sub"].(string))
	}
	if err != nil {
		log.Errorf("Can't %s catalogue %s/%s, error:%s", action, namespace, name, err)
		sendError(w, action, err)
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "action": action, "message": ""})
	}
}
