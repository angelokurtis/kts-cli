package kubernetes

import (
	"github.com/angelokurtis/kts-cli/cmd/common"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/resources"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/services"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "kubernetes",
	Short: "Kubernetes container-orchestration utilities",
	Run:   common.Help,
}

func init() {
	Command.AddCommand(resources.Command)
	Command.AddCommand(services.Command)
}
