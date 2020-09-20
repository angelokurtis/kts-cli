package google

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/pkg/app/gcloud"
	"github.com/angelokurtis/kts-cli/pkg/bash"
	"github.com/gookit/color"
	changeCase "github.com/ku/go-change-case"
	"strings"
)

type ContainerNodePool struct {
	Name     string `hcl:"name"`
	Cluster  string `hcl:"cluster"`
	Location string `hcl:"location"`
	Project  string `hcl:"project"`
}

func (c *ContainerNodePool) GetType() string {
	return "google_container_node_pool"
}

func (c *ContainerNodePool) GetID() string {
	return fmt.Sprintf("%s.%s", c.GetType(), c.GetName())
}

func (c *ContainerNodePool) GetName() string {
	return changeCase.Snake(c.Name)
}

func (c *ContainerNodePool) Import() error {
	cmd := fmt.Sprintf("terraform import %s %s/%s/%s/%s", c.GetID(), c.Project, c.Location, c.Cluster, c.Name)
	if _, err := bash.RunAndLogWrite(cmd); err != nil {
		if strings.Contains(err.Error(), "Resource already managed by Terraform") {
			color.Yellow.Printf("[WARN] the %s is already managed by Terraform\n", c.GetID())
		} else {
			return err
		}
	}
	return nil
}

func NewContainerNodePool() (*ContainerNodePool, error) {
	cluster, err := gcloud.SelectGKECluster()
	if err != nil {
		return nil, err
	}
	return &ContainerNodePool{
		Name:     cluster.Name,
		Project:  cluster.Project.ID,
		Location: cluster.Location,
	}, nil
}
