package kubernetes

import (
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/clusters"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/containers"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/resources"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/services"
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "kube",
	Short: "Kubernetes container-orchestration utilities",
	Run:   system.Help,
}

func init() {
	Command.AddCommand(&cobra.Command{Use: "hosts", Run: hosts})
	Command.AddCommand(clusters.Command)
	Command.AddCommand(containers.Command)
	Command.AddCommand(resources.Command)
	Command.AddCommand(services.Command)
}
