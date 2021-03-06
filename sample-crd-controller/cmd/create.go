// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/kubernetes"
	clientset "github.com/shudipta/k8s-practice/sample-crd-controller/pkg/client/clientset/versioned"
	"time"
	kubeinformers "k8s.io/client-go/informers"
	informers "github.com/shudipta/k8s-practice/sample-crd-controller/pkg/client/informers/externalversions"
	"github.com/shudipta/k8s-practice/sample-crd-controller/controller"
	//"k8s.io/client-go/rest"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
		fmt.Println(">>>>>>>>>>>>> kubeconfigfile: \"", kubeconfig, "\"")

		stopCh := make(chan struct{})
		defer close(stopCh)

		cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatalf("Error building kubeconfig: %s", err.Error())
		}
		kubeClient, err := kubernetes.NewForConfig(cfg)
		if err != nil {
			log.Fatalf("Error building kubernetes clientset: %s", err.Error())
		}
		exampleClient, err := clientset.NewForConfig(cfg)
		if err != nil {
			log.Fatalf("Error building example clientset: %s", err.Error())
		}

		kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)
		exampleInformerFactory := informers.NewSharedInformerFactory(exampleClient, time.Second*30)

		controller := controller.NewController(kubeClient, exampleClient, kubeInformerFactory, exampleInformerFactory)
		//go controller.Run(2, stopCh)

		go kubeInformerFactory.Start(stopCh)
		go exampleInformerFactory.Start(stopCh)

		if err = controller.Run(2, stopCh); err != nil {
			log.Fatalf("Error running controller: %s", err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
