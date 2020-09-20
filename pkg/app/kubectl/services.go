package kubectl

import (
	"encoding/json"
	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"sort"
	"strings"
)

func ListServices() (*Services, error) {
	out, err := run("get", "services", "--all-namespaces", "-o=json")
	if err != nil {
		return nil, err
	}

	var services *Services
	if err := json.Unmarshal(out, &services); err != nil {
		return nil, errors.WithStack(err)
	}

	return services, nil
}

type Services struct {
	Items []*Service `json:"items"`
}

func (s *Services) FilterByType(t string) *Services {
	items := make([]*Service, 0, 0)
	for _, item := range s.Items {
		if item.Spec.Type == t {
			items = append(items, item)
		}
	}
	return &Services{Items: items}
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
			namespaces = dedupeStr(namespaces, ns...)
		}
	}
	return namespaces
}

func (s *Services) Labels() map[string][]string {
	labels := make(map[string][]string, 0)
	for _, service := range s.Items {
		for k, v := range service.Metadata.Labels {
			values := labels[k]
			values = dedupeStr(values, v)
			labels[k] = values
		}
	}
	return labels
}

func (s *Services) LabelKeys() []string {
	keys := make([]string, 0)
	for _, service := range s.Items {
		for k := range service.Metadata.Labels {
			keys = dedupeStr(keys, k)
		}
	}
	sort.Strings(keys)
	return keys
}

type Service struct {
	Metadata struct {
		Labels    map[string]string `json:"labels"`
		Name      string            `json:"name"`
		Namespace string            `json:"namespace"`
	} `json:"metadata"`
	Spec struct {
		Ports []struct {
			NodePort int    `json:"nodePort"`
			Port     int    `json:"port"`
			Protocol string `json:"protocol"`
		} `json:"ports"`
		Selector map[string]string `json:"selector"`
		Type     string            `json:"type"`
	} `json:"spec"`
	Status struct {
		LoadBalancer struct {
			Ingresses []struct {
				Hostname string `json:"hostname"`
				IP       string `json:"ip"`
			} `json:"ingress"`
		} `json:"loadBalancer"`
	} `json:"status"`
}

func (s *Service) ExternalIP() string {
	for _, ingress := range s.Status.LoadBalancer.Ingresses {
		if ingress.Hostname != "" {
			return ingress.Hostname
		}
		if ingress.IP != "" {
			return ingress.IP
		}
	}
	return ""
}
