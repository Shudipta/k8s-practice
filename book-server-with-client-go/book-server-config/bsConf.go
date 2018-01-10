/*
Copyright 2017 The Kubernetes Authors.
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

// Note: the example only works with the code within the same release/branch.
package book_server_config

import (
	"fmt"

	appsv1beta2 "k8s.io/api/apps/v1beta2"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/apimachinery/pkg/util/intstr"
	//"k8s.io/client-go/util/retry"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"log"
	//"k8s.io/client-go/scale/scheme/appsv1beta1"
)

func CreateDeployment(clientset *kubernetes.Clientset) {
	deploymentsClient := clientset.AppsV1beta2().Deployments(apiv1.NamespaceDefault)

	deployment := &appsv1beta2.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "book-server-deployment",
		},
		Spec: appsv1beta2.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "book-server",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "book-server",
					},
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

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(deployment)
	if err != nil {
		log.Fatal("error in creating deployment", err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
}

// Delete Deployment & Service
func DeleteDeploymentAndService(clientset *kubernetes.Clientset) {
	// delete deployment
	deploymentsClient := clientset.AppsV1beta2().Deployments(apiv1.NamespaceDefault)
	fmt.Println("Deleting deployment...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete("book-server-deployment", &metav1.DeleteOptions {
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted deployment.")

	// delete service
	serviceClient := clientset.CoreV1().Services(apiv1.NamespaceDefault)
	log.Println("Deleting service...")
	if err := serviceClient.Delete("book-server-service", &metav1.DeleteOptions {
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		log.Fatal(err)
	}
	log.Println("service deleted")
}

func CreateService(clientset *kubernetes.Clientset) {
	servicesClient := clientset.CoreV1().Services(apiv1.NamespaceDefault)

	svc := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "book-server-service",
			Labels: map[string]string{
				"app": "book-server",
			},
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Protocol: apiv1.ProtocolTCP,
					Port: 80,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 10000,
					},
				},
			},
			Type:     apiv1.ServiceTypeNodePort,
			Selector: map[string]string {
				"app": "book-server",
			},
		},
	}

	// Create Service
	fmt.Println("Creating service...")
	result, err := servicesClient.Create(svc)
	if err != nil {
		log.Fatal("error in creating service", err)
	}
	fmt.Printf("Created service %q.\n", result.GetObjectMeta().GetName())

	// The url at which we can access the now
	node, err := clientset.CoreV1().Nodes().Get("minikube", metav1.GetOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("The url at which we can access the now is")
	fmt.Printf("%v:%v\n", node.Status.Addresses[0].Address, result.Spec.Ports[0].NodePort)
}

func int32Ptr(i int32) *int32 { return &i }