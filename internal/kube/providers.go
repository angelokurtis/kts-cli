package kube

import (
	"flag"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/angelokurtis/kts-cli/internal/system"
)

var (
	svcOnce         sync.Once
	clientset       *kubernetes.Clientset
	discoveryClient discovery.DiscoveryInterface
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

func NewClientSet(config *rest.Config) (*kubernetes.Clientset, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return clientset, nil
}

func GetClientset() *kubernetes.Clientset {
	svcOnce.Do(func() {
		config, err := NewRestConfig()
		if err != nil {
			system.Exit(err)
		}

		clientset, err = NewClientSet(config)
		if err != nil {
			system.Exit(err)
		}
	})

	return clientset
}

func NewDiscoveryClient(config *rest.Config) discovery.DiscoveryInterface {
	client := discovery.NewDiscoveryClientForConfigOrDie(config)
	return client
}

func GetDiscoveryClient() discovery.DiscoveryInterface {
	svcOnce.Do(func() {
		config, err := NewRestConfig()
		if err != nil {
			system.Exit(err)
		}

		discoveryClient = NewDiscoveryClient(config)
	})

	return discoveryClient
}
