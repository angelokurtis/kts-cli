package kube

import (
	"github.com/pkg/errors"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
	extensionsv1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewRestConfig() (*rest.Config, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	overrides := &clientcmd.ConfigOverrides{}

	cfg, err := clientcmd.
		NewNonInteractiveDeferredLoadingClientConfig(rules, overrides).
		ClientConfig()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return cfg, err
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
