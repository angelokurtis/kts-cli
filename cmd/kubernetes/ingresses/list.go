package ingresses

import (
	log "log/slog"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
)

func list(cmd *cobra.Command, args []string) {
	ingresses, err := kubectl.ListIngresses()
	if err != nil {
		log.Error(err.Error())
		return
	}

	services, err := kubectl.ListServices()
	if err != nil {
		log.Error(err.Error())
		return
	}

	services = services.FilterByType("LoadBalancer")

	gateways, err := kubectl.ListAllIstioGateways()
	if err != nil {
		log.Error(err.Error())
		return
	}

	if len(gateways) > 0 {
		log.Debug("found gateways", "gateways-count", len(gateways))
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetColumnSeparator("")
	table.SetBorder(false)
	table.SetHeaderLine(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader([]string{"Kind", "Namespace", "Name", "External IP", "Host"})

	for _, ingress := range ingresses {
		hosts := make([]string, 0)
		for _, rule := range ingress.Spec.Rules {
			hosts = append(hosts, rule.Host)
		}

		table.Append([]string{"Ingress", ingress.Namespace, ingress.Name, ingress.ExternalIP(), strings.Join(hosts, "\n")})
	}

	for _, service := range services.Items {
		m := service.Metadata
		table.Append([]string{"Service", m.Namespace, m.Name, service.ExternalIP()})
	}

	var istioIngress string
	if len(gateways) > 0 {
		istioIngress, err = kubectl.IstioIngress()
		if err != nil {
			log.Error(err.Error())
			return
		}
	}

	for _, gtw := range gateways {
		m := gtw.Metadata
		table.Append([]string{"Gateway", m.Namespace, m.Name, istioIngress})
	}

	table.Render()
}
