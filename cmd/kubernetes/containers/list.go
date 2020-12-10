package containers

import (
	"github.com/andanhm/go-prettytime"
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
)

// kube containers list
func list(cmd *cobra.Command, args []string) {
	containers, err := kubectl.ListContainers(namespace, allNamespaces, selector)
	if err != nil {
		system.Exit(err)
	}
	if sortUpdated {
		//sort.Slice(containers.Items, func(i, j int) bool {
		//	it := containers.Items[i].LastUpdateTime()
		//	jt := containers.Items[j].LastUpdateTime()
		//	return it.Before(*jt)
		//})
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnSeparator("")
	table.SetBorder(false)
	table.SetHeaderLine(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader([]string{"", "Container", "Ports", "Image", "Pod", "Age"})
	for _, container := range containers.Items {
		ports := make([]string, 0, len(container.Ports))
		for _, port := range container.Ports {
			ports = append(ports, strconv.Itoa(port.ContainerPort))
		}
		state := container.GetState()
		color := ""
		if state != nil {
			color = state.Color()
		}
		timeStr := ""
		updateTime := container.LastUpdateTime()
		if updateTime != nil {
			timeStr = prettytime.Format(*updateTime)
		}
		table.Append([]string{color, container.Name, strings.Join(ports, ","), container.Image, container.Pod, timeStr})
	}
	table.Render()
}
