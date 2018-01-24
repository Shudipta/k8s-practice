package controller

import (
	"github.com/golang/glog"
	"fmt"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/apimachinery/pkg/api/errors"
	//"k8s.io/client-go/scale/scheme/appsv1beta2"
	"k8s.io/apimachinery/pkg/runtime/schema"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	samplecrdcontrollerv1alpha1 "k8s-practice/sample-crd-controller/pkg/apis/samplecrdcontroller.crd.com/v1alpha1"
)

// Run will start the controller.
// StopCh channel is used to send interrupt signal to stop it.
func (c *Controller) Run(threadiness int, stopCh chan struct{}) error {
	// don't let panics crash the process
	defer runtime.HandleCrash()
	// make sure the work queue is shutdown which will trigger workers to end
	defer c.deploymentsQueue.ShutDown()
	defer c.somethingsQueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	glog.Info("Starting Something controller")

	// wait for the caches to synchronize before starting the worker
	glog.Info("Waiting for informer caches to sync")
	if !cache.WaitForCacheSync(stopCh, c.deploymentsSynced, c.somethingsSynced) {
		return fmt.Errorf("Timed out waiting for caches to sync")
	}

	glog.Info("Starting workers")
	// Launch two workers to process Something resources
	// runWorker will loop until "something bad" happens.  The .Until will
	// then rekick the worker after one second
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runSomethingWorker, time.Second, stopCh)
		go wait.Until(c.runDeploymentWorker, time.Second, stopCh)
	}

	glog.Info("Started workers")
	<-stopCh
	glog.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runSomethingWorker() {
	for c.processNextSomethingWorkItem() {
	}
}

func (c *Controller) runDeploymentWorker() {
	for c.processNextDeploymetWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextSomethingWorkItem() bool {
	obj, shutdown := c.somethingsQueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.somethingsQueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.somethingsQueue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Something resource to be synced.
		if err := c.somethingSyncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.somethingsQueue.Forget(obj)
		glog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}


// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Something resource
// with the current status of the resource.
func (c *Controller) somethingSyncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the Something resource with this namespace/name
	something, err := c.somethingsLister.Somethings(namespace).Get(name)
	if err != nil {
		// The Something resource may no longer exist, in which case we stop processing.
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("something '%s' in work queue (somethingsQueue) no longer exists", key))
			return nil
		}

		return err
	}

	deploymentName := something.Spec.DeploymentName
	if deploymentName == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		runtime.HandleError(fmt.Errorf("%s: deployment name must be specified", key))
		return nil
	}

	// Get the deployment with the name specified in Something.spec
	deployment, err := c.deploymentsLister.Deployments(something.Namespace).Get(deploymentName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		deployment, err = c.kubeclientset.AppsV1beta2().Deployments(something.Namespace).Create(newDeployment(something))
	}

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	// If the Deployment is not controlled by this Something resource, we should log
	// a warning to the event recorder and ret
	if !metav1.IsControlledBy(deployment, something) {
		msg := fmt.Sprintf("Resource %q already exists and is not managed by Something", deployment.Name)
		//c.recorder.Event(foo, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}

	// If the number of the replicas on the Something resource is specified, and the
	// number does not equal the current desired replicas on the Deployment, we
	// should update the Deployment resource.
	if something.Spec.Replicas != nil && *something.Spec.Replicas != *deployment.Spec.Replicas {
		glog.V(4).Infof("SomethingR: %d, deployR: %d", *something.Spec.Replicas, *deployment.Spec.Replicas)
		deployment, err = c.kubeclientset.AppsV1beta2().Deployments(something.Namespace).Update(newDeployment(something))
	}

	// If an error occurs during Update, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	// Finally, we update the status block of the Something resource to reflect the
	// current state of the world
	err = c.updateSomethingStatus(something, deployment)
	if err != nil {
		return err
	}

	//c.recorder.Event(foo, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)
	return nil
}

func (c *Controller) updateSomethingStatus(
		something *samplecrdcontrollerv1alpha1.Something,
		deployment *appsv1beta2.Deployment) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	somethingCopy := something.DeepCopy()
	somethingCopy.Status.AvailableReplicas = deployment.Status.AvailableReplicas
	// Until #38113 is merged, we must use Update instead of UpdateStatus to
	// update the Status block of the Foo resource. UpdateStatus will not
	// allow changes to the Spec of the resource, which is ideal for ensuring
	// nothing other than resource status has been updated.
	_, err := c.clientset.SamplecrdcontrollerV1alpha1().Somethings(something.Namespace).Update(somethingCopy)
	return err
}

// newDeployment creates a new Deployment for a Foo resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the Foo resource that 'owns' it.
func newDeployment(something *samplecrdcontrollerv1alpha1.Something) *appsv1beta2.Deployment {
	labels := map[string]string{
		"app":        "book-server",
		"controller": something.Name,
	}
	return &appsv1beta2.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      something.Spec.DeploymentName,
			Namespace: something.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(something, schema.GroupVersionKind{
					Group:   samplecrdcontrollerv1alpha1.SchemeGroupVersion.Group,
					Version: samplecrdcontrollerv1alpha1.SchemeGroupVersion.Version,
					Kind:    "Something",
				}),
			},
		},
		Spec: appsv1beta2.DeploymentSpec{
			Replicas: something.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "book-server",
							Image: "shudipta/book_server:v1",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 10000,
								},
							},
						},
					},
				},
			},
		},
	}
}
