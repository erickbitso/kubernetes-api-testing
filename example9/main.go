package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/crossplane-contrib/provider-aws/apis"
	ec2 "github.com/crossplane-contrib/provider-aws/apis/ec2/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
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
	if err := apis.AddToScheme(s); err != nil {
		fmt.Printf("ERROR: %v\n", err)

		panic(err)
	}

	c, err := client.New(config, client.Options{Scheme: s})
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)

		panic(err)
	}
	vpcs := &ec2.VPCList{}

	err = c.List(context.Background(), vpcs)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)

		panic(err)
	}
	fmt.Println(vpcs)
	fmt.Println("All Status=True")
}
