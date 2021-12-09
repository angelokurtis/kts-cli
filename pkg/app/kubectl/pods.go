package kubectl

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/enescakir/emoji"
	"github.com/pkg/errors"
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

	err := survey.AskOne(prompt, &selected, survey.WithPageSize(10), survey.WithKeepFilter(true))
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

	err := survey.AskOne(prompt, &selected, survey.WithPageSize(10), survey.WithKeepFilter(true))
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

	err := survey.AskOne(prompt, &selects, survey.WithPageSize(10), survey.WithKeepFilter(true))
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

	err := survey.AskOne(prompt, &selected, survey.WithPageSize(10), survey.WithKeepFilter(true))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return pods[selected], nil
}

// type AutoGenerated struct {
//	Annotations struct {
//		KubernetesIoPsp      string `json:"kubernetes.io/psp"`
//		PrometheusIoPath     string `json:"prometheus.io/path"`
//		PrometheusIoPort     string `json:"prometheus.io/port"`
//		PrometheusIoScrape   string `json:"prometheus.io/scrape"`
//		SidecarIstioIoStatus string `json:"sidecar.istio.io/status"`
//	} `json:"annotations"`
//	CreationTimestamp time.Time `json:"creationTimestamp"`
//	GenerateName      string    `json:"generateName"`
//	Labels            struct {
//		ControllerUID                   string `json:"controller-uid"`
//		IstioIoRev                      string `json:"istio.io/rev"`
//		JobName                         string `json:"job-name"`
//		SecurityIstioIoTLSMode          string `json:"security.istio.io/tlsMode"`
//		ServiceIstioIoCanonicalName     string `json:"service.istio.io/canonical-name"`
//		ServiceIstioIoCanonicalRevision string `json:"service.istio.io/canonical-revision"`
//	} `json:"labels"`
//	Name            string `json:"name"`
//	Namespace       string `json:"namespace"`
//	OwnerReferences []struct {
//		APIVersion         string `json:"apiVersion"`
//		BlockOwnerDeletion bool   `json:"blockOwnerDeletion"`
//		Controller         bool   `json:"controller"`
//		Kind               string `json:"kind"`
//		Name               string `json:"name"`
//		UID                string `json:"uid"`
//	} `json:"ownerReferences"`
//	ResourceVersion string `json:"resourceVersion"`
//	SelfLink        string `json:"selfLink"`
//	UID             string `json:"uid"`
// }

type (
	Pod struct {
		Metadata struct {
			Labels            map[string]string `json:"labels"`
			Name              string            `json:"name"`
			Namespace         string            `json:"namespace"`
			OwnerReferences   []*OwnerReference `json:"ownerReferences"`
			CreationTimestamp time.Time         `json:"creationTimestamp"`
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
	OwnerReference struct {
		APIVersion         string `json:"apiVersion"`
		BlockOwnerDeletion bool   `json:"blockOwnerDeletion"`
		Controller         bool   `json:"controller"`
		Kind               string `json:"kind"`
		Name               string `json:"name"`
		UID                string `json:"uid"`
	}
	ContainerStatus struct {
		ContainerID  string          `json:"containerID"`
		Image        string          `json:"image"`
		ImageID      string          `json:"imageID"`
		Name         string          `json:"name"`
		Ready        bool            `json:"ready"`
		RestartCount int             `json:"restartCount"`
		Started      bool            `json:"started"`
		LastState    *ContainerState `json:"lastState"`
		State        *ContainerState `json:"state"`
	}
	ContainerState struct {
		Running    *ContainerStateRunning    `json:"running"`
		Terminated *ContainerStateTerminated `json:"terminated"`
		Waiting    *ContainerStateWaiting    `json:"waiting"`
	}
)

func (s *ContainerState) Event() ContainerStateEvent {
	if s.Running != nil {
		return s.Running
	} else if s.Terminated != nil {
		return s.Terminated
	} else if s.Waiting != nil {
		return s.Waiting
	}
	return nil
}

func (cs *ContainerStatus) CurrentState() ContainerStateEvent {
	states := make([]*ContainerState, 0, 0)
	if cs.LastState != nil {
		states = append(states, cs.LastState)
	}
	if cs.State != nil {
		states = append(states, cs.State)
	}
	sort.Slice(states, func(i, j int) bool {
		future := time.Unix(1<<63-1, 0)

		ie := states[i].Event()
		it := &future
		if ie != nil && ie.GetTime() != nil {
			it = ie.GetTime()
		}

		je := states[j].Event()
		jt := &future
		if je != nil && je.GetTime() != nil {
			jt = je.GetTime()
		}

		return it.Before(*jt)
	})
	state := states[len(states)-1]
	event := state.Event()
	return event
}

func (p *Pod) RestartCount() int {
	restart := 0
	for _, container := range p.Status.ContainerStatuses {
		if container.RestartCount > 0 {
			restart = container.RestartCount
		}
	}
	return restart
}

func (p *Pod) CurrentStatus() string {
	for _, container := range p.Status.ContainerStatuses {
		s := container.CurrentState()
		if _, ok := s.(*ContainerStateWaiting); ok || s.GetReason() == "Error" {
			return s.GetReason()
		}
	}
	if p.Status.Phase == "Succeeded" {
		return "Completed"
	}
	return p.Status.Phase
}

func (p *Pod) LastState() ContainerStateEvent {
	states := make([]ContainerStateEvent, 0, 0)
	for _, container := range p.Status.ContainerStatuses {
		state := container.CurrentState()
		if state != nil {
			states = append(states, state)
		}
	}
	sort.Slice(states, func(i, j int) bool {
		t1 := states[i].GetTime()
		t2 := states[j].GetTime()
		if t1 == nil {
			t1 = &time.Time{}
		}
		if t2 == nil {
			t2 = &time.Time{}
		}
		return t1.Before(*t2)
	})
	if len(states) == 0 {
		return nil
	}
	return states[len(states)-1]
}

func (p *Pod) StatusColor() string {
	containers := p.Status.ContainerStatuses
	restart := 0
	ready := 0
	desired := len(p.Spec.Containers)
	for _, container := range containers {
		state := container.CurrentState()
		if state.GetReason() == "Error" {
			return emoji.RedCircle.String()
		}
		if _, ok := state.(*ContainerStateRunning); ok {
			ready++
		}
		if container.RestartCount > 0 {
			restart = container.RestartCount
		}
	}
	if p.CurrentStatus() == "Completed" {
		return emoji.BlackCircle.String()
	}
	if ready < desired {
		return emoji.YellowCircle.String()
	}
	if restart > 0 {
		return emoji.YellowCircle.String()
	}
	return emoji.GreenCircle.String()
}

func (p *Pod) Ready() string {
	containers := p.Status.ContainerStatuses
	running := 0
	for _, container := range containers {
		state := container.CurrentState()
		if _, ok := state.(*ContainerStateRunning); ok {
			running++
		}
	}
	return fmt.Sprintf("%d/%d", running, len(p.Spec.Containers))
}

func (p *Pod) LastUpdate() *time.Time {
	state := p.LastState()
	if state == nil {
		return &p.Metadata.CreationTimestamp
	}
	t := state.GetTime()
	if t == nil {
		return &p.Metadata.CreationTimestamp
	}
	return t
}

func (p *Pod) EnvironmentVariables() []*Env {
	envvars := make([]*Env, 0, 0)
	for _, container := range p.Spec.Containers {
		envvars = append(envvars, container.Env...)
	}
	return envvars
}

func (p *Pod) IsJob() bool {
	owners := p.Metadata.OwnerReferences
	if owners == nil || len(owners) == 0 {
		return false
	}
	for _, owner := range owners {
		if owner.Kind == "Job" {
			return true
		}
	}
	return false
}

func (p *Pod) HasIstioSidecar() bool {
	containers := p.Spec.Containers
	if len(containers) < 2 {
		// log.Debugf("%s/%s has not sidecar", d.Metadata.Namespace, d.Metadata.Name)
		return false
	}
	for _, container := range containers {
		if container.Name == "istio-proxy" {
			// log.Debugf("%s/%s has the sidecar", d.Metadata.Namespace, d.Metadata.Name)
			return true
		}
	}
	// log.Debugf("%s/%s has not sidecar", d.Metadata.Namespace, d.Metadata.Name)
	return false
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

	err := survey.AskOne(prompt, &selected, survey.WithPageSize(10), survey.WithKeepFilter(true))
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return strconv.Atoi(selected)
}
