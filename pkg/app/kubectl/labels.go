package kubectl

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/gookit/color"
	"github.com/pkg/errors"

	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func SelectLabel(resources string, namespace string, allNamespaces bool) (string, error) {
	labels, err := ListLabels(resources, namespace, allNamespaces)
	if err != nil {
		return "", err
	}

	if len(labels) == 1 {
		return labels[0], nil
	}

	var selected string
	prompt := &survey.Select{
		Message: "Select the label:",
		Options: labels,
	}

	err = survey.AskOne(prompt, &selected, survey.WithPageSize(10), survey.WithKeepFilter(true))
	if err != nil {
		return "", errors.WithStack(err)
	}

	return selected, nil
}

func ListLabels(resources string, namespace string, allNamespaces bool) ([]string, error) {
	cmd := []string{"kubectl", "get", resources}
	if allNamespaces {
		cmd = append(cmd, "--all-namespaces")
	} else if namespace != "" {
		cmd = append(cmd, "-n", namespace)
	}
	color.Secondary.Println(strings.Join(cmd, " ") + " --show-labels")
	out, err := bash.Run(strings.Join(cmd, " ") + " -o json")
	if err != nil {
		return nil, err
	}
	var r *GenericResources
	if err := json.Unmarshal(out, &r); err != nil {
		return nil, errors.WithStack(err)
	}

	labels := make([]string, 0, 0)
	for _, item := range r.Items {
		for k, v := range item.Metadata.Labels {
			labels = dedupeStr(labels, fmt.Sprintf("%s=%s", k, v))
		}
	}
	sort.Strings(labels)
	return labels, nil
}

func RemoveLabels(resources string, label string, namespace string, allNamespaces bool) error {
	cmd := []string{"kubectl", "get", resources, fmt.Sprintf("-l \"%s\"", label)}
	if allNamespaces {
		cmd = append(cmd, "--all-namespaces")
	} else if namespace != "" {
		cmd = append(cmd, "-n", namespace)
	}
	out, err := bash.Run(strings.Join(cmd, " ") + " -o json")
	if err != nil {
		return err
	}
	var r *GenericResources
	if err := json.Unmarshal(out, &r); err != nil {
		return errors.WithStack(err)
	}
	key := strings.Split(label, "=")[0]
	for _, item := range r.Items {
		if item.Metadata.Namespace == "" {
			color.Secondary.Printf("kubectl label %s %s %s-\n", item.Kind, item.Metadata.Name, key)
		} else {
			color.Secondary.Printf("kubectl label %s %s %s- -n %s\n", item.Kind, item.Metadata.Name, key, item.Metadata.Namespace)
		}
	}
	return nil
}

type GenericResources struct {
	Items []*GenericResource `json:"items"`
}

type GenericResource struct {
	Kind     string `json:"kind"`
	Metadata struct {
		Labels    map[string]string `json:"labels"`
		Name      string            `json:"name"`
		Namespace string            `json:"namespace"`
	} `json:"metadata"`
}
