package kubectl

import (
	"encoding/json"
	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

func ListPods(namespace string, allNamespaces bool, selector string) (*Pods, error) {
	cmd := []string{"get", "pods", "-o=json"}
	if allNamespaces {
		cmd = append(cmd, "--all-namespaces")
	} else if namespace != "" {
		cmd = append(cmd, "-n", namespace)
	}
	if selector != "" {
		cmd = append(cmd, "-l", selector)
	}
	out, err := run(cmd...)
	if err != nil {
		return nil, err
	}

	var pods *Pods
	if err := json.Unmarshal(out, &pods); err != nil {
		return nil, errors.WithStack(err)
	}

	return pods, nil
}

func ListAllPods() (*Pods, error) {
	out, err := run("get", "pods", "--all-namespaces", "-o=json")
	if err != nil {
		return nil, err
	}

	var pods *Pods
	if err := json.Unmarshal(out, &pods); err != nil {
		return nil, errors.WithStack(err)
	}

	return pods, nil
}

type Pods struct {
	Items []*Pod `json:"items"`
}

func (s *Pods) Containers() *Containers {
	c := make([]*Container, 0, 0)
	for _, pod := range s.Items {
		podName := pod.Metadata.Name
		podNamespace := pod.Metadata.Namespace
		podTemplateHash := pod.Metadata.Labels["pod-template-hash"]
		for _, container := range pod.Spec.Containers {
			container.Pod = podName
			container.Namespace = podNamespace
			container.PodTemplateHash = podTemplateHash
			container.Status = pod.GetContainerStatus(container.Name)
			c = append(c, container)
		}
		for _, container := range pod.Spec.InitContainers {
			container.Pod = podName
			container.Namespace = podNamespace
			container.PodTemplateHash = podTemplateHash
			container.Status = pod.GetContainerStatus(container.Name)
			c = append(c, container)
		}
	}
	containers := &Containers{Items: c}
	for _, container := range containers.Items {
		container.Single = containers.CountByPod(container.Pod) == 1
	}
	return containers
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
		return "", errors.WithStack(err)
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
		return 0, errors.WithStack(err)
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

func (s *Pods) SelectMany() (*Pods, error) {
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
		return nil, errors.WithStack(err)
	}

	result := make([]*Pod, 0, 0)
	for _, s := range selects {
		result = append(result, pods[s])
	}

	return &Pods{result}, nil
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
		return nil, errors.WithStack(err)
	}

	return pods[selected], nil
}

type Pod struct {
	Metadata struct {
		Labels    map[string]string `json:"labels"`
		Name      string            `json:"name"`
		Namespace string            `json:"namespace"`
	} `json:"metadata"`
	Spec struct {
		Containers     []*Container `json:"containers"`
		InitContainers []*Container `json:"initContainers"`
	} `json:"spec"`
	Status struct {
		Conditions []struct {
			LastTransitionTime time.Time `json:"lastTransitionTime"`
			Status             string    `json:"status"`
			Type               string    `json:"type"`
		} `json:"conditions"`
		ContainerStatuses     []*ContainerStatus `json:"containerStatuses"`
		HostIP                string             `json:"hostIP"`
		InitContainerStatuses []*ContainerStatus `json:"initContainerStatuses"`
		Phase                 string             `json:"phase"`
		PodIP                 string             `json:"podIP"`
		PodIPs                []struct {
			IP string `json:"ip"`
		} `json:"podIPs"`
		QosClass  string    `json:"qosClass"`
		StartTime time.Time `json:"startTime"`
	} `json:"status"`
}

type ContainerStatus struct {
	ContainerID  string `json:"containerID"`
	Image        string `json:"image"`
	ImageID      string `json:"imageID"`
	Name         string `json:"name"`
	Ready        bool   `json:"ready"`
	RestartCount int    `json:"restartCount"`
	Started      bool   `json:"started"`
	LastState    struct {
		Running    *ContainerStateRunning    `json:"running"`
		Terminated *ContainerStateTerminated `json:"terminated"`
		Waiting    *ContainerStateWaiting    `json:"waiting"`
	} `json:"lastState"`
	State struct {
		Running    *ContainerStateRunning    `json:"running"`
		Terminated *ContainerStateTerminated `json:"terminated"`
		Waiting    *ContainerStateWaiting    `json:"waiting"`
	} `json:"state"`
}

func (cs *ContainerStatus) GetLastState() ContainerState {
	if cs.LastState.Running != nil {
		return cs.LastState.Running
	} else if cs.LastState.Terminated != nil {
		return cs.LastState.Terminated
	} else if cs.LastState.Waiting != nil {
		return cs.LastState.Waiting
	}
	return nil
}

func (cs *ContainerStatus) GetState() ContainerState {
	if cs.State.Running != nil {
		return cs.State.Running
	} else if cs.State.Terminated != nil {
		return cs.State.Terminated
	} else if cs.State.Waiting != nil {
		return cs.State.Waiting
	}
	return nil
}

func (p *Pod) GetContainerStatus(container string) *ContainerStatus {
	for _, status := range p.Status.ContainerStatuses {
		if status.Name == container {
			return status
		}
	}
	for _, status := range p.Status.InitContainerStatuses {
		if status.Name == container {
			return status
		}
	}
	return nil
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
		return 0, errors.WithStack(err)
	}

	return strconv.Atoi(selected)
}
