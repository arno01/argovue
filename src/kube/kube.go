package kube

import (
	"io"
	"os"
	"path/filepath"

	versioned "github.com/argoproj/argo/pkg/client/clientset/versioned"
	"github.com/argoproj/argo/pkg/client/clientset/versioned/typed/workflow/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetConfig() (*rest.Config, error) {
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		return rest.InClusterConfig()
	}
	kubeConfigPath := os.Getenv("KUBECONFIG")
	if kubeConfigPath == "" {
		kubeConfigPath = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}
	return clientcmd.BuildConfigFromFlags("", kubeConfigPath)
}

func GetClient() (*kubernetes.Clientset, error) {
	config, _ := GetConfig()
	return kubernetes.NewForConfig(config)
}

func GetPodLogs(name, namespace, container string) (io.ReadCloser, error) {
	config, _ := GetConfig()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	podLogOpts := corev1.PodLogOptions{Container: container, Follow: false}
	req := clientset.CoreV1().Pods(namespace).GetLogs(name, &podLogOpts)
	return req.Stream()
}

func GetWfClient(namespace string) v1alpha1.WorkflowInterface {
	config, _ := GetConfig()
	wfClientset := versioned.NewForConfigOrDie(config)
	wfClient := wfClientset.ArgoprojV1alpha1().Workflows(namespace)
	return wfClient
}
