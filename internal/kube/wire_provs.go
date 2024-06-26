package kube

import (
	"github.com/google/wire"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"
)

var Providers = wire.NewSet(
	discovery.NewDiscoveryClientForConfig,
	kubernetes.NewForConfig,
	NewExtensions,
	NewRestConfig,
)
