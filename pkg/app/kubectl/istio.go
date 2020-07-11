package kubectl

import (
	"errors"
	"strings"
)

func IstioIngress() (string, error) {
	out, err := runAndLog("get", "service", "--all-namespaces", "-l", "app=istio-ingressgateway", "-o=jsonpath='{.items[*].status.loadBalancer.ingress[*].ip}'")
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
