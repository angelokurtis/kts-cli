package terraform

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"github.com/rodaine/hclencoder"
)

func ListProviders() []string {
	return []string{
		"aws",
		"azurerm",
		"google",
		"helm",
		"kubernetes",
		"vault",
	}
}

func SelectProvider() (*Provider, error) {
	providers := ListProviders()

	var selected string
	if len(providers) == 0 {
		return nil, nil
	} else if len(providers) > 1 {
		prompt := &survey.Select{
			Message: "Select the Terraform Provider:",
			Options: providers,
		}

		err := survey.AskOne(prompt, &selected, survey.WithPageSize(10))
		if err != nil {
			return nil, errors.WithStack(err)
		}
	} else {
		selected = providers[0]
	}

	return NewProvider(selected), nil
}

type (
	Provider struct {
		ProviderDefinition `hcl:"provider"`
	}
	ProviderDefinition struct {
		Name string      `hcl:",key"`
		Spec interface{} `hcl:",squash"`
	}
)

func NewProvider(name string) *Provider {
	return &Provider{ProviderDefinition: ProviderDefinition{Name: name}}
}

func (p *Provider) Encode() ([]byte, error) {
	hcl, err := hclencoder.Encode(p)
	if err != nil {
		return nil, err
	}
	return hcl, nil
}
