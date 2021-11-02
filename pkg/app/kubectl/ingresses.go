package kubectl

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	extensions "k8s.io/api/extensions/v1beta1"
	apis "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/angelokurtis/kts-cli/internal/kube"
)

func ListIngresses() ([]*Ingress, error) {
	ingresses, err := kube.GetClientset().ExtensionsV1beta1().Ingresses("").List(context.TODO(), apis.ListOptions{})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	items := make([]*Ingress, 0)
	for _, ing := range ingresses.Items {
		ingress := Ingress(ing)
		items = append(items, &ingress)
	}

	return items, nil
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
	OwnerReferences []*OwnerReference `json:"ownerReferences"`
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

type Ingress extensions.Ingress

func (i *Ingress) ExternalIP() string {
	for _, ingress := range i.Status.LoadBalancer.Ingress {
		if ingress.Hostname != "" {
			return ingress.Hostname
		}
		if ingress.IP != "" {
			return ingress.IP
		}
	}
	return ""
}
