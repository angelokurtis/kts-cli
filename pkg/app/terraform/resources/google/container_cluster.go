package google

import (
	"fmt"
	"strings"

	"github.com/gookit/color"
	changeCase "github.com/ku/go-change-case"

	"github.com/angelokurtis/kts-cli/pkg/app/gcloud"
	"github.com/angelokurtis/kts-cli/pkg/bash"
)

type ContainerCluster struct {
	Location string `hcl:"location"`
	Name     string `hcl:"name"`
	Project  string `hcl:"project"`
}

func (c *ContainerCluster) GetType() string {
	return "google_container_cluster"
}

func (c *ContainerCluster) GetID() string {
	return fmt.Sprintf("%s.%s", c.GetType(), c.GetName())
}

func (c *ContainerCluster) GetName() string {
	return changeCase.Snake(c.Name)
}

func (c *ContainerCluster) Import() error {
	cmd := fmt.Sprintf("terraform import %s projects/%s/locations/%s/clusters/%s", c.GetID(), c.Project, c.Location, c.Name)
	if _, err := bash.RunAndLogWrite(cmd); err != nil {
		if strings.Contains(err.Error(), "Resource already managed by Terraform") {
			color.Yellow.Printf("[WARN] the %s is already managed by Terraform\n", c.GetID())
		} else {
			return err
		}
	}

	return nil
}

func NewContainerCluster() (*ContainerCluster, error) {
	cluster, err := gcloud.SelectGKECluster()
	if err != nil {
		return nil, err
	}

	return &ContainerCluster{
		Name:     cluster.Name,
		Project:  cluster.Project.ID,
		Location: cluster.Location,
	}, nil
}
