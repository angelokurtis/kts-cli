package kube

import (
	"flag"
	"path/filepath"

	"github.com/pkg/errors"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	extensionsv1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func NewRestConfig() (*rest.Config, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return config, nil
}

func NewExtensions(clientset *kubernetes.Clientset) extensionsv1beta1.ExtensionsV1beta1Interface {
	return clientset.ExtensionsV1beta1()
}

func NewDiscoveryClient(cfg *rest.Config) (*discovery.DiscoveryClient, error) {
	discvy, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return discvy, nil
}

func NewClientset(cfg *rest.Config) (*kubernetes.Clientset, error) {
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return clientset, nil
}
