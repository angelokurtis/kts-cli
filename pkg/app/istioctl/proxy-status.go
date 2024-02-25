package istioctl

import (
	"fmt"

	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func ProxyStatus(pod *kubectl.Pod) error {
	ns, err := findIstioNamespace()
	if err != nil {
		return err
	}

	out, err := bash.RunAndLogRead(fmt.Sprintf("istioctl proxy-status %s.%s -i %s", pod.Metadata.Name, pod.Metadata.Namespace, ns))
	fmt.Println(string(out))

	return err
}
