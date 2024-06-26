package kube

import (
	"github.com/google/wire"
)

var Providers = wire.NewSet(
	NewDiscoveryClient,
	NewClientset,
	NewExtensions,
	NewRestConfig,
)
