package kiali

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
)

type (
	Nodes map[string]*Node
	Node  struct {
		Workload string `json:"workload"`
		App      string `json:"app"`
		Version  string `json:"version"`
		Traffic  []struct {
			Protocol string `json:"protocol"`
			Rates    struct {
				HTTPIn  string `json:"httpIn"`
				HTTPOut string `json:"httpOut"`
				TCPOut  string `json:"tcpOut"`
			} `json:"rates,omitempty"`
		} `json:"traffic"`
		ID           string `json:"id"`
		NodeType     string `json:"nodeType"`
		Namespace    string `json:"namespace"`
		Service      string `json:"service"`
		DestServices []struct {
			Namespace string `json:"namespace"`
			Name      string `json:"name"`
		} `json:"destServices"`
		IsServiceEntry string `json:"isServiceEntry"`
	}
)

func (n *Node) GetName() string {
	name := ""
	switch n.NodeType {
	case "workload":
		name = n.Workload
	case "service":
		name = n.Service
	}
	return n.Namespace + "/" + n.NodeType + "/" + name
}

func (n *Node) Selector() string {
	label := ""
	name := ""
	switch n.NodeType {
	case "workload", "unknown":
		label = "app"
		name = n.App
	case "service":
		label = "svc"
		name = n.Service
	}
	return label + " != " + name
}

func (n Nodes) Join(o Nodes) Nodes {
	res := make(Nodes, 0)
	for k, v := range n {
		res[k] = v
	}
	for k, v := range o {
		res[k] = v
	}
	return res
}

func (n Nodes) Get(name string) *Node {
	for _, node := range n {
		if name == node.GetName() {
			return node
		}
	}
	return nil
}

func (n Nodes) FullNames() []string {
	names := make([]string, 0, len(n))
	for _, node := range n {
		names = append(names, node.GetName())
	}
	return names
}

func (n Nodes) SelectOne() (*Node, error) {
	names := n.FullNames()

	if len(names) == 1 {
		return n.Get(names[0]), nil
	}

	var selected string
	prompt := &survey.Select{
		Message: "Select the Node:",
		Options: names,
	}

	err := survey.AskOne(prompt, &selected, survey.WithPageSize(10))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return n.Get(selected), nil
}
