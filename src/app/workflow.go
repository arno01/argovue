package app

import (
	"encoding/json"
	"fmt"
	"argovue/kube"
	"net/http"

	"github.com/argoproj/argo/workflow/util"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gorilla/mux"
)

func sendError(w http.ResponseWriter, action string, err error) {
	json.NewEncoder(w).Encode(map[string]string{"status": "error", "action": action, "message": fmt.Sprintf("%s", err)})
}

func (a *App) CommandWorkflow(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	namespace := mux.Vars(r)["namespace"]
	action := mux.Vars(r)["action"]

	wfClientset, err := kube.GetWfClientset()
	if err != nil {
		log.Errorf("Can't get argo clientset, error:%s", err)
		sendError(w, action, err)
		return
	}
	wfClient := kube.GetWfClient(wfClientset, namespace)
	wf, err := wfClient.Get(name, metav1.GetOptions{})
	if err != nil {
		log.Errorf("Can't get workflow %s/%s, error:%s", namespace, name, err)
		sendError(w, action, err)
		return
	}
	kubeClient, _ := kube.GetClient()
	switch action {
	case "retry":
		_, err = util.RetryWorkflow(kubeClient, wfClient, wf)
	case "resubmit":
		newWF, err := util.FormulateResubmitWorkflow(wf, false)
		if err == nil {
			_, err = util.SubmitWorkflow(wfClient, wfClientset, namespace, newWF, nil)
		}
	case "delete":
		err = wfClient.Delete(name, &metav1.DeleteOptions{})
	case "suspend":
		err = util.SuspendWorkflow(wfClient, name)
	case "resume":
		err = util.ResumeWorkflow(wfClient, name)
	case "terminate":
		err = util.TerminateWorkflow(wfClient, name)
	case "default":
		err = fmt.Errorf("unrecognized command %s", action)
	}
	if err != nil {
		log.Errorf("Can't %s workflow %s/%s, error:%s", action, namespace, name, err)
		sendError(w, action, err)
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "action": action, "message": ""})
	}
}
