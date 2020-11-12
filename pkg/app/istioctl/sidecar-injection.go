package istioctl

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/angelokurtis/kts-cli/pkg/bash"
	"strings"
)

func AddToMesh(deployment *kubectl.Deployment) error {
	name := deployment.Metadata.Name
	namespace := deployment.Metadata.Namespace
	_, err := bash.RunAndLogWrite(fmt.Sprintf("istioctl experimental add-to-mesh deployment %s -n %s", name, namespace))
	if err != nil {
		if !strings.Contains(err.Error(), "0 errors occurred") {
			return err
		}
	}
	return nil
}

func RemoveFromMesh(deployment *kubectl.Deployment) error {
	name := deployment.Metadata.Name
	namespace := deployment.Metadata.Namespace
	_, err := bash.RunAndLogWrite(fmt.Sprintf("istioctl experimental remove-from-mesh deployment %s -n %s", name, namespace))
	if err != nil {
		if !strings.Contains(err.Error(), "0 errors occurred") {
			return err
		}
	}
	return nil
}

func KubeInject(deployment *kubectl.Deployment) error {
	name := deployment.Metadata.Name
	namespace := deployment.Metadata.Namespace
	_, err := bash.RunAndLogWrite(fmt.Sprintf("kubectl get deployment %s -n %s -o yaml | istioctl kube-inject -f - | kubectl apply -f -", name, namespace))
	return err
}

func KubeUninject(deployment *kubectl.Deployment) error {
	name := deployment.Metadata.Name
	namespace := deployment.Metadata.Namespace
	_, err := bash.RunAndLogWrite(fmt.Sprintf("kubectl get deployment %s -n %s -o yaml | istioctl x kube-uninject -f - | kubectl apply -f -", name, namespace))
	return err
}
