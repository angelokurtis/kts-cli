package kubectl

import (
	"encoding/json"
	"time"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"

	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func ListNamespaces() (*Namespaces, error) {
	out, err := bash.RunAndLogRead("kubectl get namespaces -o=json")
	if err != nil {
		return nil, err
	}

	var namespaces *Namespaces
	if err := json.Unmarshal(out, &namespaces); err != nil {
		return nil, errors.WithStack(err)
	}

	return namespaces, nil
}

type Namespaces struct {
	Items []*Namespace `json:"items"`
}

func (n *Namespaces) Names() []string {
	namespaces := n.Items
	names := make([]string, 0, len(namespaces))

	for _, ns := range namespaces {
		names = append(names, ns.Metadata.Name)
	}

	return names
}

func (n *Namespaces) Get(name string) *Namespace {
	for _, namespace := range n.Items {
		if namespace.Metadata.Name == name {
			return namespace
		}
	}

	return nil
}

func (n *Namespaces) SelectOne() (*Namespace, error) {
	names := n.Names()

	if len(names) == 1 {
		return n.Get(names[0]), nil
	}

	var selected string

	prompt := &survey.Select{
		Message: "Select the namespace:",
		Options: names,
	}

	err := survey.AskOne(prompt, &selected, survey.WithPageSize(10), survey.WithKeepFilter(true))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return n.Get(selected), nil
}

func (n *Namespaces) SelectMany() (*Namespaces, error) {
	if len(n.Items) == 0 {
		return &Namespaces{}, nil
	}

	names := n.Names()
	prompt := &survey.MultiSelect{
		Message: "Select the Namespaces:",
		Options: names,
	}

	var selects []string

	err := survey.AskOne(prompt, &selects, survey.WithPageSize(10), survey.WithKeepFilter(true))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	ns := make([]*Namespace, 0, len(selects))
	for _, name := range selects {
		ns = append(ns, n.Get(name))
	}

	return &Namespaces{Items: ns}, nil
}

type Namespace struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Metadata   struct {
		CreationTimestamp time.Time         `json:"creationTimestamp"`
		Name              string            `json:"name"`
		ResourceVersion   string            `json:"resourceVersion"`
		SelfLink          string            `json:"selfLink"`
		UID               string            `json:"uid"`
		Labels            map[string]string `json:"labels"`
	} `json:"metadata"`
	Spec struct {
		Finalizers []string `json:"finalizers"`
	} `json:"spec"`
	Status struct {
		Phase string `json:"phase"`
	} `json:"status"`
}
