package crd

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"argovue/kube"
	"argovue/msg"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

// Crd
type Crd struct {
	id       string
	group    string
	version  string
	resource string
	selector string
	notify   chan *msg.Msg
	stop     chan struct{}
	informer informers.GenericInformer
}

// Notify channel on new message
func (crd *Crd) Notify(action string, obj interface{}) {
	mObj := obj.(metav1.Object)
	log.Debugf("CRD: %s %s %s@%s uid:%s", crd.id, action, mObj.GetName(), mObj.GetNamespace(), mObj.GetUID())
	crd.notify <- msg.New(action, obj)
}

// New crd
func New(group, version, resource, selector string) *Crd {
	crd := new(Crd)
	crd.id = MakeId(group, version, resource, selector)
	crd.group = group
	crd.version = version
	crd.resource = resource
	crd.selector = selector
	crd.notify = make(chan *msg.Msg)
	crd.stop = make(chan struct{})
	crd.Watch()
	return crd
}

func MakeId(group, version, resource, selector string) string {
	return fmt.Sprintf("%s/%s/%s selector:%s", group, version, resource, selector)
}

func (crd *Crd) Id() string {
	return crd.id
}

// Stop watching
func (crd *Crd) Stop() {
	log.Debugf("CRD: %s stop", crd.id)

	close(crd.notify)
	close(crd.stop)
}

func (crd *Crd) tweakListOptions(opts *metav1.ListOptions) {
	opts.LabelSelector = crd.selector
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
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dc, 0, metav1.NamespaceAll, crd.tweakListOptions)
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
