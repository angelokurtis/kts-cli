package kubectl

import "encoding/json"

func ListAllIngresses() ([]*Ingress, error) {
	out, err := run("get", "ingress", "--all-namespaces", "-o=json", "--request-timeout=5s")
	if err != nil {
		return nil, err
	}

	var ingresses *Ingresses
	if err := json.Unmarshal(out, &ingresses); err != nil {
		return nil, err
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
		return nil, err
	}

	return ingresses, nil
}

type Ingresses struct {
	Items []*Ingress `json:"items"`
}

type Metadata struct {
	Labels    map[string]string `json:"labels"`
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
}

type Rule struct {
	Host string `json:"host"`
}

type Spec struct {
	Rules []*Rule `json:"rules"`
}

type IngressIP struct {
	IP string `json:"ip"`
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
