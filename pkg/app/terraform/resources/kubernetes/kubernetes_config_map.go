package kubernetes

import (
	"fmt"
	"strings"

	"github.com/gookit/color"
	changeCase "github.com/ku/go-change-case"
	"github.com/pkg/errors"

	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/angelokurtis/kts-cli/pkg/bash"
)

type (
	ConfigMap struct {
		Metadata Metadata `hcl:"metadata"`
	}
)

func (c *ConfigMap) GetType() string {
	return "kubernetes_config_map"
}

func (c *ConfigMap) GetID() string {
	return fmt.Sprintf("%s.%s", c.GetType(), c.GetName())
}

func (c *ConfigMap) GetName() string {
	return changeCase.Snake(c.Metadata.Name)
}

func (c *ConfigMap) Import() error {
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

func NewConfigMap() (*ConfigMap, error) {
	configMaps, err := kubectl.ListConfigMaps()
	if err != nil {
		return nil, err
	}

	configMap, err := configMaps.SelectOne()
	if err != nil {
		return nil, err
	} else if configMap == nil {
		return nil, errors.New("configmap was not found")
	}

	return &ConfigMap{Metadata: Metadata{Name: configMap.Metadata.Name, Namespace: configMap.Metadata.Namespace}}, nil
}
