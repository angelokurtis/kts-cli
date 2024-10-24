package pods

import (
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

// kube pods list
func list(cmd *cobra.Command, args []string) {
	pods, err := kubectl.ListPods(namespace, allNamespaces, selector)
	if err != nil {
		log.Error(err.Error())
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnSeparator("")
	table.SetBorder(false)
	table.SetHeaderLine(false)
	table.SetColWidth(100)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	if allNamespaces {
		table.SetHeader([]string{"", "NAMESPACE", "NAME", "READY", "STATUS", "RESTARTS", "AGE", "LAST UPDATE"})
	} else {
		table.SetHeader([]string{"", "NAME", "READY", "STATUS", "RESTARTS", "AGE", "LAST UPDATE"})
	}

	if sortUpdated {
		sort.Slice(pods.Items, func(i, j int) bool {
			it := pods.Items[i].LastUpdate()
			jt := pods.Items[j].LastUpdate()

			return it.Before(*jt)
		})
	}

	for _, pod := range pods.Items {
		job := func() string {
			if pod.IsJob() {
				return " " + emoji.Robot.String()
			}

			return ""
		}()
		istio := func() string {
			if pod.HasIstioSidecar() {
				return " " + emoji.AlienMonster.String()
			}

			return ""
		}()
		system := func() string {
			if pod.Metadata.Namespace == "kube-system" || pod.Metadata.Namespace == "local-path-storage" {
				return " " + emoji.Laptop.String()
			}

			return ""
		}()

		if allNamespaces {
			table.Append([]string{pod.StatusColor(), pod.Metadata.Namespace, pod.Metadata.Name + job + istio + system, pod.Ready(), pod.CurrentStatus(), strconv.Itoa(pod.RestartCount()), prettytime.Format(pod.Metadata.CreationTimestamp), prettytime.Format(*pod.LastUpdate())})
		} else {
			table.Append([]string{pod.StatusColor(), pod.Metadata.Name + job + istio + system, pod.Ready(), pod.CurrentStatus(), strconv.Itoa(pod.RestartCount()), prettytime.Format(pod.Metadata.CreationTimestamp), prettytime.Format(*pod.LastUpdate())})
		}
	}

	table.Render()
}
