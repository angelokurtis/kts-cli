package kubectl

import (
	"encoding/json"
	"github.com/pkg/errors"
)

func ListIngresses() ([]*Ingress, error) {
	out, err := run("get", "ingress", "--all-namespaces", "-o=json", "--request-timeout=5s")
	if err != nil {
		return nil, err
	}

	var ingresses *Ingresses
	if err := json.Unmarshal(out, &ingresses); err != nil {
		return nil, errors.WithStack(err)
	}

	return ingresses.Items, nil
}

func SearchIngress(label string) (*Ingresses, error) {
	out, err := run("get", "ingress", "--all-namespaces", "-l", label, "-o=json")
	if err != nil {
		return nil, err
	}

	var ingresses *Ingresses
	if err := json.Unmarshal(out, &ingresses); err != nil {
		return nil, errors.WithStack(err)
	}

	return ingresses, nil
}

type Ingresses struct {
	Items []*Ingress `json:"items"`
}

type Metadata struct {
	Labels          map[string]string `json:"labels"`
	Name            string            `json:"name"`
	Namespace       string            `json:"namespace"`
	Generation      int               `json:"generation"`
	ResourceVersion string            `json:"resourceVersion"`
	UID             string            `json:"uid"`
}

type Rule struct {
	Host string `json:"host"`
}

type Spec struct {
	Rules []*Rule `json:"rules"`
}

type IngressIP struct {
	Hostname string `json:"hostname"`
	IP       string `json:"ip"`
}

type LoadBalancer struct {
	Ingresses []*IngressIP `json:"ingress"`
}

type Status struct {
	LoadBalancer LoadBalancer `json:"loadBalancer"`
}

type Ingress struct {
	Metadata Metadata `json:"metadata"`
	Spec     Spec     `json:"spec"`
	Status   Status   `json:"status"`
}

func (i *Ingress) ExternalIP() string {
	for _, ingress := range i.Status.LoadBalancer.Ingresses {
		if ingress.Hostname != "" {
			return ingress.Hostname
		}
		if ingress.IP != "" {
			return ingress.IP
		}
	}
	return ""
}
