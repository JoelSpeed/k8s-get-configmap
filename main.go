package main

import (
	"flag"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	// Import K8s auth providers
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
)

func getClient(kubeconfig string) (kubernetes.Interface, error) {
	var config *rest.Config
	if kubeconfig != "" {
		var err error
		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, err
		}
	} else {
		var err error
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func main() {
	// get a kubeconfig variable from a flag (if not set expecting to be running within a Pod)
	kubeconfig := flag.String("kubeconfig", "", "kubeconfig file to use (uses current context)")
	flag.Parse()

	client, err := getClient(*kubeconfig)
	if err != nil {
		panic(err)
	}

	// client.<APIGroupVersion>().<Kind>(<Namespace>).Get(<Name>, metav1.GetOptions{})
	cm, err := client.CoreV1().ConfigMaps("kube-system").Get("extension-apiserver-authentication", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", cm)
}
