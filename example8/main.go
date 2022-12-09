package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"time"

	ec2 "github.com/crossplane-contrib/provider-aws/apis/ec2/v1beta1"
	xpv1 "github.com/crossplane/crossplane/apis/pkg/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
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

	// New scheme
	s := runtime.NewScheme()
	if err := xpv1.AddToScheme(s); err != nil {
		fmt.Printf("ERROR: %v\n", err)

		panic(err)
	}

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)

		panic(err)
	}

	resourceScheme := schema.GroupVersionResource{
		Group:    "ec2.aws.crossplane.io",
		Version:  "v1beta1",
		Resource: "vpcs",
	}

	var vpcList ec2.VPCList

	resp, err := dynClient.Resource(resourceScheme).List(ctx, v1.ListOptions{})

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)

		panic(err)
	}

	unstructured := resp.UnstructuredContent()
	err = runtime.DefaultUnstructuredConverter.
		FromUnstructured(unstructured, &vpcList)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)

		panic(err)
	}

	for _, vpc := range vpcList.Items {
		err := wait.PollImmediate(3*time.Second, 10*time.Minute, func() (done bool, err error) {

			for _, condition := range vpc.Status.Conditions {
				if condition.Status != "True" {
					return false, nil
				}
			}

			return true, nil
		})
		if err != nil {
			fmt.Printf("wait.PollImmediate error waiting for condition to be true: %v\n", err)

			panic(err)
		}

	}

	fmt.Println("All Status=True")
}
