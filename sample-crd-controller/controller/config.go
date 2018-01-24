package controller

import (
	"k8s.io/client-go/util/workqueue"
	"github.com/golang/glog"
	"k8s.io/client-go/scale/scheme/appsv1beta2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/record"

	clientset "k8s-practice/sample-crd-controller/pkg/client/clientset/versioned"
	kubeinformers "k8s.io/client-go/informers"
	informers "k8s-practice/sample-crd-controller1/pkg/client/informers/externalversions"
)

