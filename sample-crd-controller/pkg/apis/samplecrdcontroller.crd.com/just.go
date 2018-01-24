package samplecrdcontroller_crd_com

import (
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type ResourceEventHandlerFuncs struct {
	AddFunc    func(obj interface{})
	UpdateFunc func(oldObj, newObj interface{})
	DeleteFunc func(obj interface{})
}

func main() {
	for {
		desired := getDesiredState()
		current := getCurrentState()
		makeChanges(desired, current)
	}

	lw := cache.NewListWatchFromClient(
		client,
		&v1.Pod{},
		api.NamespaceAll,
		fieldSelector
	)

	store, controller := cache.NewInformer {
		&cache.ListWatch{},
		&cache.ListWatch {
			ListFunc:listFunc := func(options metav1.ListOptions) (runtime.Object, error) {
				return client.Get().
					Namespace(namespace).
					Resource(resource).
					VersionedParams(&options, metav1.ParameterCodec).
					FieldsSelectorParam(fieldSelector).
					Do().
					Get()
			},
			WatchFunc: watchFunc := func(options metav1.ListOptions) (watch.Interface, error) {
				options.Watch = true
				return client.Get().
					Namespace(namespace).
					Resource(resource).
					VersionedParams(&options, metav1.ParameterCodec).
					FieldsSelectorParam(fieldSelector).
					Watch()
			}
		},
		&v1.Pod{},
		resyncPeriod,
		cache.ResourceEventHandlerFuncs{},
	}

	lw := cache.NewListWatchFromClient(…)
	sharedInformer := cache.NewSharedInformer(lw, &api.Pod{}, resyncPeriod)

	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	controller.informer = cache.NewSharedInformer(...)
	controller.queue = workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	controller.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, controller.HasSynched)
	{
		log.Errorf("Timed out waiting for caches to sync"))
	}

	// Now start processing
	controller.runWorker()
}
