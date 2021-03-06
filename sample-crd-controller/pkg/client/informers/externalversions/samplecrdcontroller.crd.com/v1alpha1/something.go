/*
Copyright 2018 The Kubernetes Authors.

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

// This file was automatically generated by informer-gen

package v1alpha1

import (
	samplecrdcontroller_crd_com_v1alpha1 "github.com/shudipta/k8s-practice/sample-crd-controller/pkg/apis/samplecrdcontroller.crd.com/v1alpha1"
	versioned "github.com/shudipta/k8s-practice/sample-crd-controller/pkg/client/clientset/versioned"
	internalinterfaces "github.com/shudipta/k8s-practice/sample-crd-controller/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/shudipta/k8s-practice/sample-crd-controller/pkg/client/listers/samplecrdcontroller.crd.com/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	time "time"
)

// SomethingInformer provides access to a shared informer and lister for
// Somethings.
type SomethingInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.SomethingLister
}

type somethingInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewSomethingInformer constructs a new informer for Something type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewSomethingInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredSomethingInformer(client, namespace, resyncPeriod, indexers)
}

// NewFilteredSomethingInformer constructs a new informer for Something type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredSomethingInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				//if tweakListOptions != nil {
				//	tweakListOptions(&options)
				//}
				return client.SamplecrdcontrollerV1alpha1().Somethings(namespace).List(v1.ListOptions{IncludeUninitialized:true})
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				//if tweakListOptions != nil {
				//	tweakListOptions(&options)
				//}
				return client.SamplecrdcontrollerV1alpha1().Somethings(namespace).Watch(v1.ListOptions{IncludeUninitialized:true})
			},
		},
		&samplecrdcontroller_crd_com_v1alpha1.Something{},
		resyncPeriod,
		indexers,
	)
}

func (f *somethingInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredSomethingInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc})
}

func (f *somethingInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&samplecrdcontroller_crd_com_v1alpha1.Something{}, f.defaultInformer)
}

func (f *somethingInformer) Lister() v1alpha1.SomethingLister {
	return v1alpha1.NewSomethingLister(f.Informer().GetIndexer())
}
