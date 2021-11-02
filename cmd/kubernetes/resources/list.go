package resources

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
)

func list(cmd *cobra.Command, args []string) {
	resources := ""
	if len(args) == 0 {
		rd, err := kubectl.ListResourceDefinitions()
		if err != nil {
			log.Fatal(err)
		}
		rd = rd.FilterVerbs("list")
		if !allNamespaces {
			rd = rd.FilterNamespaced()
		}
		if group != "" {
			rd = rd.FilterAPIGroup(group)
		}
		resources = rd.String()
	} else {
		resources = args[0]
	}
	if owners {
		ro, err := kubectl.ListResourcesOwners(resources, namespace, allNamespaces)
		if err != nil {
			log.Fatal(err)
		}
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.SetColumnSeparator("")
		table.SetBorder(false)
		table.SetHeaderLine(false)
		table.SetColWidth(100)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		if allNamespaces {
			table.SetHeader([]string{"NAMESPACE", "Kind", "Name", "Dependents"})
		} else {
			table.SetHeader([]string{"Kind", "Name", "Dependents"})
		}
		for _, item := range ro {
			m := item.Metadata
			if allNamespaces {
				table.Append([]string{m.Namespace, item.Kind, m.Name, fmt.Sprintf("%d", item.Dependents)})
			} else {
				table.Append([]string{item.Kind, m.Name, fmt.Sprintf("%d", item.Dependents)})
			}
		}
		table.Render()
		return
	}
	results, err := kubectl.ListResources(resources, namespace, allNamespaces)
	if err != nil {
		log.Fatal(err)
	}
	for _, result := range results {
		fmt.Println(result)
	}
}
