package crd

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"

	"kubevue/msg"

	apiv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

// Crd
type Crd struct {
	group     string
	version   string
	resource  string
	namespace string
	notify    chan *msg.Msg
	stop      chan struct{}
	informer  informers.GenericInformer
}

// Notify channel on new message
func (crd *Crd) Notify(action string, obj interface{}) {
	mObj := obj.(apiv1.Object)
	log.Debugf("%s %s@%s uid:%s", action, mObj.GetName(), mObj.GetNamespace(), mObj.GetUID())
	crd.notify <- msg.New(action, obj)
}

// New crd
func New(group, version, resource, namespace string) *Crd {
	crd := new(Crd)
	crd.group = group
	crd.version = version
	crd.resource = resource
	crd.namespace = namespace
	crd.notify = make(chan *msg.Msg)
	crd.stop = make(chan struct{})
	crd.Watch()
	return crd
}

// Stop watching
func (crd *Crd) Stop() {
	close(crd.notify)
	close(crd.stop)
}

func getConfig() (*rest.Config, error) {
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		return rest.InClusterConfig()
	}
	kubeConfigPath := os.Getenv("KUBECONFIG")
	if kubeConfigPath == "" {
		kubeConfigPath = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}
	return clientcmd.BuildConfigFromFlags("", kubeConfigPath)
}

// Watch resources
func (crd *Crd) Watch() *Crd {
	log.Debugf("Starting kubernetes watcher: %s/%s/%s/%s", crd.resource, crd.version, crd.group, crd.namespace)

	cfg, err := getConfig()
	dc, err := dynamic.NewForConfig(cfg)
	if err != nil {
		log.WithError(err).Fatal("Could not generate dynamic client for config")
	}
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dc, 0, crd.namespace, nil)
	gvr := schema.GroupVersionResource{Group: crd.group, Version: crd.version, Resource: crd.resource}
	crd.informer = factory.ForResource(gvr)

	crd.informer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				crd.Notify("add", obj)
			},
			DeleteFunc: func(obj interface{}) {
				crd.Notify("delete", obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				crd.Notify("update", newObj)
			},
		},
	)
	go crd.informer.Informer().Run(crd.stop)
	return crd
}

// Notifier returns notifier channel
func (crd *Crd) Notifier() chan *msg.Msg {
	return crd.notify
}
