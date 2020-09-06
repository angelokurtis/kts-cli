package kubectl

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
)

func ListContainers(namespace string, allNamespaces bool) (*Containers, error) {
	containers := make([]*Container, 0, 0)
	pods, err := ListPods(namespace, allNamespaces)
	if err != nil {
		return nil, err
	}
	for _, pod := range pods.Items {
		for _, container := range pod.Spec.Containers {
			container.Pod = pod.Metadata.Name
			container.Namespace = pod.Metadata.Namespace
			containers = append(containers, container)
		}
		for _, container := range pod.Spec.InitContainers {
			container.Pod = pod.Metadata.Name
			container.Namespace = pod.Metadata.Namespace
			containers = append(containers, container)
		}
	}
	return &Containers{Items: containers}, nil
}

type (
	Containers struct {
		Items []*Container
	}
	Container struct {
		Namespace string
		Pod       string
		Name      string `json:"name"`
		Ports     []*struct {
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
		name := container.Namespace + "/" + container.Pod + "/" + container.Name
		containers[name] = container
		names = append(names, name)
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
