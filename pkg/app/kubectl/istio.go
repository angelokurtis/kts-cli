package kubectl

import (
	"encoding/json"
	"github.com/pkg/errors"
	"strings"
)

func IstioIngress() (string, error) {
	out, err := run("get", "service", "--all-namespaces", "-l", "app=istio-ingressgateway", "-o=jsonpath='{.items[*].status.loadBalancer.ingress[*].ip}'")
	if err != nil {
		return "", err
	}
	str := string(out)
	str = str[1:]
	str = str[:len(str)-1]
	ingresses := strings.Split(str, " ")
	if len(ingresses) == 1 {
		return ingresses[0], nil
	} else if len(ingresses) > 1 {
		return "", errors.New("found multiple Istio ingresses addresses")
	}
	return "", nil
}

func ListAllIstioGateways() ([]*IstioGateway, error) {
	out, err := run("get", "gateways.networking.istio.io", "--all-namespaces", "-o=json", "--request-timeout=5s")
	if err != nil {
		if strings.Contains(err.Error(), "the server doesn't have a resource type") {
			return nil, nil
		}
		return nil, err
	}

	var gateways *IstioGateways
	if err := json.Unmarshal(out, &gateways); err != nil {
		return nil, errors.WithStack(err)
	}

	return gateways.Items, nil
}

type IstioGateways struct {
	Items []*IstioGateway `json:"items"`
}

type IstioGateway struct {
	Metadata struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
	} `json:"metadata"`
	Spec struct {
		Selector map[string]string `json:"selector"`
		Servers  []struct {
			Hosts []string `json:"hosts"`
		} `json:"servers"`
	} `json:"spec"`
}
