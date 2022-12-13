package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config.yaml"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// Create configuration
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: *kubeconfig},
		&clientcmd.ConfigOverrides{
			CurrentContext: "kind-kind",
		}).ClientConfig()
	handleError(err)

	// Create client
	c, err := client.New(config, client.Options{})
	handleError(err)

	// List pods using k8s.io/api types
	pods := &corev1.PodList{}
	err = c.List(context.Background(), pods)
	handleError(err)

	fmt.Println("PODS Using k8s.io/api types:")
	fmt.Println(pods)

	// List pods using unstructured.Unstructured
	u := &unstructured.UnstructuredList{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "",
		Kind:    "PodList",
		Version: "v1",
	})
	err = c.List(context.Background(), u)
	handleError(err)

	fmt.Println("Unstructured: ")
	for _, pod := range u.Items {
		fmt.Println(pod)
	}
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		panic(err)
	}
}
