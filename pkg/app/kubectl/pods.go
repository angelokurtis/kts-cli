package kubectl

import (
	"encoding/json"
	"github.com/AlecAivazis/survey/v2"
	"strconv"
	"strings"
)

func ListAllPods() (*Pods, error) {
	out, err := run("get", "pods", "--all-namespaces", "-o=json")
	if err != nil {
		return nil, err
	}

	var pods *Pods
	if err := json.Unmarshal(out, &pods); err != nil {
		return nil, err
	}

	return pods, nil
}

type Pods struct {
	Items []*Pod `json:"items"`
}

func (s *Pods) Labels() map[string][]string {
	labels := make(map[string][]string, 0)
	for _, pod := range s.Items {
		for k, v := range pod.Metadata.Labels {
			values := labels[k]
			values = dedupeStr(values, v)
			labels[k] = values
		}
	}
	return labels
}

func (s *Pods) SelectLabels() (map[string][]string, error) {
	var options []string
	for key, values := range s.Labels() {
		for _, value := range values {
			options = append(options, key+"="+value)
		}
	}

	var selects []string
	prompt := &survey.MultiSelect{
		Message: "Select the pod labels:",
		Options: options,
	}

	err := survey.AskOne(prompt, &selects, survey.WithPageSize(10))
	if err != nil {
		return nil, err
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

func (s *Pods) Namespaces(labels map[string][]string) []string {
	m := make(map[string][]string)
	for _, service := range s.Items {
		for k, v := range service.Metadata.Labels {
			label := k + "=" + v
			ns := m[label]
			ns = append(ns, service.Metadata.Namespace)
			m[label] = ns
		}
	}
	namespaces := make([]string, 0)
	for k, values := range labels {
		for _, v := range values {
			label := k + "=" + v
			ns := m[label]
			namespaces = dedupeStr(namespaces, ns...)
		}
	}
	return namespaces
}

func (s *Pods) SelectNamespace(labels map[string][]string) (string, error) {
	namespaces := s.Namespaces(labels)

	if len(namespaces) == 1 {
		return namespaces[0], nil
	}

	var selected string
	prompt := &survey.Select{
		Message: "Select the namespace:",
		Options: namespaces,
	}

	err := survey.AskOne(prompt, &selected, survey.WithPageSize(10))
	if err != nil {
		return "", err
	}

	return selected, nil
}

func (s *Pods) SelectContainerPort(namespace string, labels map[string][]string) (int, error) {
	pods := s.Pods(namespace, labels)
	ports := make([]string, 0, 0)
	for _, pod := range pods {
		for _, container := range pod.Spec.Containers {
			for _, port := range container.Ports {
				ports = dedupeStr(ports, strconv.Itoa(port.ContainerPort))
			}
		}
	}

	if len(ports) == 1 {
		return strconv.Atoi(ports[0])
	}

	var selected string
	prompt := &survey.Select{
		Message: "Select the container port:",
		Options: ports,
	}

	err := survey.AskOne(prompt, &selected, survey.WithPageSize(10))
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(selected)
}

func (s *Pods) Pods(namespace string, labels map[string][]string) []*Pod {
	pods := make([]*Pod, 0, 0)
	for _, pod := range s.Items {
		if pod.Metadata.Namespace != namespace {
			continue
		}
		if !pod.ContainsLabels(labels) {
			continue
		}
		pods = append(pods, pod)
	}
	return pods
}

func (s *Pods) SelectMany() ([]*Pod, error) {
	pods := make(map[string]*Pod, 0)
	names := make([]string, 0, 0)
	for _, pod := range s.Items {
		name := pod.Metadata.Namespace + "/" + pod.Metadata.Name
		pods[name] = pod
		names = append(names, name)
	}

	var selects []string
	prompt := &survey.MultiSelect{
		Message: "Select the pods:",
		Options: names,
	}

	err := survey.AskOne(prompt, &selects, survey.WithPageSize(10))
	if err != nil {
		return nil, err
	}

	result := make([]*Pod, 0, 0)
	for _, s := range selects {
		result = append(result, pods[s])
	}

	return result, nil
}

func (s *Pods) SelectOne() (*Pod, error) {
	pods := make(map[string]*Pod, 0)
	names := make([]string, 0, 0)
	for _, pod := range s.Items {
		name := pod.Metadata.Namespace + "/" + pod.Metadata.Name
		pods[name] = pod
		names = append(names, name)
	}

	if len(s.Items) == 1 {
		return s.Items[0], nil
	}

	var selected string
	prompt := &survey.Select{
		Message: "Select the pod:",
		Options: names,
	}

	err := survey.AskOne(prompt, &selected, survey.WithPageSize(10))
	if err != nil {
		return nil, err
	}

	return pods[selected], nil
}

type Containers struct {
	Name  string `json:"name"`
	Ports []*struct {
		ContainerPort int    `json:"containerPort"`
		Name          string `json:"name"`
		Protocol      string `json:"protocol"`
	} `json:"ports"`
}

type Pod struct {
	Metadata struct {
		Labels    map[string]string `json:"labels"`
		Name      string            `json:"name"`
		Namespace string            `json:"namespace"`
	} `json:"metadata"`
	Spec struct {
		Containers []*Containers `json:"containers"`
	} `json:"spec"`
}

func (p *Pod) ContainsLabels(labels map[string][]string) bool {
	l := make([]string, 0, 0)
	for k, v := range p.Metadata.Labels {
		l = append(l, k+"="+v)
	}
	s := make([]string, 0, 0)
	for k, values := range labels {
		for _, v := range values {
			s = append(s, k+"="+v)
		}
	}
	for _, askLabel := range s {
		for _, foundLabel := range l {
			if foundLabel == askLabel {
				return true
			}
		}
	}
	return false
}

func (p *Pod) SelectContainerPort() (int, error) {
	ports := make([]string, 0, 0)
	for _, container := range p.Spec.Containers {
		for _, port := range container.Ports {
			ports = dedupeStr(ports, strconv.Itoa(port.ContainerPort))
		}
	}

	if len(ports) == 1 {
		return strconv.Atoi(ports[0])
	}

	var selected string
	prompt := &survey.Select{
		Message: "Select the container port:",
		Options: ports,
	}

	err := survey.AskOne(prompt, &selected, survey.WithPageSize(10))
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(selected)
}