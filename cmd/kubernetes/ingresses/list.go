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

	gateways, err := kubectl.ListAllIstioGateways()
	if err != nil {
		log.Fatal(err)
	}
	if len(gateways) > 0 {
		log.Debugf("found %d gateways\n", len(gateways))
	}

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

	var istioIngress string
	if len(gateways) > 0 {
		istioIngress, err = kubectl.IstioIngress()
		if err != nil {
			log.Fatal(err)
		}
	}
	for _, gtw := range gateways {
		m := gtw.Metadata
		table.Append([]string{"Gateway", m.Namespace, m.Name, istioIngress})
	}
	table.Render()
}
