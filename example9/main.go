package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"time"

	"github.com/crossplane-contrib/provider-aws/apis"
	ec2 "github.com/crossplane-contrib/provider-aws/apis/ec2/v1beta1"
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
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

	objects := []client.ObjectList{
		&ec2.VPCList{},
		&ec2.SubnetList{},
		&ec2.InternetGatewayList{},
		&ec2.NATGatewayList{},
		&ec2.RouteTableList{},
	}

	for _, obj := range objects {
		verifyObjectList(c, obj)
		//getObjectList(c, obj)
	}

}

func verifyObjectList(client client.Client, list client.ObjectList) {
	err := client.List(context.Background(), list)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)

		panic(err)
	}

	//x := list.(*ec2.VPCList)
	switch v := list.(type) {
	case *ec2.VPCList:
		fmt.Println("VPC: ")
		for _, item := range v.Items {
			_, err := resourceStatus(item.Status.ResourceStatus)
			processError(err)

		}

	case *ec2.SubnetList:
		fmt.Println("Subnets: ")
		for _, item := range v.Items {
			_, err := resourceStatus(item.Status.ResourceStatus)
			processError(err)
		}
	case *ec2.NATGatewayList:
		fmt.Println("NAT Gateways: ")
		for _, item := range v.Items {
			_, err := resourceStatus(item.Status.ResourceStatus)
			processError(err)
		}
	case *ec2.AddressList:
		fmt.Println("Elastic IPs: ")
		for _, item := range v.Items {
			_, err := resourceStatus(item.Status.ResourceStatus)
			processError(err)
		}
	case *ec2.InternetGatewayList:
		fmt.Println("Internet Gateways: ")
		for _, item := range v.Items {
			_, err := resourceStatus(item.Status.ResourceStatus)
			processError(err)
		}
	case *ec2.RouteTableList:
		fmt.Println("Route Tables: ")
		for _, item := range v.Items {
			_, err := resourceStatus(item.Status.ResourceStatus)
			processError(err)
		}
	default:
		fmt.Println("not valid")
	}

}

func resourceStatus(item xpv1.ResourceStatus) (bool, error) {
	err := wait.PollImmediate(3*time.Second, 15*time.Second, func() (done bool, err error) {

		for _, condition := range item.Conditions {
			if condition.Status != "True" {
				return false, nil
			}
		}
		fmt.Println("Status=True")
		return true, nil
	})
	if err != nil {
		fmt.Printf("wait.PollImmediate error waiting for condition to be true: %v\n", err)

		return false, err
	}
	return true, nil
}

func processError(err error) {
	if err != nil {
		fmt.Printf("Error: %v\n", err)

		panic(err)
	}
}
