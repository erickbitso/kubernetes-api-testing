package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	var kubeconfig *string
	ctx := context.Background()
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config.yaml"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: *kubeconfig},
		&clientcmd.ConfigOverrides{
			CurrentContext: "kind-kind",
		}).ClientConfig()
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)

		panic(err)
	}

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)

		panic(err)
	}

	namespace := "kube-system"

	resourceId := schema.GroupVersionResource{
		Group:    "apps",
		Version:  "v1",
		Resource: "deployments",
	}

	list, err := dynClient.Resource(resourceId).Namespace(namespace).
		List(ctx, metav1.ListOptions{})

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)

		panic(err)
	}

	fmt.Printf("Deployments: %v\n", list)

}
