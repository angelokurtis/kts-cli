package istio

import (
	"fmt"
	log "log/slog"
	"os"
	"sort"
	"strconv"

	prettytime "github.com/andanhm/go-prettytime"
	"github.com/enescakir/emoji"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
)

// istio list
func list(cmd *cobra.Command, args []string) {
	deployments, err := kubectl.ListDeployments(namespace, allNamespaces)
	if err != nil {
		log.Error(err.Error())
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnSeparator("")
	table.SetBorder(false)
	table.SetHeaderLine(false)
	table.SetColWidth(50)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	if allNamespaces {
		table.SetHeader([]string{"", "NAMESPACE", "NAME", "READY", "UP-TO-DATE", "AVAILABLE", "AGE", "LAST UPDATE"})
	} else {
		table.SetHeader([]string{"", "NAME", "READY", "UP-TO-DATE", "AVAILABLE", "AGE", "LAST UPDATE"})
	}

	if sortUpdated {
		sort.Slice(deployments.Items, func(i, j int) bool {
			return deployments.Items[i].LastUpdateTime().Before(*deployments.Items[j].LastUpdateTime())
		})
	}

	for _, deployment := range deployments.Items {
		istio := func() string {
			if deployment.HasIstioSidecar() {
				return " " + emoji.AlienMonster.String()
			}

			return ""
		}()
		ready := fmt.Sprintf("%d/%d", deployment.Status.ReadyReplicas, deployment.Spec.Replicas)
		updated := strconv.Itoa(deployment.Status.UpdatedReplicas)
		available := strconv.Itoa(deployment.Status.AvailableReplicas)

		if allNamespaces {
			table.Append([]string{deployment.StatusColor(), deployment.Metadata.Namespace, deployment.Metadata.Name + istio, ready, updated, available, prettytime.Format(deployment.Metadata.CreationTimestamp), prettytime.Format(*deployment.LastUpdateTime())})
		} else {
			table.Append([]string{deployment.StatusColor(), deployment.Metadata.Name + istio, ready, updated, available, prettytime.Format(deployment.Metadata.CreationTimestamp), prettytime.Format(*deployment.LastUpdateTime())})
		}
	}

	table.Render()
}
