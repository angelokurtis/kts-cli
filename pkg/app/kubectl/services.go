package kubectl

import (
	"encoding/json"
	"github.com/AlecAivazis/survey/v2"
	"sort"
	"strings"
)

func ListAllServices() (*Services, error) {
	out, err := run("get", "services", "--all-namespaces", "-o=json")
	if err != nil {
		return nil, err
	}

	var services *Services
	if err := json.Unmarshal(out, &services); err != nil {
		return nil, err
	}

	return services, nil
}

type Services struct {
	Items []*Service `json:"items"`
}

func (s *Services) SelectLabels() (map[string][]string, error) {
	var options []string
	for key, values := range s.Labels() {
		for _, value := range values {
			options = append(options, key+"="+value)
		}
	}

	var selects []string
	prompt := &survey.MultiSelect{
		Message: "Select the service labels:",
		Options: options,
	}

	err := survey.AskOne(prompt, &selects, survey.WithPageSize(10))
	if err != nil {
		return nil, err
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

func (s *Services) Namespaces(labels map[string][]string) []string {
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
			namespaces = dedupe(namespaces, ns...)
		}
	}
	return namespaces
}

func (s *Services) Labels() map[string][]string {
	labels := make(map[string][]string, 0)
	for _, service := range s.Items {
		for k, v := range service.Metadata.Labels {
			values := labels[k]
			values = dedupe(values, v)
			labels[k] = values
		}
	}
	return labels
}

func (s *Services) LabelKeys() []string {
	keys := make([]string, 0)
	for _, service := range s.Items {
		for k := range service.Metadata.Labels {
			keys = dedupe(keys, k)
		}
	}
	sort.Strings(keys)
	return keys
}

type Metadata struct {
	Labels    map[string]string `json:"labels"`
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
}

type Port struct {
	Name     string `json:"name"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol"`
}

type Spec struct {
	Ports []Port `json:"ports"`
}

type Service struct {
	Metadata Metadata `json:"metadata"`
	Spec     Spec     `json:"spec"`
}

func dedupe(a []string, b ...string) []string {

	check := make(map[string]int)
	d := append(a, b...)
	res := make([]string, 0)
	for _, val := range d {
		check[val] = 1
	}

	for letter, _ := range check {
		res = append(res, letter)
	}

	return res
}
