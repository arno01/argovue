package app

import (
	"kubevue/kube"
	"net/http"

	"github.com/argoproj/argo/workflow/util"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/gorilla/mux"
)

func (a *App) RetryWorkflow(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	namespace := mux.Vars(r)["namespace"]

	wfClient := kube.GetWfClient(namespace)
	wf, err := wfClient.Get(name, metav1.GetOptions{})
	if err != nil {
		log.Errorf("Can't get workflow %s/%s, error:%s", namespace, name, err)
		http.Error(w, "Can't get workflow", http.StatusInternalServerError)
	}
	kubeClient, _ := kube.GetClient()
	wf, err = util.RetryWorkflow(kubeClient, wfClient, wf)
	if err != nil {
		log.Errorf("Can't retry workflow %s/%s, error:%s", namespace, name, err)
		http.Error(w, "Can't retry workflow", http.StatusInternalServerError)
	}
}
