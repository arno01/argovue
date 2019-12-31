package kube

import (
	"fmt"

	argovuev1 "argovue/apis/argovue.io/v1"
	v1alpha1 "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetByKind(kind, name, namespace string) (metav1.Object, error) {
	switch kind {
	case "service":
		return GetService(name, namespace)
	case "pod":
		return GetPod(name, namespace)
	case "deployment":
		return GetDeployment(name, namespace)
	case "workflow":
		return GetWorkflow(name, namespace)
	case "argovue":
		return GetArgovueService(name, namespace)
	default:
		return nil, fmt.Errorf("Unknown kubernetes kind %s", kind)
	}
}

func GetArgovueService(name, namespace string) (*argovuev1.Service, error) {
	clientset, err := GetV1Clientset()
	if err != nil {
		return nil, err
	}
	return clientset.ArgovueV1().Services(namespace).Get(name, metav1.GetOptions{})
}

func GetWorkflow(name, namespace string) (*v1alpha1.Workflow, error) {
	clientset, err := GetWfClientset()
	if err != nil {
		return nil, err
	}
	return clientset.ArgoprojV1alpha1().Workflows(namespace).Get(name, metav1.GetOptions{})
}

func GetService(name, namespace string) (*corev1.Service, error) {
	clientset, err := GetClient()
	if err != nil {
		return nil, err
	}
	return clientset.CoreV1().Services(namespace).Get(name, metav1.GetOptions{})
}

func GetPod(name, namespace string) (*corev1.Pod, error) {
	clientset, err := GetClient()
	if err != nil {
		return nil, err
	}
	return clientset.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
}

func GetDeployment(name, namespace string) (*appsv1.Deployment, error) {
	clientset, err := GetClient()
	if err != nil {
		return nil, err
	}
	return clientset.AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{})
}
