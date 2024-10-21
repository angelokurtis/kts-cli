package kubectl

import (
	"encoding/json"
	"strings"
	"time"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
)

func ListNodes(selector string) (*Nodes, error) {
	cmd := []string{"get", "nodes", "-o=json"}

	if selector != "" {
		cmd = append(cmd, "-l", selector)
	}

	out, err := run(cmd...)
	if err != nil {
		return nil, err
	}

	var nodes *Nodes
	if err = json.Unmarshal(out, &nodes); err != nil {
		return nil, errors.WithStack(err)
	}

	return nodes, nil
}

type Node struct {
	APIVersion string       `json:"apiVersion"`
	Kind       string       `json:"kind"`
	Metadata   NodeMetadata `json:"metadata"`
	Spec       NodeSpec     `json:"spec"`
	Status     NodeStatus   `json:"status"`
}

type NodeMetadata struct {
	Annotations       Annotations `json:"annotations"`
	CreationTimestamp time.Time   `json:"creationTimestamp"`
	Labels            Labels      `json:"labels"`
	Name              string      `json:"name"`
	ResourceVersion   string      `json:"resourceVersion"`
	Uid               string      `json:"uid"`
}

type Annotations map[string]string

type Labels map[string]string

type NodeSpec struct {
	PodCIDR    string   `json:"podCIDR"`
	PodCIDRs   []string `json:"podCIDRs"`
	ProviderID string   `json:"providerID"`
	Taints     []Taint  `json:"taints"`
}

type Taint struct {
	Effect string `json:"effect"`
	Key    string `json:"key"`
}

type NodeStatus struct {
	Addresses       []Address       `json:"addresses"`
	Allocatable     Allocatable     `json:"allocatable"`
	Capacity        Allocatable     `json:"capacity"`
	Conditions      []Condition     `json:"conditions"`
	DaemonEndpoints DaemonEndpoints `json:"daemonEndpoints"`
	Images          []Image         `json:"images"`
	NodeInfo        NodeInfo        `json:"nodeInfo"`
}

type Address struct {
	Address string `json:"address"`
	Type    string `json:"type"`
}

type Allocatable struct {
	CPU              string `json:"cpu"`
	EphemeralStorage string `json:"ephemeral-storage"`
	Hugepages1Gi     string `json:"hugepages-1Gi"`
	Hugepages2Mi     string `json:"hugepages-2Mi"`
	Memory           string `json:"memory"`
	Pods             string `json:"pods"`
}

type Condition struct {
	LastHeartbeatTime  time.Time `json:"lastHeartbeatTime"`
	LastTransitionTime time.Time `json:"lastTransitionTime"`
	Message            string    `json:"message"`
	Reason             string    `json:"reason"`
	Status             string    `json:"status"`
	Type               string    `json:"type"`
}

type DaemonEndpoints struct {
	KubeletEndpoint KubeletEndpoint `json:"kubeletEndpoint"`
}

type KubeletEndpoint struct {
	Port int64 `json:"Port"`
}

type Image struct {
	Names     []string `json:"names"`
	SizeBytes int64    `json:"sizeBytes"`
}

type NodeInfo struct {
	Architecture            string `json:"architecture"`
	BootID                  string `json:"bootID"`
	ContainerRuntimeVersion string `json:"containerRuntimeVersion"`
	KernelVersion           string `json:"kernelVersion"`
	KubeProxyVersion        string `json:"kubeProxyVersion"`
	KubeletVersion          string `json:"kubeletVersion"`
	MachineID               string `json:"machineID"`
	OperatingSystem         string `json:"operatingSystem"`
	OSImage                 string `json:"osImage"`
	SystemUUID              string `json:"systemUUID"`
}

type Nodes struct {
	Items []*Node `json:"items"`
}

func (n *Nodes) Labels() map[string][]string {
	labels := make(map[string][]string, 0)

	for _, pod := range n.Items {
		for k, v := range pod.Metadata.Labels {
			values := labels[k]
			values = dedupeStr(values, v)
			labels[k] = values
		}
	}

	return labels
}

func (n *Nodes) SelectLabels() (map[string][]string, error) {
	var options []string

	for key, values := range n.Labels() {
		for _, value := range values {
			options = append(options, key+"="+value)
		}
	}

	var selects []string

	prompt := &survey.MultiSelect{
		Message: "Select the node labels:",
		Options: options,
	}

	err := survey.AskOne(prompt, &selects, survey.WithPageSize(10), survey.WithKeepFilter(true))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	labels := make(map[string][]string, 0)

	for _, s := range selects {
		spt := strings.Split(s, "=")
		key := spt[0]
		value := spt[len(spt)-1]

		values := labels[key]
		values = append(values, value)
		labels[key] = values
	}

	return labels, nil
}
