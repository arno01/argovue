package app

import (
	"argovue/crd"
	"argovue/util"
	"encoding/json"
	"fmt"
	"net/http"

	argovuev1 "argovue/apis/argovue.io/v1"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

func (a *App) ServiceExists(sessionId, name, namespace string) bool {
	broker := a.getBroker(sessionId, "catalogue")
	if broker == nil {
		return false
	}
	wf := broker.Broker().Find(name, namespace)
	if wf == nil {
		return false
	}
	return true
}

func (a *App) checkServiceExists(sessionId, name, namespace string, w http.ResponseWriter) bool {
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
	if !a.checkServiceExists(session.ID, name, namespace, w) {
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
	if !a.checkServiceExists(session.ID, name, namespace, w) {
		return
	}
	crd := crd.CatalogueInstances(name)
	cb := a.maybeNewSubsetBroker(session.ID, crd)
	a.watchBroker(cb, w, r)
	log.Debugf("SSE: stop catalogue/%s/%s instances", namespace, name)
}

func (a *App) watchCatalogueResources(w http.ResponseWriter, r *http.Request) {
	session, _ := a.Store().Get(r, "auth-session")
	v := mux.Vars(r)
	name, namespace := v["name"], v["namespace"]
	log.Debugf("SSE: start catalogue/%s/%s resources", namespace, name)
	if !a.checkServiceExists(session.ID, name, namespace, w) {
		return
	}
	crd := crd.CatalogueResources(name)
	cb := a.maybeNewSubsetBroker(session.ID, crd)
	a.watchBroker(cb, w, r)
	log.Debugf("SSE: stop catalogue/%s/%s resources", namespace, name)
}

func (a *App) watchCatalogueInstanceResources(w http.ResponseWriter, r *http.Request) {
	session, _ := a.Store().Get(r, "auth-session")
	v := mux.Vars(r)
	name, namespace, instance := v["name"], v["namespace"], v["instance"]
	log.Debugf("SSE: start catalogue/%s/%s/%s resources", namespace, name, instance)
	if !a.checkServiceExists(session.ID, name, namespace, w) {
		return
	}
	id := fmt.Sprintf("%s-%s-%s-resources", namespace, name, instance)
	cb := a.maybeNewIdSubsetBroker(session.ID, id)
	cb.AddCrd(crd.CatalogueInstancePods(instance))
	cb.AddCrd(crd.CatalogueInstancePvcs(instance))
	cb.AddCrd(crd.CatalogueInstanceServices(instance))
	a.watchBroker(cb, w, r)
	log.Debugf("SSE: stop catalogue/%s/%s/%s resources", namespace, name, instance)
}

func (a *App) watchCatalogueInstance(w http.ResponseWriter, r *http.Request) {
	session, _ := a.Store().Get(r, "auth-session")
	v := mux.Vars(r)
	name, namespace, instance := v["name"], v["namespace"], v["instance"]
	log.Debugf("SSE: start catalogue/%s/%s/%s", namespace, name, instance)
	if !a.checkServiceExists(session.ID, name, namespace, w) {
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
	v := mux.Vars(r)
	name, namespace, action := v["name"], v["namespace"], v["action"]
	if !a.checkServiceExists(session.ID, name, namespace, w) {
		return
	}
	// Access granted, perform action
	switch action {
	case "deploy":
		var svc *argovuev1.Service
		var label, owner string

		profile := session.Values["profile"].(map[string]interface{})
		svc, err = crd.Typecast(a.getBroker(session.ID, "catalogue").Broker().Find(name, namespace))
		if err != nil {
			goto err
		}
		data := struct {
			Owner string                 `json:"owner"`
			Input []argovuev1.InputValue `json:"input"`
		}{}

		err = json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			goto err
		}
		label, owner, err = verifyOwner(profile, data.Owner)
		if err != nil {
			goto err
		}
		log.Debugf("Deploy service %s/%s label:%s value:%s", namespace, name, label, owner)
		err = crd.Deploy(svc, label, owner, data.Input)
		if err != nil {
			goto err
		}
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "ok", "action": action, "message": ""})
	return
err:
	log.Errorf("Can't %s catalogue %s/%s, error:%s", action, namespace, name, err)
	sendError(w, action, err)
}

func (a *App) controlCatalogueInstance(w http.ResponseWriter, r *http.Request) {
	session, err := a.Store().Get(r, "auth-session")
	if err != nil {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	v := mux.Vars(r)
	name, namespace, instance, action := v["name"], v["namespace"], v["instance"], v["action"]
	if !a.checkServiceExists(session.ID, name, namespace, w) {
		return
	}
	holderCrd := crd.CatalogueInstances(name)
	cb := a.getSubsetBroker(session.ID, holderCrd.Id())
	if cb == nil {
		log.Debugf("authCatalogueInstance: no broker for:%s", holderCrd.Id())
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}
	inst := cb.Broker().Find(instance, namespace)
	if inst == nil {
		log.Debugf("authCatalogueInstance: no instance for:%s/%s", namespace, instance)
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}
	// Access granted, perform action
	switch action {
	case "delete":
		crd.DeleteInstance(namespace, instance)
	}
	if err != nil {
		log.Errorf("Can't %s catalogue %s/%s instance:%s, error:%s", action, namespace, name, instance, err)
		sendError(w, action, err)
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "action": action, "message": ""})
	}
}

func verifyOwner(profile map[string]interface{}, owner string) (string, string, error) {
	if util.I2s(profile["effective_id"]) == owner {
		return "oidc.argovue.io/id", util.EncodeLabel(owner), nil
	}
	if groups, ok := profile["effective_groups"].([]string); ok && len(groups) > 0 {
		log.Debugf("Checking owner:%s in groups:%s", owner, groups)
		for _, g := range groups {
			if g == owner {
				return "oidc.argovue.io/group", owner, nil
			}
		}
	}
	return "", "", fmt.Errorf("Can't verify owner:%s", owner)
}
