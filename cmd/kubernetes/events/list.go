package events

import (
	"fmt"
	"os"
	"sort"
	"time"

	prettytime "github.com/andanhm/go-prettytime"
	"github.com/olekukonko/tablewriter"
	"github.com/samber/lo"
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

	if since != time.Second*0 {
		events.Items = lo.Filter(events.Items, func(item *kubectl.Event, index int) bool {
			from := time.Now().Add(since * -1)
			return item.LastSeenTimestamp().After(from)
		})
	}

	sort.Slice(events.Items, func(a, b int) bool {
		timeA := events.Items[a].LastSeenTimestamp()
		timeB := events.Items[b].LastSeenTimestamp()

		return timeA.Before(timeB)
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnSeparator("")
	table.SetBorder(false)
	table.SetHeaderLine(false)
	table.SetColWidth(1000)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	if allNamespaces {
		table.SetHeader([]string{"NAMESPACE", "LAST SEEN", "TYPE", "REASON", "OBJECT", "MESSAGE"})
	} else {
		table.SetHeader([]string{"LAST SEEN", "TYPE", "REASON", "OBJECT", "MESSAGE"})
	}

	for _, event := range events.Items {
		lastSeenTimestamp := event.LastSeenTimestamp().In(spTimeZone)
		prettyTimestamp := fmt.Sprintf("%s (%s)", lastSeenTimestamp.Format("02/01/2006 15:04"), prettytime.Format(lastSeenTimestamp))
		eventResource := event.InvolvedObject.Kind + "/" + event.InvolvedObject.Name

		if allNamespaces {
			table.Append([]string{event.InvolvedObject.Namespace, prettyTimestamp, event.Type, event.Reason, eventResource, event.Message})
		} else {
			table.Append([]string{prettyTimestamp, event.Type, event.Reason, eventResource, event.Message})
		}
	}

	table.Render()
}
