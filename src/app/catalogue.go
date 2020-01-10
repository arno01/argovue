package app

import (
	"argovue/crd"
	"argovue/kube"
	"argovue/util"
	"encoding/json"
	"fmt"
	"net/http"

	argovuev1 "argovue/apis/argovue.io/v1"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

func (a *App) authObj(kind, name, namespace string, w http.ResponseWriter, r *http.Request) bool {
	obj, err := kube.GetByKind(kind, name, namespace)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		log.Errorf("Can't find object %s/%s/%s", kind, namespace, name)
		return false
	}
	if !authHTTP(obj, a.Store(), r) {
		http.Error(w, "Not authorized", http.StatusForbidden)
		log.Errorf("Not authorized to access object %s/%s/%s", kind, namespace, name)
		return false
	}
	return true
}

func (a *App) watchCatalogue(w http.ResponseWriter, r *http.Request) {
	session, _ := a.Store().Get(r, "auth-session")
	v := mux.Vars(r)
	name, namespace := v["name"], v["namespace"]
	if !a.authObj("argovue", name, namespace, w, r) {
		return
	}
	log.Debugf("SSE: start catalogue/%s/%s", namespace, name)
	crd := crd.Catalogue(name)
	cb := a.maybeNewSubsetBroker(session.ID, crd)
	a.watchBroker(cb, w, r)
	log.Debugf("SSE: stop catalogue/%s/%s", namespace, name)
}

func (a *App) watchCatalogueInstances(w http.ResponseWriter, r *http.Request) {
	session, _ := a.Store().Get(r, "auth-session")
	v := mux.Vars(r)
	name, namespace := v["name"], v["namespace"]
	if !a.authObj("argovue", name, namespace, w, r) {
		return
	}
	log.Debugf("SSE: start catalogue/%s/%s instances", namespace, name)
	crd := crd.CatalogueInstances(name)
	cb := a.maybeNewSubsetBroker(session.ID, crd)
	a.watchBroker(cb, w, r)
	log.Debugf("SSE: stop catalogue/%s/%s instances", namespace, name)
}

func (a *App) watchCatalogueResources(w http.ResponseWriter, r *http.Request) {
	session, _ := a.Store().Get(r, "auth-session")
	v := mux.Vars(r)
	name, namespace := v["name"], v["namespace"]
	if !a.authObj("argovue", name, namespace, w, r) {
		return
	}
	log.Debugf("SSE: start catalogue/%s/%s resources", namespace, name)
	crd := crd.CatalogueResources(name)
	cb := a.maybeNewSubsetBroker(session.ID, crd)
	a.watchBroker(cb, w, r)
	log.Debugf("SSE: stop catalogue/%s/%s resources", namespace, name)
}

func (a *App) watchCatalogueInstanceResources(w http.ResponseWriter, r *http.Request) {
	session, _ := a.Store().Get(r, "auth-session")
	v := mux.Vars(r)
	name, namespace, instance := v["name"], v["namespace"], v["instance"]
	if !a.authObj("helmrelease", instance, namespace, w, r) {
		return
	}
	log.Debugf("SSE: start catalogue/%s/%s/%s resources", namespace, name, instance)
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
	if !a.authObj("helmrelease", instance, namespace, w, r) {
		return
	}
	log.Debugf("SSE: start catalogue/%s/%s/%s", namespace, name, instance)
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
	if !a.authObj("argovue", name, namespace, w, r) {
		return
	}
	switch action {
	case "deploy":
		var svc *argovuev1.Service
		var label, owner string

		profile := session.Values["profile"].(map[string]interface{})
		svc, err = kube.GetArgovueService(name, namespace)
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
	v := mux.Vars(r)
	name, namespace, instance, action := v["name"], v["namespace"], v["instance"], v["action"]
	if !a.authObj("helmrelease", instance, namespace, w, r) {
		return
	}
	switch action {
	case "delete":
		err := crd.DeleteInstance(namespace, instance)
		if err != nil {
			log.Errorf("Can't %s catalogue %s/%s instance:%s, error:%s", action, namespace, name, instance, err)
			sendError(w, action, err)
		} else {
			json.NewEncoder(w).Encode(map[string]string{"status": "ok", "action": action, "message": ""})
		}
	}
}

func verifyOwner(profile map[string]interface{}, owner string) (string, string, error) {
	if util.I2s(profile["effective_id"]) == owner {
		return "oidc.argovue.io/id", util.EncodeLabel(owner), nil
	}
	if groups, ok := profile["effective_groups"].([]string); ok && len(groups) > 0 {
		for _, g := range groups {
			if g == owner {
				return "oidc.argovue.io/group", owner, nil
			}
		}
	}
	return "", "", fmt.Errorf("Can't verify owner:%s", owner)
}
