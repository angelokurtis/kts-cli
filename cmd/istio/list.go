package istio

import (
	"fmt"
	"github.com/andanhm/go-prettytime"
	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/enescakir/emoji"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

// istio list
func list(cmd *cobra.Command, args []string) {
	deployments, err := kubectl.ListDeployments(namespace, allNamespaces)
	if err != nil {
		log.Fatal(err)
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnSeparator("")
	table.SetBorder(false)
	table.SetHeaderLine(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader([]string{"", "NAME", "READY", "UP-TO-DATE", "AVAILABLE", "AGE"})
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
		table.Append([]string{deployment.StatusColor(), deployment.Metadata.Name + istio, ready, updated, available, prettytime.Format(deployment.Metadata.CreationTimestamp)})
	}
	table.Render()
}
