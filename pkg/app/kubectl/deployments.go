package kubectl

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/enescakir/emoji"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func ListDeployments(namespace string, allNamespaces bool) (*Deployments, error) {
	cmd := []string{"get", "deployments", "-o=json"}
	if allNamespaces {
		cmd = append(cmd, "--all-namespaces")
	} else if namespace != "" {
		cmd = append(cmd, "-n", namespace)
	}

	out, err := run(cmd...)
	if err != nil {
		return nil, err
	}

	var deploys *Deployments
	if err := json.Unmarshal(out, &deploys); err != nil {
		return nil, errors.WithStack(err)
	}

	return deploys, nil
}

func ListAllDeployments() (*Deployments, error) {
	out, err := bash.RunAndLogRead("kubectl get deployments --all-namespaces -o=json")
	if err != nil {
		return nil, err
	}

	var deployments *Deployments
	if err := json.Unmarshal(out, &deployments); err != nil {
		return nil, errors.WithStack(err)
	}

	return deployments, nil
}

type Deployments struct {
	Items []*Deployment `json:"items"`
}

func (d *Deployments) FilterInjected() *Deployments {
	deployments := make([]*Deployment, 0, 0)

	for _, deployment := range d.Items {
		if deployment.HasIstioSidecar() {
			deployments = append(deployments, deployment)
		}
	}

	return &Deployments{Items: deployments}
}

func (d *Deployments) FilterUninjected() *Deployments {
	deployments := make([]*Deployment, 0, 0)

	for _, deployment := range d.Items {
		if !deployment.HasIstioSidecar() {
			deployments = append(deployments, deployment)
		}
	}

	return &Deployments{Items: deployments}
}

func (d *Deployments) FilterInjectable() *Deployments {
	deployments := make([]*Deployment, 0, 0)

	for _, deployment := range d.Items {
		if deployment.IsInjectable() {
			deployments = append(deployments, deployment)
		}
	}

	return &Deployments{Items: deployments}
}

func (d *Deployments) Names() []string {
	deployments := d.Items
	names := make([]string, 0, len(deployments))

	for _, release := range deployments {
		names = append(names, release.Metadata.Name)
	}

	return names
}

func (d *Deployments) FullNames() []string {
	deployments := d.Items
	names := make([]string, 0, len(deployments))

	for _, release := range deployments {
		names = append(names, release.Metadata.Namespace+"/"+release.Metadata.Name)
	}

	return names
}

func (d *Deployments) Namespaces() []string {
	deployments := d.Items
	namespaces := make([]string, 0, len(deployments))

	for _, release := range deployments {
		namespaces = dedupeStr(namespaces, release.Metadata.Namespace)
	}

	return namespaces
}

func (d *Deployments) Get(name string) *Deployment {
	for _, deployment := range d.Items {
		if deployment.Metadata.Namespace+"/"+deployment.Metadata.Name == name {
			return deployment
		}
	}

	return nil
}

func (d *Deployments) SelectOne() (*Deployment, error) {
	names := d.FullNames()

	if len(names) == 1 {
		return d.Get(names[0]), nil
	}

	var selected string

	prompt := &survey.Select{
		Message: "Select the Deployment:",
		Options: names,
	}

	err := survey.AskOne(prompt, &selected, survey.WithPageSize(10), survey.WithKeepFilter(true))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return d.Get(selected), nil
}

func (d *Deployments) SelectMany() (*Deployments, error) {
	if len(d.Items) == 0 {
		return &Deployments{}, nil
	}

	names := d.FullNames()
	prompt := &survey.MultiSelect{
		Message: "Select the Deployments:",
		Options: names,
	}

	var selects []string

	err := survey.AskOne(prompt, &selects, survey.WithPageSize(10), survey.WithKeepFilter(true))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	deploys := make([]*Deployment, 0, len(selects))
	for _, name := range selects {
		deploys = append(deploys, d.Get(name))
	}

	return &Deployments{Items: deploys}, nil
}

func (d *Deployments) SelectContainers() (*Containers, error) {
	if len(d.Items) == 0 {
		return &Containers{}, nil
	}

	containers, err := d.ListContainers()
	if err != nil {
		return nil, err
	}

	return containers.SelectMany()
}

func (d *Deployments) ListContainers() (*Containers, error) {
	if len(d.Items) == 0 {
		return &Containers{}, nil
	}

	c := make([]*Container, 0, 0)

	for _, deploy := range d.Items {
		containers := deploy.Spec.Template.Spec.Containers
		c = append(c, containers...)
	}

	return &Containers{Items: c}, nil
}

func (d *Deployments) Rollout() error {
	byNamespace := lo.GroupBy(d.Items, func(item *Deployment) string {
		return item.Metadata.Namespace
	})

	for namespace, deployments := range byNamespace {
		names := lo.Map(deployments, func(deployment *Deployment, index int) string {
			return deployment.Metadata.Name
		})

		_, err := bash.RunAndLogWrite(fmt.Sprintf("kubectl rollout restart deployment %s -n %s", strings.Join(names, " "), namespace))
		if err != nil {
			return err
		}
	}

	return nil
}

type (
	Deployment struct {
		APIVersion string `json:"apiVersion"`
		Kind       string `json:"kind"`
		Metadata   struct {
			Annotations struct {
				DeploymentKubernetesIoRevision              string `json:"deployment.kubernetes.io/revision"`
				KubectlKubernetesIoLastAppliedConfiguration string `json:"kubectl.kubernetes.io/last-applied-configuration"`
			} `json:"annotations"`
			CreationTimestamp time.Time `json:"creationTimestamp"`
			Generation        int       `json:"generation"`
			Labels            struct {
				App     string `json:"app"`
				Release string `json:"release"`
				Version string `json:"version"`
			} `json:"labels"`
			Name            string `json:"name"`
			Namespace       string `json:"namespace"`
			ResourceVersion string `json:"resourceVersion"`
			SelfLink        string `json:"selfLink"`
			UID             string `json:"uid"`
		} `json:"metadata"`
		Spec struct {
			Replicas int `json:"replicas"`
			Selector struct {
				MatchLabels map[string]string `json:"matchLabels"`
			} `json:"selector"`
			Template struct {
				Metadata struct {
					Annotations map[string]string `json:"annotations"`
					Labels      map[string]string `json:"labels"`
				} `json:"metadata"`
				Spec struct {
					Containers []*Container `json:"containers"`
				} `json:"spec"`
			} `json:"template"`
		} `json:"spec"`
		Status struct {
			AvailableReplicas  int                `json:"availableReplicas"`
			Conditions         []*StatusCondition `json:"conditions"`
			ObservedGeneration int                `json:"observedGeneration"`
			ReadyReplicas      int                `json:"readyReplicas"`
			Replicas           int                `json:"replicas"`
			UpdatedReplicas    int                `json:"updatedReplicas"`
		} `json:"status"`
	}
	StatusCondition struct {
		LastTransitionTime time.Time `json:"lastTransitionTime"`
		LastUpdateTime     time.Time `json:"lastUpdateTime"`
		Message            string    `json:"message"`
		Reason             string    `json:"reason"`
		Status             string    `json:"status"`
		Type               string    `json:"type"`
	}
)

func (d *Deployment) LastUpdateTime() *time.Time {
	conditions := d.Status.Conditions
	if len(conditions) < 1 {
		return nil
	}

	sort.Slice(conditions, func(i, j int) bool {
		return conditions[i].LastTransitionTime.After(conditions[j].LastTransitionTime)
	})

	return &conditions[0].LastUpdateTime
}

func (d *Deployment) StatusColor() string {
	ready := d.Status.ReadyReplicas
	desired := d.Spec.Replicas

	if ready >= desired {
		return emoji.GreenCircle.String()
	}

	if ready > 0 {
		return emoji.YellowCircle.String()
	}

	return emoji.RedCircle.String()
}

func (d *Deployment) HasIstioSidecar() bool {
	containers := d.Spec.Template.Spec.Containers
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

func (d *Deployment) GetContainer(name string) *Container {
	containers := d.Spec.Template.Spec.Containers
	for _, container := range containers {
		if container.Name == name {
			return container
		}
	}

	return nil
}

func (d *Deployment) IsInjectable() bool {
	annotations := d.Spec.Template.Metadata.Annotations
	inject := annotations["sidecar.istio.io/inject"]
	// log.Debugf("deployment %s/%s should inject? %s", d.Metadata.Namespace, d.Metadata.Name, inject)
	if inject == "false" {
		return false
	}

	return true
}
