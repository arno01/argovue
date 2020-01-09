package app

import (
	"argovue/crd"
	"argovue/kube"
	u "argovue/util"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/argoproj/argo/workflow/util"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gorilla/mux"
)

func sendError(w http.ResponseWriter, action string, err error) {
	json.NewEncoder(w).Encode(map[string]string{"status": "error", "action": action, "message": fmt.Sprintf("%s", err)})
}

func (a *App) authWorkflow(sessionId, name, namespace string, w http.ResponseWriter) bool {
	broker := a.getBroker(sessionId, "workflows")
	if broker == nil {
		log.Debugf("authWorkflow: no workflow broker")
		http.Error(w, "Access denied", http.StatusForbidden)
		return false
	}
	wf := broker.Broker().Find(name, namespace)
	if wf == nil {
		log.Debugf("authWorkflow: no workflow %s/%s in broker", namespace, name)
		http.Error(w, "Access denied", http.StatusForbidden)
		return false
	}
	return true
}

func (a *App) watchWorkflow(w http.ResponseWriter, r *http.Request) {
	session, _ := a.Store().Get(r, "auth-session")
	v := mux.Vars(r)
	name, namespace := v["name"], v["namespace"]
	log.Debugf("SSE: start workflow/%s/%s", namespace, name)
	if !a.authWorkflow(session.ID, name, namespace, w) {
		return
	}
	crd := crd.Workflow(name)
	cb := a.maybeNewSubsetBroker(session.ID, crd)
	a.watchBroker(cb, w, r)
	log.Debugf("SSE: stop workflow/%s/%s", namespace, name)
}

func (a *App) watchWorkflowPods(w http.ResponseWriter, r *http.Request) {
	session, _ := a.Store().Get(r, "auth-session")
	v := mux.Vars(r)
	name, namespace, pod := v["name"], v["namespace"], v["pod"]
	log.Debugf("SSE: start workflow/%s/%s/%s", namespace, name, pod)
	if !a.authWorkflow(session.ID, name, namespace, w) {
		return
	}
	crd := crd.WorkflowPods(name, pod)
	cb := a.maybeNewSubsetBroker(session.ID, crd)
	a.watchBroker(cb, w, r)
	log.Debugf("SSE: stop workflow/%s/%s/%s", namespace, name, pod)
}

func (a *App) watchWorkflowMounts(w http.ResponseWriter, r *http.Request) {
	session, _ := a.Store().Get(r, "auth-session")
	v := mux.Vars(r)
	name, namespace := v["name"], v["namespace"]
	log.Debugf("SSE: start workflow/%s/%s mounts", namespace, name)
	if !a.authWorkflow(session.ID, name, namespace, w) {
		return
	}
	wf, err := getWorkflow(name, namespace)
	if err != nil {
		log.Errorf("Can't get workflow, error:%s", err)
		http.Error(w, "Can't get workflow", http.StatusInternalServerError)
		return
	}
	id := fmt.Sprintf("%s-%s-mounts", namespace, name)
	cb := a.maybeNewIdSubsetBroker(session.ID, id)
	for _, svc := range crd.GetWorkflowFilebrowserNames(wf) {
		cb.AddCrd(crd.WorkflowMounts(svc))
	}
	a.watchBroker(cb, w, r)
	log.Debugf("SSE: stop workflow/%s/%s mounts", namespace, name)
}

func (a *App) watchWorkflowServices(w http.ResponseWriter, r *http.Request) {
	session, _ := a.Store().Get(r, "auth-session")
	v := mux.Vars(r)
	name, namespace := v["name"], v["namespace"]
	log.Debugf("SSE: start workflow/%s/%s services", namespace, name)
	if !a.authWorkflow(session.ID, name, namespace, w) {
		return
	}
	crd := crd.WorkflowServices(name)
	cb := a.maybeNewSubsetBroker(session.ID, crd)
	a.watchBroker(cb, w, r)
	log.Debugf("SSE: stop workflow/%s/%s services", namespace, name)
}

func (a *App) controlWorkflowService(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	name, namespace, service, action := v["name"], v["namespace"], v["service"], v["action"]
	session, err := a.Store().Get(r, "auth-session")
	if !a.authWorkflow(session.ID, name, namespace, w) {
		return
	}
	switch action {
	case "delete":
		err = crd.DeleteService(namespace, service)
	default:
	}
	if err != nil {
		log.Errorf("Can't %s workflow %s/%s, error:%s", action, namespace, name, err)
		sendError(w, action, err)
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "action": action, "message": ""})
	}
}

func (a *App) watchWorkflowPodLogs(w http.ResponseWriter, r *http.Request) {
	session, _ := a.Store().Get(r, "auth-session")
	v := mux.Vars(r)
	name, namespace, pod, container := v["name"], v["namespace"], v["pod"], v["container"]
	log.Debugf("SSE: start workflow/%s/%s/%s/%s", namespace, name, pod, container)
	if !a.authWorkflow(session.ID, name, namespace, w) {
		return
	}
	crd := crd.WorkflowPods(name, pod)
	broker := a.getSubsetBroker(session.ID, crd.Id())
	if broker == nil {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}
	podObj := broker.Broker().Find(pod, namespace)
	if podObj == nil {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}
	wfLabel, ok := podObj.(metav1.Object).GetLabels()["workflows.argoproj.io/workflow"]
	if !ok {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}
	if wfLabel != name {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}
	a.streamLogs(w, r)
	log.Debugf("SSE: stop workflow/%s/%s/%s/%s", namespace, name, pod, container)
}

func (a *App) commandWorkflow(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	name, namespace, action := v["name"], v["namespace"], v["action"]

	session, _ := a.Store().Get(r, "auth-session")
	if !a.authWorkflow(session.ID, name, namespace, w) {
		return
	}

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
	case "mount":
		profile := session.Values["profile"].(map[string]interface{})
		err = crd.DeployFilebrowser(wf, a.Args().Namespace(), u.EncodeLabel(u.I2s(profile["effective_id"])))
	default:
		err = fmt.Errorf("unrecognized command %s", action)
	}
	if err != nil {
		log.Errorf("Can't %s workflow %s/%s, error:%s", action, namespace, name, err)
		sendError(w, action, err)
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "action": action, "message": ""})
	}
}

func getWorkflow(name, namespace string) (*v1alpha1.Workflow, error) {
	wfClientset, err := kube.GetWfClientset()
	if err != nil {
		log.Errorf("Can't get argo clientset, error:%s", err)
		return nil, err
	}
	wfClient := kube.GetWfClient(wfClientset, namespace)
	wf, err := wfClient.Get(name, metav1.GetOptions{})
	if err != nil {
		log.Errorf("Can't get workflow %s/%s, error:%s", namespace, name, err)
		return nil, err
	}
	return wf, nil
}
