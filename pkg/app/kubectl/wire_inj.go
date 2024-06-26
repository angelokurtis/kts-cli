//go:build wireinject
// +build wireinject

package kubectl

import (
	"github.com/google/wire"
	"k8s.io/client-go/discovery"
	extensionsv1beta1 "k8s.io/client-go/kubernetes/typed/extensions/v1beta1"

	"github.com/angelokurtis/kts-cli/internal/kube"
)

func newDiscovery() (discovery.DiscoveryInterface, error) {
	wire.Build(
		kube.Providers,
		wire.Bind(new(discovery.DiscoveryInterface), new(*discovery.DiscoveryClient)),
	)

	return nil, nil
}

func newExtensions() (extensionsv1beta1.ExtensionsV1beta1Interface, error) {
	wire.Build(kube.Providers)

	return nil, nil
}
