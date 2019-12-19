package crd

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"kubevue/kube"
	"kubevue/msg"

	apiv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

// Crd
type Crd struct {
	id        string
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
	log.Debugf("CRD: %s %s %s@%s uid:%s", crd.id, action, mObj.GetName(), mObj.GetNamespace(), mObj.GetUID())
	crd.notify <- msg.New(action, obj)
}

// New crd
func New(group, version, resource, namespace string) *Crd {
	crd := new(Crd)
	crd.id = fmt.Sprintf("%s/%s/%s/%s", group, version, resource, namespace)
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
	log.Debugf("CRD: %s stop", crd.id)

	close(crd.notify)
	close(crd.stop)
}

// Watch resources
func (crd *Crd) Watch() *Crd {
	log.Debugf("CRD: %s start", crd.id)

	cfg, err := kube.GetConfig()
	if err != nil {
		log.Fatalf("CRD: %s could not configure kubernetes access, error:%s", crd.id, err)
	}
	dc, err := dynamic.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("CRD: %s could not generate dynamic client for config, error:%s", crd.id, err)
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
