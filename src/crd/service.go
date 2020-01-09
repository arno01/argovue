package crd

import (
	"fmt"
	"strconv"

	"argovue/kube"

	argovuev1 "argovue/apis/argovue.io/v1"
	fluxv1 "github.com/fluxcd/helm-operator/pkg/apis/helm.fluxcd.io/v1"

	wfv1alpha1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func int32Ptr(i int32) *int32 { return &i }

func Deploy(s *argovuev1.Service, owner string, input []argovuev1.InputValue) error {
	clientset, err := kube.GetFluxV1Clientset()
	if err != nil {
		return err
	}
	releaseName := fmt.Sprintf("%s-%s", s.Name, getInstanceId(s))
	release := &fluxv1.HelmRelease{
		ObjectMeta: metav1.ObjectMeta{
			Name:      releaseName,
			Namespace: s.Namespace,
			Labels: map[string]string{
				"service.argovue.io/name": s.Name,
				"oidc.argovue.io/id":      owner,
			},
			OwnerReferences: []metav1.OwnerReference{{APIVersion: "argovue.io/v1", Kind: "Service", Name: s.Name, UID: s.UID}},
		},
		Spec: s.Spec.HelmRelease,
	}
	release.Spec.ReleaseName = releaseName
	release.Spec.Values["argovue"] = map[string]string{"owner": owner}
	_, err = clientset.HelmV1().HelmReleases(s.GetNamespace()).Create(release)
	return err
}

func Delete(s *argovuev1.Service, instance string) error {
	return nil
}

func DeleteService(namespace, name string) error {
	clientset, err := kube.GetV1Clientset()
	if err != nil {
		return err
	}
	deletePolicy := metav1.DeletePropagationForeground
	opts := &metav1.DeleteOptions{PropagationPolicy: &deletePolicy}
	return clientset.ArgovueV1().Services(namespace).Delete(name, opts)
}

func DeployFilebrowser(wf *wfv1alpha1.Workflow, owner string) error {
	return nil
}

func GetWorkflowFilebrowserNames(wf *wfv1alpha1.Workflow) (re []string) {
	clientset, err := kube.GetV1Clientset()
	if err != nil {
		log.Errorf("Can't get clientset, error:%s", err)
		return
	}
	iface := clientset.ArgovueV1().Services(wf.Namespace)

	list, err := iface.List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("workflows.argoproj.io/workflow=%s,service.argovue.io/type=%s", wf.Name, "filebrowser")})
	if err != nil {
		return
	}
	for _, svc := range list.Items {
		re = append(re, svc.GetName())
	}
	return
}

func getInstanceId(s *argovuev1.Service) string {
	clientset, err := kube.GetV1Clientset()
	if err != nil {
		log.Errorf("Can't get clientset, error:%s", err)
		return "0"
	}
	freshCopy, err := clientset.ArgovueV1().Services(s.Namespace).Get(s.Name, metav1.GetOptions{})
	if err != nil {
		log.Errorf("Can't get object, error:%s", err)
		return "0"
	}
	ann := freshCopy.GetAnnotations()
	if ann == nil {
		ann = make(map[string]string)
	}
	id, ok := ann["instance.argovue.io/id"]
	if !ok {
		id = "1"
	} else {
		idI, err := strconv.Atoi(id)
		if err != nil {
			idI = 1
		}
		id = strconv.Itoa(idI + 1)
	}
	ann["instance.argovue.io/id"] = id
	freshCopy.SetAnnotations(ann)
	_, err = clientset.ArgovueV1().Services(s.Namespace).Update(freshCopy)
	if err != nil {
		log.Errorf("Can't update object, error:%s", err)
	}
	return id
}

func getFilebrowserInstanceId(wf *wfv1alpha1.Workflow) string {
	clientset, err := kube.GetWfClientset()
	if err != nil {
		log.Errorf("Can't get clientset, error:%s", err)
		return "0"
	}
	iface := clientset.ArgoprojV1alpha1().Workflows(wf.Namespace)
	freshCopy, err := iface.Get(wf.Name, metav1.GetOptions{})
	if err != nil {
		log.Errorf("Can't get object, error:%s", err)
		return "0"
	}
	ann := freshCopy.GetAnnotations()
	if ann == nil {
		ann = make(map[string]string)
	}
	id, ok := ann["instance.argovue.io/id"]
	if !ok {
		id = "1"
	} else {
		idI, err := strconv.Atoi(id)
		if err != nil {
			idI = 1
		}
		id = strconv.Itoa(idI + 1)
	}
	ann["instance.argovue.io/id"] = id
	freshCopy.SetAnnotations(ann)
	_, err = iface.Update(freshCopy)
	if err != nil {
		log.Errorf("Can't update object, error:%s", err)
	}
	return id
}
