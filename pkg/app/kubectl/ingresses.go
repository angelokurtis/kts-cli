package kubectl

import "encoding/json"

func SearchIngress(label string) (*Ingresses, error) {
	out, err := run("get", "ingress", "--all-namespaces", "-l", label)
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
	Items []Ingress `json:"items"`
}

type Metadata struct {
	Labels    map[string]string `json:"labels"`
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
}

type IngressIP struct {
	IP string `json:"ip"`
}

type LoadBalancer struct {
	Ingress []IngressIP `json:"ingress"`
}

type Status struct {
	LoadBalancer LoadBalancer `json:"loadBalancer"`
}

type Ingress struct {
	Metadata Metadata `json:"metadata"`
	Status   Status   `json:"status"`
}
