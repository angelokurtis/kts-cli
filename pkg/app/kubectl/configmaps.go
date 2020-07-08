package kubectl

import (
	"encoding/json"
	"errors"
)

func SearchConfigMap(label string) (*ConfigMaps, error) {
	out, err := runAndLog("get", "ConfigMap", "-o=json", "--all-namespaces", "-l", label)
	if err != nil {
		return nil, err
	}

	var configMaps *ConfigMaps
	if err := json.Unmarshal(out, &configMaps); err != nil {
		return nil, err
	}

	return configMaps, nil
}

type ConfigMaps struct {
	Items []*ConfigMap `json:"items"`
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
