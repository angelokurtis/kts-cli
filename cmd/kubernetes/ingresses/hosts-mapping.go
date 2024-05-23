package ingresses

import (
	"fmt"
	log "log/slog"
	"strings"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/angelokurtis/kts-cli/pkg/app/linux"
)

// kube ingresses hosts-mapping
func hostsMapping(cmd *cobra.Command, args []string) {
	context, err := kubectl.CurrentContext()
	if err != nil {
		system.Exit(err)
	}

	ingresses, err := kubectl.ListIngresses()
	if err != nil {
		system.Exit(err)
	}

	if len(ingresses) > 0 {
		log.Debug("found ingresses", "ingresses-count", len(ingresses))
	}

	gateways, err := kubectl.ListAllIstioGateways()
	if err != nil {
		system.Exit(err)
	}

	if len(gateways) > 0 {
		log.Debug("found gateways", "gateways-count", len(gateways))
	}

	hosts, err := linux.LoadHostsFile()
	if err != nil {
		system.Exit(err)
	}

	err = hosts.Add(context, ingresses, gateways)
	if err != nil {
		system.Exit(err)
	}

	err = hosts.Write()
	if err != nil {
		if !strings.Contains(err.Error(), "open /etc/hosts: permission denied") {
			system.Exit(err)
		}

		fmt.Printf("This command requires superuser privileges to run. These\nprivileges are required to add IP address aliases to your\nloopback interface.\n\nTry:\n - sudo -E kts kube ingresses hosts-mapping\n")
	} else {
		log.Info("/etc/hosts file has been rewritten!")
	}
}
