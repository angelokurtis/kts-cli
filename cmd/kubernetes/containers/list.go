package containers

import (
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
)

func list(cmd *cobra.Command, args []string) {
	containers, err := kubectl.ListContainers(namespace, allNamespaces)
	if err != nil {
		system.Exit(err)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetHeader([]string{"Container", "Pod", "Ports", "Image", "Namespace"})
	for _, container := range containers.Items {
		ports := make([]string, 0, len(container.Ports))
		for _, port := range container.Ports {
			ports = append(ports, strconv.Itoa(port.ContainerPort))
		}
		table.Append([]string{container.Name, container.Pod, strings.Join(ports, ","), container.Image, container.Namespace})
	}
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Render()
}
