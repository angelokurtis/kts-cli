package kubectl

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"strings"
)

func ListContainers(namespace string, allNamespaces bool) (*Containers, error) {
	c := make([]*Container, 0, 0)
	pods, err := ListPods(namespace, allNamespaces)
	if err != nil {
		return nil, err
	}
	for _, pod := range pods.Items {
		podName := pod.Metadata.Name
		podNamespace := pod.Metadata.Namespace
		podTemplateHash := pod.Metadata.Labels["pod-template-hash"]
		for _, container := range pod.Spec.Containers {
			container.Pod = podName
			container.Namespace = podNamespace
			container.PodTemplateHash = podTemplateHash
			c = append(c, container)
		}
		for _, container := range pod.Spec.InitContainers {
			container.Pod = podName
			container.Namespace = podNamespace
			container.PodTemplateHash = podTemplateHash
			c = append(c, container)
		}
	}
	containers := &Containers{Items: c}
	for _, container := range containers.Items {
		container.Single = containers.CountByPod(container.Pod) == 1
	}
	return containers, nil
}

type (
	Containers struct {
		Items []*Container
	}
	Container struct {
		Namespace       string
		Pod             string
		PodTemplateHash string
		Single          bool
		Name            string `json:"name"`
		Ports           []*struct {
			ContainerPort int    `json:"containerPort"`
			Name          string `json:"name"`
			Protocol      string `json:"protocol"`
		} `json:"ports"`
		Args    []string `json:"args"`
		Command []string `json:"command"`
		Env     []*struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"env"`
		Image           string `json:"image"`
		ImagePullPolicy string `json:"imagePullPolicy"`
		Resources       struct {
			Limits struct {
				CPU    string `json:"cpu"`
				Memory string `json:"memory"`
			} `json:"limits"`
			Requests struct {
				CPU    string `json:"cpu"`
				Memory string `json:"memory"`
			} `json:"requests"`
		} `json:"resources"`
		SecurityContext struct {
			AllowPrivilegeEscalation bool `json:"allowPrivilegeEscalation"`
			Capabilities             struct {
				Drop []string `json:"drop"`
			} `json:"capabilities"`
			ReadOnlyRootFilesystem bool `json:"readOnlyRootFilesystem"`
			RunAsGroup             int  `json:"runAsGroup"`
			RunAsNonRoot           bool `json:"runAsNonRoot"`
			RunAsUser              int  `json:"runAsUser"`
		} `json:"securityContext"`
		TerminationMessagePath   string `json:"terminationMessagePath"`
		TerminationMessagePolicy string `json:"terminationMessagePolicy"`
		VolumeMounts             []*struct {
			MountPath string `json:"mountPath"`
			Name      string `json:"name"`
			ReadOnly  bool   `json:"readOnly,omitempty"`
		} `json:"volumeMounts"`
	}
)

func (c *Containers) Contains(containerName string) bool {
	for _, container := range c.Items {
		if container.Name == containerName {
			return true
		}
	}
	return false
}

func (c *Containers) Namespaces() []string {
	ns := make([]string, 0, 0)
	for _, container := range c.Items {
		ns = dedupeStr(ns, container.Namespace)
	}
	return ns
}

func (c *Containers) Names() []string {
	n := make([]string, 0, 0)
	for _, container := range c.Items {
		n = dedupeStr(n, container.Name)
	}
	return n
}

func (c *Containers) Pods() []string {
	p := make([]string, 0, 0)
	for _, container := range c.Items {
		prefix := strings.Split(container.Pod, "-"+container.PodTemplateHash)[0]
		p = dedupeStr(p, prefix)
		//p = dedupeStr(p, container.Pod)
	}
	return p
}

func (c *Containers) CountByPod(pod string) int {
	count := 0
	for _, container := range c.Items {
		if container.Pod == pod {
			count++
		}
	}
	return count
}

func (c *Containers) SelectMany() (*Containers, error) {
	containers := make(map[string]*Container, 0)
	names := make([]string, 0, 0)
	for _, container := range c.Items {
		name := ""
		if container.Namespace != "" {
			name = container.Namespace + "/"
		}
		if container.Pod != "" {
			name = container.Pod + "/"
		}
		name = name + container.Name
		containers[name] = container
		names = dedupeStr(names, name)
	}

	var selects []string
	prompt := &survey.MultiSelect{
		Message: "Select the containers:",
		Options: names,
	}

	err := survey.AskOne(prompt, &selects, survey.WithPageSize(10))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := make([]*Container, 0, 0)
	for _, s := range selects {
		result = append(result, containers[s])
	}

	return &Containers{Items: result}, nil
}

func (c *Containers) FilterExposed() *Containers {
	containers := make([]*Container, 0, 0)
	for _, container := range c.Items {
		if len(container.Ports) > 0 {
			containers = append(containers, container)
		}
	}
	return &Containers{Items: containers}
}

//func (c *Containers) SelectOnePort() *Container {
//	containers := make(map[string]*Container, 0)
//	names := make([]string, 0, 0)
//	for _, container := range c.Items {
//		for _, port := range container.Ports {
//			name := container.Namespace + "/" + container.Pod + "/" + container.Name
//			containers[name] = container
//			names = append(names, name)
//		}
//	}
//
//	if len(s.Items) == 1 {
//		return s.Items[0], nil
//	}
//
//	var selected string
//	prompt := &survey.Select{
//		Message: "Select the container:",
//		Options: names,
//	}
//
//	err := survey.AskOne(prompt, &selected, survey.WithPageSize(10))
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//
//	return containers[selected], nil
//}
