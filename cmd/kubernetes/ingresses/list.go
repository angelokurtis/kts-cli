package ingresses

import (
	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

func list(cmd *cobra.Command, args []string) {
	ingresses, err := kubectl.ListIngresses()
	if err != nil {
		log.Fatal(err)
	}
	services, err := kubectl.ListServices()
	if err != nil {
		log.Fatal(err)
	}
	services = services.FilterByType("LoadBalancer")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Kind", "Namespace", "Name", "External IP"})
	table.SetBorder(false)
	for _, ingress := range ingresses {
		m := ingress.Metadata
		table.Append([]string{"Ingress", m.Namespace, m.Name, ingress.ExternalIP()})
	}
	for _, service := range services.Items {
		m := service.Metadata
		table.Append([]string{"Service", m.Namespace, m.Name, service.ExternalIP()})
	}
	table.Render()
}
