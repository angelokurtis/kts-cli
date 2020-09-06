package terraform

import (
	changeCase "github.com/ku/go-change-case"
	"github.com/rodaine/hclencoder"
)

type (
	Resource struct {
		ResourceDefinition `hcl:"resource"`
	}
	ResourceDefinition struct {
		Type         string `hcl:",key"`
		Name         string `hcl:",key"`
		ResourceSpec `hcl:",squash"`
	}
	ResourceSpec interface {
		GetType() string
		GetName() string
		GetID() string
		Import() error
	}
)

func newResource(t string) *Resource {
	return &Resource{ResourceDefinition: ResourceDefinition{Type: t}}
}

func (r *Resource) set(s ResourceSpec) {
	r.Name = changeCase.Snake(s.GetName())
	r.ResourceSpec = s
}

func (r *Resource) Encode() ([]byte, error) {
	hcl, err := hclencoder.Encode(r)
	if err != nil {
		return nil, err
	}
	return hcl, nil
}
