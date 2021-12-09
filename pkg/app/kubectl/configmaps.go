package kubectl

import (
	"encoding/json"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"

	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func ListConfigMaps() (*ConfigMaps, error) {
	out, err := bash.RunAndLogRead("kubectl get configmap --all-namespaces -o=json")
	if err != nil {
		return nil, err
	}

	var configMaps *ConfigMaps
	if err := json.Unmarshal(out, &configMaps); err != nil {
		return nil, errors.WithStack(err)
	}

	return configMaps, nil
}

func SearchConfigMap(label string) (*ConfigMaps, error) {
	out, err := runAndLogRead("get", "ConfigMap", "-o=json", "--all-namespaces", "-l", label)
	if err != nil {
		return nil, err
	}

	var configMaps *ConfigMaps
	if err := json.Unmarshal(out, &configMaps); err != nil {
		return nil, errors.WithStack(err)
	}

	return configMaps, nil
}

func GetConfigMapKeyValue(ref *KeyRef, namespace string) (string, error) {
	out, err := bash.Run(fmt.Sprintf("kubectl get configmap --namespace %s %s -o jsonpath=\"{.data.%s}\" | base64 --decode", namespace, ref.Name, ref.Key))
	if err != nil {
		return "", err
	}

	return string(out), nil
}

type ConfigMaps struct {
	Items []*ConfigMap `json:"items"`
}

func (m *ConfigMaps) Names() []string {
	configMaps := m.Items
	names := make([]string, 0, len(configMaps))
	for _, release := range configMaps {
		names = append(names, release.Metadata.Namespace+"/"+release.Metadata.Name)
	}
	return names
}

func (m *ConfigMaps) Get(name string) *ConfigMap {
	for _, configMap := range m.Items {
		if configMap.Metadata.Namespace+"/"+configMap.Metadata.Name == name {
			return configMap
		}
	}
	return nil
}

func (m *ConfigMaps) SelectOne() (*ConfigMap, error) {
	names := m.Names()

	if len(names) == 1 {
		return m.Get(names[0]), nil
	}

	var selected string
	prompt := &survey.Select{
		Message: "Select the ConfigMap:",
		Options: names,
	}

	err := survey.AskOne(prompt, &selected, survey.WithPageSize(10), survey.WithKeepFilter(true))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return m.Get(selected), nil
}

func (m *ConfigMaps) SingleResult() (*ConfigMap, error) {
	if len(m.Items) == 0 {
		return nil, nil
	}
	if len(m.Items) == 1 {
		return m.Items[0], nil
	}
	return nil, errors.New("found more than one ConfigMap")
}

type ConfigMap struct {
	Data     map[string]string `json:"data"`
	Kind     string            `json:"kind"`
	Metadata struct {
		Labels    map[string]string `json:"labels"`
		Name      string            `json:"name"`
		Namespace string            `json:"namespace"`
	} `json:"metadata"`
}
