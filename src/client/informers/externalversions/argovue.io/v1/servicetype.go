/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	versioned "argovue/client/clientset/versioned"
	internalinterfaces "argovue/client/informers/externalversions/internalinterfaces"
	time "time"

	argovueiov1 "argovue/apis/argovue.io/v1"
	v1 "argovue/client/listers/argovue.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ServiceTypeInformer provides access to a shared informer and lister for
// ServiceTypes.
type ServiceTypeInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.ServiceTypeLister
}

type serviceTypeInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewServiceTypeInformer constructs a new informer for ServiceType type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewServiceTypeInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredServiceTypeInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredServiceTypeInformer constructs a new informer for ServiceType type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredServiceTypeInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ArgovueV1().ServiceTypes(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ArgovueV1().ServiceTypes(namespace).Watch(options)
			},
		},
		&argovueiov1.ServiceType{},
		resyncPeriod,
		indexers,
	)
}

func (f *serviceTypeInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredServiceTypeInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *serviceTypeInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&argovueiov1.ServiceType{}, f.defaultInformer)
}

func (f *serviceTypeInformer) Lister() v1.ServiceTypeLister {
	return v1.NewServiceTypeLister(f.Informer().GetIndexer())
}
