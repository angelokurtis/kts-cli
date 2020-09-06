package kubernetes

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/angelokurtis/kts-cli/pkg/app/linux"
	"github.com/spf13/cobra"
	"strings"
)

func hosts(cmd *cobra.Command, args []string) {
	context, err := kubectl.CurrentContext()
	if err != nil {
		system.Exit(err)
	}
	log.Infof("identified that the current context is %s\n", context)

	ingresses, err := kubectl.ListAllIngresses()
	if err != nil {
		system.Exit(err)
	}
	log.Infof("found %d ingresses\n", len(ingresses))

	gateways, err := kubectl.ListAllIstioGateways()
	if err != nil {
		system.Exit(err)
	}
	log.Infof("found %d gateways\n", len(gateways))

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
		fmt.Printf("This command requires superuser privileges to run. These\nprivileges are required to add IP address aliases to your\nloopback interface.\n\nTry:\n - sudo -E kts kube hosts\n")
	} else {
		log.Info("/etc/hosts file has been rewritten!")
	}
}
