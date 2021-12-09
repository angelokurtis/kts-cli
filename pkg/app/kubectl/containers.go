package kubectl

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
)

func ListContainersByDeployment(deploy *Deployment) (*Containers, error) {
	labels := make([]string, 0, 0)
	for key, value := range deploy.Spec.Selector.MatchLabels {
		labels = append(labels, key+"="+value)
	}
	cmd := []string{"get", "pods", "-o=json", "-n", deploy.Metadata.Namespace, "-l", strings.Join(labels, ", ")}
	out, err := run(cmd...)
	if err != nil {
		return nil, err
	}

	var pods *Pods
	if err := json.Unmarshal(out, &pods); err != nil {
		return nil, errors.WithStack(err)
	}

	return pods.Containers(), nil
}

func ListContainers(namespace string, allNamespaces bool, selector string) (*Containers, error) {
	pods, err := ListPods(namespace, allNamespaces, selector)
	if err != nil {
		return nil, err
	}
	return pods.Containers(), nil
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
		Args            []string `json:"args"`
		Command         []string `json:"command"`
		Env             []*Env   `json:"env"`
		Image           string   `json:"image"`
		ImagePullPolicy string   `json:"imagePullPolicy"`
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
		LivenessProbe struct {
			FailureThreshold    *int `json:"failureThreshold"`
			InitialDelaySeconds *int `json:"initialDelaySeconds"`
			PeriodSeconds       *int `json:"periodSeconds"`
			SuccessThreshold    *int `json:"successThreshold"`
			TimeoutSeconds      *int `json:"timeoutSeconds"`
		} `json:"livenessProbe"`
		ReadinessProbe struct {
			FailureThreshold    *int `json:"failureThreshold"`
			InitialDelaySeconds *int `json:"initialDelaySeconds"`
			PeriodSeconds       *int `json:"periodSeconds"`
			SuccessThreshold    *int `json:"successThreshold"`
			TimeoutSeconds      *int `json:"timeoutSeconds"`
		} `json:"readinessProbe"`
		Status *ContainerStatus
	}
	Env struct {
		Name      string     `json:"name,omitempty"`
		Value     string     `json:"value,omitempty"`
		ValueFrom *ValueFrom `json:"valueFrom,omitempty"`
	}
	ValueFrom struct {
		ConfigMapKeyRef *KeyRef `json:"configMapKeyRef,omitempty"`
		SecretKeyRef    *KeyRef `json:"secretKeyRef,omitempty"`
	}
	KeyRef struct {
		Key      string `json:"key,omitempty"`
		Name     string `json:"name,omitempty"`
		Optional bool   `json:"optional,omitempty"`
	}
)

func (c *Container) LastUpdateTime() *time.Time {
	status := c.Status
	if status == nil {
		return nil
	}
	state := status.CurrentState()
	if state == nil {
		return nil
	}
	return state.GetTime()
}

func (c *Container) GetState() ContainerStateEvent {
	if c.Status == nil {
		return nil
	}
	return c.Status.CurrentState()
}

func (c *Container) GetEnv(key string) string {
	for _, v := range c.Env {
		if key == v.Name {
			return v.Value
		}
	}
	return ""
}

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
		// p = dedupeStr(p, container.Pod)
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

	err := survey.AskOne(prompt, &selects, survey.WithPageSize(10), survey.WithKeepFilter(true))
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

// func (c *Containers) SelectOnePort() *Container {
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
//	err := survey.AskOne(prompt, &selected, survey.WithPageSize(10), survey.WithKeepFilter(true))
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//
//	return containers[selected], nil
// }
