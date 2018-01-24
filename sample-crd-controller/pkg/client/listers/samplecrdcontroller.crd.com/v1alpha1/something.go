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

// This file was automatically generated by lister-gen

package v1alpha1

import (
	v1alpha1 "k8s-practice/sample-crd-controller/pkg/apis/samplecrdcontroller.crd.com/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// SomethingLister helps list Somethings.
type SomethingLister interface {
	// List lists all Somethings in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.Something, err error)
	// Somethings returns an object that can list and get Somethings.
	Somethings(namespace string) SomethingNamespaceLister
	SomethingListerExpansion
}

// somethingLister implements the SomethingLister interface.
type somethingLister struct {
	indexer cache.Indexer
}

// NewSomethingLister returns a new SomethingLister.
func NewSomethingLister(indexer cache.Indexer) SomethingLister {
	return &somethingLister{indexer: indexer}
}

// List lists all Somethings in the indexer.
func (s *somethingLister) List(selector labels.Selector) (ret []*v1alpha1.Something, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Something))
	})
	return ret, err
}

// Somethings returns an object that can list and get Somethings.
func (s *somethingLister) Somethings(namespace string) SomethingNamespaceLister {
	return somethingNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// SomethingNamespaceLister helps list and get Somethings.
type SomethingNamespaceLister interface {
	// List lists all Somethings in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.Something, err error)
	// Get retrieves the Something from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.Something, error)
	SomethingNamespaceListerExpansion
}

// somethingNamespaceLister implements the SomethingNamespaceLister
// interface.
type somethingNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Somethings in the indexer for a given namespace.
func (s somethingNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Something, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Something))
	})
	return ret, err
}

// Get retrieves the Something from the indexer for a given namespace and name.
func (s somethingNamespaceLister) Get(name string) (*v1alpha1.Something, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("something"), name)
	}
	return obj.(*v1alpha1.Something), nil
}
