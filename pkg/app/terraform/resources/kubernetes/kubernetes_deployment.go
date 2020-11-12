package kubernetes

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/angelokurtis/kts-cli/pkg/bash"
	"github.com/gookit/color"
	changeCase "github.com/ku/go-change-case"
	"github.com/pkg/errors"
	"strings"
)

type (
	Deployment struct {
		Metadata Metadata `hcl:"metadata"`
	}
)

func (c *Deployment) GetType() string {
	return "kubernetes_deployment"
}

func (c *Deployment) GetID() string {
	return fmt.Sprintf("%s.%s", c.GetType(), c.GetName())
}

func (c *Deployment) GetName() string {
	return changeCase.Snake(c.Metadata.Name)
}

func (c *Deployment) Import() error {
	cmd := fmt.Sprintf("terraform import %s %s/%s", c.GetID(), c.Metadata.Namespace, c.Metadata.Name)
	if _, err := bash.RunAndLogWrite(cmd); err != nil {
		if strings.Contains(err.Error(), "Resource already managed by Terraform") {
			color.Yellow.Printf("[WARN] the %s is already managed by Terraform\n", c.GetID())
		} else {
			return err
		}
	}
	return nil
}

func NewDeployment() (*Deployment, error) {
	deployments, err := kubectl.ListAllDeployments()
	if err != nil {
		return nil, err
	}
	deployment, err := deployments.SelectOne()
	if err != nil {
		return nil, err
	} else if deployment == nil {
		return nil, errors.New("Deployment was not found")
	}
	return &Deployment{Metadata: Metadata{Name: deployment.Metadata.Name, Namespace: deployment.Metadata.Namespace}}, nil
}
