package containers

import (
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"strconv"
)

// kube containers forwarding
func forwarding(cmd *cobra.Command, args []string) {
	containers, err := kubectl.ListContainers(namespace, allNamespaces, "")
	if err != nil {
		system.Exit(err)
	}
	containers = containers.FilterExposed()
	containers, err = containers.SelectMany()
	if err != nil {
		system.Exit(err)
	}
	for _, container := range containers.Items {
		for _, port := range container.Ports {
			n := container.Pod
			ns := container.Namespace
			rp := strconv.Itoa(port.ContainerPort)
			lp := rp
			if port.ContainerPort == 80 {
				lp = "8000"
			}
			color.Secondary.Println("kubectl port-forward " + n + " " + lp + ":" + rp + " -n " + ns)
		}
	}
}
