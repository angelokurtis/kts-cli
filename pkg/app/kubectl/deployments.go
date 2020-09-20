package kubectl

import (
	"encoding/json"
	"github.com/AlecAivazis/survey/v2"
	"github.com/angelokurtis/kts-cli/pkg/bash"
	"github.com/pkg/errors"
	"time"
)

func ListDeployments() (*Deployments, error) {
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

func (d *Deployments) Names() []string {
	deployments := d.Items
	names := make([]string, 0, len(deployments))
	for _, release := range deployments {
		names = append(names, release.Metadata.Namespace+"/"+release.Metadata.Name)
	}
	return names
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
	names := d.Names()

	if len(names) == 1 {
		return d.Get(names[0]), nil
	}

	var selected string
	prompt := &survey.Select{
		Message: "Select the Deployment:",
		Options: names,
	}

	err := survey.AskOne(prompt, &selected, survey.WithPageSize(10))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return d.Get(selected), nil
}

type Deployment struct {
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
	Status struct {
		AvailableReplicas int `json:"availableReplicas"`
		Conditions        []struct {
			LastTransitionTime time.Time `json:"lastTransitionTime"`
			LastUpdateTime     time.Time `json:"lastUpdateTime"`
			Message            string    `json:"message"`
			Reason             string    `json:"reason"`
			Status             string    `json:"status"`
			Type               string    `json:"type"`
		} `json:"conditions"`
		ObservedGeneration int `json:"observedGeneration"`
		ReadyReplicas      int `json:"readyReplicas"`
		Replicas           int `json:"replicas"`
		UpdatedReplicas    int `json:"updatedReplicas"`
	} `json:"status"`
}
