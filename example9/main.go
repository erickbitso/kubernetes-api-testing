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

	// Load the provider-aws types into the scheme
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

	// Now we can use the client to issue lists

	objects := []interface{}{
		&ec2.VPCList{},
		&ec2.SubnetList{},
		&ec2.InternetGatewayList{},
		&ec2.NATGatewayList{},
		&ec2.RouteTableList{},
	}

	for _, obj := range objects {
		x := obj.(client.ObjectList)
		getObjectList(c, x)
	}
}

func getObjectList(client client.Client, list client.ObjectList) {
	fmt.Println("check")
	err := client.List(context.Background(), list)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)

		panic(err)
	}
	kind := list.GetObjectKind().GroupVersionKind().Kind
	fmt.Printf("%s: %v\n", kind, list)

}
