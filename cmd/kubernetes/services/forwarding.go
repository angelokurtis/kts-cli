package services

import (
	"github.com/angelokurtis/kts-cli/cmd/common"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/angelokurtis/kts-cli/pkg/app/kubefwd"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

func forwarding(cmd *cobra.Command, args []string) {
	services, err := kubectl.ListAllServices()
	if err != nil {
		common.Exit(err)
	}
	labels, err := services.SelectLabels()
	if err != nil {
		common.Exit(err)
	}
	namespaces := services.Namespaces(labels)
	command := kubefwd.NewCommand(labels, namespaces)
	color.Secondary.Println(command.String())
}
