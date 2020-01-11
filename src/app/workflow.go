package app

import (
	"argovue/constant"
	"argovue/crd"
	"argovue/kube"
	"argovue/profile"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/argoproj/argo/workflow/util"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gorilla/mux"
)

func (a *App) watchWorkflow(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace := v["name"], v["namespace"]
	if err := authObj("workflow", name, namespace, p); err != nil {
		return err
	}
	cb := a.maybeNewSubsetBroker(sid, crd.Workflow(name))
	a.watchBroker(cb, w, r)
	return nil
}

func (a *App) watchWorkflowPods(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace, pod := v["name"], v["namespace"], v["pod"]
	if err := authObj("workflow", name, namespace, p); err != nil {
		return err
	}
	crd := crd.WorkflowPods(name, pod)
	cb := a.maybeNewSubsetBroker(sid, crd)
	a.watchBroker(cb, w, r)
	return nil
}

func (a *App) watchWorkflowMounts(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace := v["name"], v["namespace"]
	if err := authObj("workflow", name, namespace, p); err != nil {
		return err
	}
	id := fmt.Sprintf("%s-%s-mounts", namespace, name)
	cb := a.maybeNewIdSubsetBroker(sid, id)
	cb.AddCrd(crd.WorkflowMounts(name))
	a.watchBroker(cb, w, r)
	return nil
}

func (a *App) watchWorkflowServices(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace := v["name"], v["namespace"]
	log.Debugf("SSE: start workflow/%s/%s services", namespace, name)
	if err := authObj("workflow", name, namespace, p); err != nil {
		return err
	}
	crd := crd.WorkflowServices(name)
	cb := a.maybeNewSubsetBroker(sid, crd)
	a.watchBroker(cb, w, r)
	log.Debugf("SSE: stop workflow/%s/%s services", namespace, name)
	return nil
}

func (a *App) controlWorkflowService(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace, service, action := v["name"], v["namespace"], v["service"], v["action"]
	if err := authObj("workflow", name, namespace, p); err != nil {
		return err
	}
	var err error
	switch action {
	case "delete":
		err = crd.DeleteInstance(namespace, service)
	default:
		err = fmt.Errorf("Unknown action:%s", action)
	}
	if err != nil {
		log.Errorf("Can't %s workflow %s/%s, error:%s", action, namespace, name, err)
		sendError(w, action, err)
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "action": action, "message": ""})
	}
	return nil
}

func (a *App) watchWorkflowPodLogs(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace, pod, container := v["name"], v["namespace"], v["pod"], v["container"]
	if err := authObj("workflow", name, namespace, p); err != nil {
		return err
	}
	crd := crd.WorkflowPods(name, pod)
	broker := a.getSubsetBroker(sid, crd.Id())
	if broker == nil {
		return makeError(http.StatusForbidden, "Can't find workflow broker")
	}
	podObj := broker.Broker().Find(pod, namespace)
	if podObj == nil {
		return makeError(http.StatusForbidden, "Can't find workflow pod")
	}
	wfLabel, ok := podObj.(metav1.Object).GetLabels()[constant.WorkflowLabel]
	if !ok || wfLabel != name {
		return makeError(http.StatusForbidden, "Can't find matching workflow label")
	}
	a.streamPodLogs(w, r, pod, namespace, container)
	return nil
}

func (a *App) controlWorkflow(sid string, p *profile.Profile, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace, action := v["name"], v["namespace"], v["action"]

	if err := authObj("workflow", name, namespace, p); err != nil {
		return err
	}

	wfClientset, err := kube.GetWfClientset()
	if err != nil {
		log.Errorf("Can't get argo clientset, error:%s", err)
		sendError(w, action, err)
		return nil
	}
	wfClient := kube.GetWfClient(wfClientset, namespace)
	wf, err := wfClient.Get(name, metav1.GetOptions{})
	if err != nil {
		log.Errorf("Can't get workflow %s/%s, error:%s", namespace, name, err)
		sendError(w, action, err)
		return nil
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
		err = crd.DeployFilebrowser(wf, a.Args().Namespace(), a.Args().Release(), p.IdLabel())
	default:
		err = fmt.Errorf("unrecognized command %s", action)
	}
	if err != nil {
		log.Errorf("Can't %s workflow %s/%s, error:%s", action, namespace, name, err)
		sendError(w, action, err)
	} else {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "action": action, "message": ""})
	}
	return nil
}
