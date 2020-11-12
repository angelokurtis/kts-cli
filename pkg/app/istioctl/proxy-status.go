package istioctl

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func ProxyStatus(pod *kubectl.Pod) error {
	_, err := bash.RunAndLogRead(fmt.Sprintf("istioctl proxy-status %s.%s", pod.Metadata.Name, pod.Metadata.Namespace))
	return err
}
