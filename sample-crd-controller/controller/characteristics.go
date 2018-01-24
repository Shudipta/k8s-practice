package controller

import (
	"k8s.io/client-go/kubernetes"
	clientset "k8s-practice/sample-crd-controller/pkg/client/clientset/versioned"
	listers "k8s-practice/sample-crd-controller/pkg/client/listers/samplecrdcontroller.crd.com/v1alpha1"
	kubelisters "k8s.io/client-go/listers/apps/v1beta2"
	"k8s.io/client-go/util/workqueue"
	//"k8s.io/client-go/tools/record"
	"k8s.io/client-go/tools/cache"
	"github.com/golang/glog"
	//"k8s.io/client-go/scale/scheme/appsv1beta2"
	kubeinformers "k8s.io/client-go/informers"
	informers "k8s-practice/sample-crd-controller/pkg/client/informers/externalversions"
	//"k8s.io/client-go/scale/scheme/appsv1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	samplecrdcontrollerv1alpha1 "k8s-practice/sample-crd-controller/pkg/apis/samplecrdcontroller.crd.com/v1alpha1"
)

// Controller implementation for Something resources
type Controller struct {
	kubeclientset		kubernetes.Interface
	clientset			clientset.Interface

	deploymentsLister	kubelisters.DeploymentLister
	somethingsLister	listers.SomethingLister

	deploymentsQueue	workqueue.RateLimitingInterface
	somethingsQueue		workqueue.RateLimitingInterface

	deploymentsInformer	cache.SharedIndexInformer
	somethingsInformer	cache.SharedIndexInformer

	deploymentsSynced	cache.InformerSynced
	somethingsSynced	cache.InformerSynced

	//kubType				string
	//recorder			record.EventRecorder
}

// NewController returns a new sample-crd-controller
func NewController(
	kubeclientset kubernetes.Interface,
	clientset clientset.Interface,
	kubeInformerFactory kubeinformers.SharedInformerFactory,
	sampleInformerFactory informers.SharedInformerFactory) *Controller {

	// obtain references to shared index informers for the Deployment and Something CRD
	// types.
	deploymentInformer := kubeInformerFactory.Apps().V1beta2().Deployments()
	somethingInformer := sampleInformerFactory.Samplecrdcontroller().V1alpha1().Somethings()

	// Create event broadcaster
	// Add sample-crd-controller types to the default Kubernetes Scheme so Events can be
	// logged for sample-crd-controller types.
	//samplescheme.AddToScheme(scheme.Scheme)
	//glog.V(4).Info("Creating event broadcaster")
	//eventBroadcaster := record.NewBroadcaster()
	//eventBroadcaster.StartLogging(glog.Infof)
	//eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	//recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset:		kubeclientset,
		clientset:			clientset,

		deploymentsLister:	deploymentInformer.Lister(),
		somethingsLister:	somethingInformer.Lister(),

		deploymentsInformer:deploymentInformer.Informer(),
		somethingsInformer:	somethingInformer.Informer(),

		deploymentsSynced:	deploymentInformer.Informer().HasSynced,
		somethingsSynced:	somethingInformer.Informer().HasSynced,

		deploymentsQueue:	workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Deployments"),
		somethingsQueue:	workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Somethings"),
		//recorder:          recorder,
	}

	glog.Info("Setting up event handlers")
	// Set up an event handler for when Foo resources change
	controller.somethingsInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		//AddFunc: controller.handleObject,
		AddFunc: func(obj interface{}) {
			if key, err := cache.MetaNamespaceKeyFunc(obj); err == nil {
				controller.somethingsQueue.AddRateLimited(key)
			}
		},
		UpdateFunc: func(old, new interface{}) {
			newSomething := new.(*samplecrdcontrollerv1alpha1.Something)
			oldSomething := old.(*samplecrdcontrollerv1alpha1.Something)

			if oldSomething.ResourceVersion == newSomething.ResourceVersion {
				// Periodic resync will send update events for all known Somethings.
				// Two different versions of the same Somthing will always have different RVs.
				return
			} else {
				if key, err := cache.MetaNamespaceKeyFunc(new); err == nil {
					controller.somethingsQueue.Add(key)
				}
			}
			//controller.handleObject(new)
		},
		//DeleteFunc: controller.handleObject,
		DeleteFunc: func(obj interface{}) {
			if key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj); err == nil {
				controller.somethingsQueue.Add(key)
			}
		},
	})
	// Set up an event handler for when Deployment resources change. This
	// handler will lookup the owner of the given Deployment, and if it is
	// owned by a Something resource will enqueue that Something resource for
	// processing. This way, we don't need to implement custom logic for
	// handling Deployment resources. More info on this pattern:
	// https://github.com/kubernetes/community/blob/8cafef897a22026d42f5e5bb3f104febe7e29830/contributors/devel/controllers.md
	controller.deploymentsInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		//AddFunc: controller.handleObject,
		AddFunc: func(obj interface{}) {
			if key, err := cache.MetaNamespaceKeyFunc(obj); err == nil {
				controller.deploymentsQueue.Add(key)
			}
		},
		UpdateFunc: func(old, new interface{}) {
			newDeploy := new.(*appsv1beta2.Deployment)
			oldDeploy := old.(*appsv1beta2.Deployment)

			if newDeploy.ResourceVersion == oldDeploy.ResourceVersion {
				// Periodic resync will send update events for all known Deployments.
				// Two different versions of the same Deployment will always have different RVs.
				return
			} else {
				if key, err := cache.MetaNamespaceKeyFunc(new); err == nil {
					controller.deploymentsQueue.Add(key)
				}
			}
			//controller.handleObject(new)
		},
		//DeleteFunc: controller.handleObject,
		DeleteFunc: func(obj interface{}) {
			if key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj); err == nil {
				controller.deploymentsQueue.Add(key)
			}
		},
	})

	return controller
}