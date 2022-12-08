package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	resourceId := schema.GroupVersionResource{
		Group:    "apiextensions.crossplane.io",
		Version:  "v1",
		Resource: "compositeresourcedefinitions",
	}

	list, err := dynClient.Resource(resourceId).
		List(ctx, v1.ListOptions{})

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)

		panic(err)
	}

	fmt.Printf("XRDS: %v/n", list)

}
