package terraform

import (
	"github.com/pkg/errors"
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

func NewResource(s ResourceSpec) *Resource {
	return &Resource{ResourceDefinition: ResourceDefinition{
		Type:         s.GetType(),
		Name:         s.GetName(),
		ResourceSpec: s,
	}}
}

func (r *Resource) Encode() ([]byte, error) {
	hcl, err := hclencoder.Encode(r)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return hcl, nil
}
