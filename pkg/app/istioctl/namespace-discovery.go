package istioctl

import (
	"github.com/pkg/errors"

	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
)

func findIstioNamespace() (string, error) {
	pods, err := kubectl.ListPods("", true, "app=istiod")
	if err != nil {
		return "", err
	}
	for _, pod := range pods.Items {
		return pod.Metadata.Namespace, nil
	}
	return "", errors.New("unable to find any Istiod instances")
}
