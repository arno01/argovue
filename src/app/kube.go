package app

import (
	log "github.com/sirupsen/logrus"
	"kubevue/msg"

	v1 "k8s.io/api/core/v1"
	apiv1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

func (a *App) Notify(action string, obj interface{}) {
	mObj := obj.(apiv1.Object)
	log.Printf("%s %s@%s uid:%s", action, mObj.GetName(), mObj.GetNamespace(), mObj.GetUID())
	a.Broker().Notifier <- msg.New(action, obj)
}

func (a *App) Watch(group, version, resource string) {
	log.Debugf("Starting kubernetes watcher: %s/%s/%s", resource, version, group)
	cfg, err := rest.InClusterConfig()
	dc, err := dynamic.NewForConfig(cfg)
	if err != nil {
		log.WithError(err).Fatal("Could not generate dynamic client for config")
	}
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dc, 0, v1.NamespaceAll, nil)
	gvr := schema.GroupVersionResource{
		Group:    group,
		Version:  version,
		Resource: resource,
	}
	a.informer = factory.ForResource(gvr)

	a.informer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				a.Notify("add", obj)
			},
			DeleteFunc: func(obj interface{}) {
				a.Notify("delete", obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				a.Notify("update", newObj)
			},
		},
	)
	go a.informer.Informer().Run(a.stop)
}
