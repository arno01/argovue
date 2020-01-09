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

func makeRelease(s *argovuev1.Service, owner string) *fluxv1.HelmRelease {
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
	baseUrl := fmt.Sprintf("/proxy/%s/%s/%d", s.Namespace, releaseName, 80)
	release.Spec.ReleaseName = releaseName
	if release.Spec.Values == nil {
		release.Spec.Values = make(map[string]interface{})
	}
	release.Spec.Values["argovue"] = map[string]string{"owner": owner, "baseurl": baseUrl}
	return release
}

func deployRelease(s *argovuev1.Service, release *fluxv1.HelmRelease, owner string) error {
	clientset, err := kube.GetFluxV1Clientset()
	if err != nil {
		return err
	}
	_, err = clientset.HelmV1().HelmReleases(s.GetNamespace()).Create(release)
	return err
}

func Deploy(s *argovuev1.Service, owner string, input []argovuev1.InputValue) error {
	release := makeRelease(s, owner)
	return deployRelease(s, release, owner)
}

func DeployFilebrowser(wf *wfv1alpha1.Workflow, namespace, owner string) error {
	clientset, err := kube.GetV1Clientset()
	if err != nil {
		return err
	}
	filebrowser, err := clientset.ArgovueV1().Services(namespace).Get("filebrowser", metav1.GetOptions{})
	if err != nil {
		return err
	}
	release := makeRelease(filebrowser, owner)
	volumes := []map[string]string{}
	for _, pvc := range wf.Status.PersistentVolumeClaims {
		volumes = append(volumes, map[string]string{"name": pvc.Name, "claim": pvc.PersistentVolumeClaim.ClaimName})
	}
	release.ObjectMeta.Labels["workflows.argoproj.io/workflow"] = wf.Name
	release.Spec.Values["volumes"] = volumes
	if av, ok := release.Spec.Values["argovue"].(map[string]string); ok {
		av["workflow"] = wf.Name
	}
	return deployRelease(filebrowser, release, owner)
}

func Delete(s *argovuev1.Service, name string) error {
	clientset, err := kube.GetFluxV1Clientset()
	if err != nil {
		return err
	}
	deletePolicy := metav1.DeletePropagationForeground
	opts := &metav1.DeleteOptions{PropagationPolicy: &deletePolicy}
	return clientset.HelmV1().HelmReleases(s.GetNamespace()).Delete(name, opts)
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

func GetWorkflowFilebrowserNames(wf *wfv1alpha1.Workflow) (re []string) {
	clientset, err := kube.GetV1Clientset()
	if err != nil {
		log.Errorf("Can't get clientset, error:%s", err)
		return
	}
	iface := clientset.ArgovueV1().Services(wf.Namespace)

	list, err := iface.List(metav1.ListOptions{
		LabelSelector: fmt.Sprintf("workflows.argoproj.io/workflow=%s,app.kubernetes.io/name=%s", wf.Name, "filebrowser")})
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
