package kubernetes

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/angelokurtis/kts-cli/pkg/bash"
	"github.com/gookit/color"
	changeCase "github.com/ku/go-change-case"
	"strings"
)

type (
	Namespace struct {
		Metadata MetadataClusterScoped `hcl:"metadata"`
	}
	MetadataClusterScoped struct {
		Name string `hcl:"name"`
	}
)

func (n *Namespace) GetType() string {
	return "kubernetes_namespace"
}

func (n *Namespace) GetID() string {
	return fmt.Sprintf("%s.%s", n.GetType(), n.GetName())
}

func (n *Namespace) GetName() string {
	return changeCase.Snake(n.Metadata.Name)
}

func (n *Namespace) Import() error {
	cmd := fmt.Sprintf("terraform import %s %s", n.GetID(), n.Metadata.Name)
	if _, err := bash.RunAndLogWrite(cmd); err != nil {
		if strings.Contains(err.Error(), "Resource already managed by Terraform") {
			color.Yellow.Printf("[WARN] the %s is already managed by Terraform\n", n.GetID())
		} else {
			return err
		}
	}
	return nil
}

func NewNamespace() (*Namespace, error) {
	namespaces, err := kubectl.ListNamespaces()
	if err != nil {
		return nil, err
	}
	namespace, err := namespaces.SelectOne()
	if err != nil {
		return nil, err
	}
	return &Namespace{Metadata: MetadataClusterScoped{Name: namespace.Metadata.Name}}, nil
}
