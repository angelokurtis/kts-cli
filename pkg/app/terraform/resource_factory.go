package terraform

import (
	"github.com/angelokurtis/kts-cli/pkg/app/terraform/resources/google"
	"github.com/pkg/errors"
)

func NewResource(t string) (*Resource, error) {
	resource := newResource(t)
	switch t {
	case "google_container_cluster":
		cluster, err := google.NewContainerCluster()
		if err != nil {
			return nil, err
		}
		resource.set(cluster)
		return resource, nil
	}
	return nil, errors.New("resource type not found")
}
