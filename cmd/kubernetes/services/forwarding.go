package services

import (
	"github.com/gookit/color"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/angelokurtis/kts-cli/pkg/app/kubefwd"
)

func forwarding(cmd *cobra.Command, args []string) {
	services, err := kubectl.ListServices()
	if err != nil {
		system.Exit(err)
	}

	labels, err := services.SelectLabels()
	if err != nil {
		system.Exit(err)
	}

	namespaces := services.Namespaces(labels)
	command := kubefwd.NewCommand(labels, namespaces)
	color.Secondary.Println(command.String())
}
