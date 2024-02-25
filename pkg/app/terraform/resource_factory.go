package terraform

import (
	"github.com/pkg/errors"

	"github.com/angelokurtis/kts-cli/pkg/app/terraform/resources/google"
	"github.com/angelokurtis/kts-cli/pkg/app/terraform/resources/helm"
	"github.com/angelokurtis/kts-cli/pkg/app/terraform/resources/kubernetes"
)

func BuildResource(t string) (*Resource, error) {
	switch t {
	case "google_container_cluster":
		cluster, err := google.NewContainerCluster()
		if err != nil {
			return nil, err
		}

		return NewResource(cluster), nil
	case "google_container_node_pool":
		cluster, err := google.NewContainerNodePool()
		if err != nil {
			return nil, err
		}

		return NewResource(cluster), nil
	case "helm_release":
		release, err := helm.NewRelease()
		if err != nil {
			return nil, err
		}

		return NewResource(release), nil
	case "kubernetes_namespace":
		namespace, err := kubernetes.NewNamespace()
		if err != nil {
			return nil, err
		}

		return NewResource(namespace), nil
	case "kubernetes_config_map":
		namespace, err := kubernetes.NewConfigMap()
		if err != nil {
			return nil, err
		}

		return NewResource(namespace), nil
	case "kubernetes_deployment":
		namespace, err := kubernetes.NewDeployment()
		if err != nil {
			return nil, err
		}

		return NewResource(namespace), nil
	}

	return nil, errors.New("resource type not found")
}
