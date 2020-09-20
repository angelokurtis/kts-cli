package containers

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

func logs(cmd *cobra.Command, args []string) {
	containers, err := kubectl.ListContainers(namespace, allNamespaces)
	if err != nil {
		system.Exit(err)
	}
	selected, err := containers.SelectMany()
	if err != nil {
		system.Exit(err)
	}
	for _, container := range selected.Items {
		if download {
			err := kubectl.SaveLogs(container)
			if err != nil {
				color.Yellow.Println("[WARN] " + err.Error())
			}
		} else {
			p := container.Pod
			ns := container.Namespace
			c := container.Name
			cmd := ""
			if containers.CountByPod(container.Pod) > 1 {
				cmd = fmt.Sprintf("kubectl logs %s -c %s -n %s", p, c, ns)
			} else {
				cmd = fmt.Sprintf("kubectl logs %s -n %s", p, ns)
			}
			color.Secondary.Println(cmd)
		}
	}
}