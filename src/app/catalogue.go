package app

import (
	"argovue/crd"
	"argovue/kube"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

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
