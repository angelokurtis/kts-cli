package events

import (
	"fmt"
	"os"
	"time"

	prettytime "github.com/andanhm/go-prettytime"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
)

// kube events list
func list(cmd *cobra.Command, args []string) {
	spTimeZone, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		system.Exit(err)
	}

	events, err := kubectl.ListEvents(namespace, allNamespaces)
	if err != nil {
		system.Exit(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnSeparator("")
	table.SetBorder(false)
	table.SetHeaderLine(false)
	table.SetColWidth(100)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	if allNamespaces {
		table.SetHeader([]string{"LAST SEEN", "TYPE", "REASON", "OBJECT", "MESSAGE"})
	} else {
		table.SetHeader([]string{"NAMESPACE", "LAST SEEN", "TYPE", "REASON", "OBJECT", "MESSAGE"})
	}

	for _, event := range events.Items {
		timeFormatted := event.FirstTimestamp
		prettyTimestamp := fmt.Sprintf("%s (%s)", timeFormatted.In(spTimeZone).Format("02/01/2006 15:04"), prettytime.Format(timeFormatted))
		eventResource := event.InvolvedObject.Kind + "/" + event.InvolvedObject.Name

		if allNamespaces {
			table.Append([]string{namespace, prettyTimestamp, event.Type, event.Reason, eventResource, event.Message})
		} else {
			table.Append([]string{prettyTimestamp, event.Type, event.Reason, eventResource, event.Message})
		}
	}

	table.Render()
}
