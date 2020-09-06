package containers

import (
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/spf13/cobra"
)

var (
	download      = false
	allNamespaces = false
	namespace     = ""
	Command       = &cobra.Command{
		Use:   "containers",
		Short: "Utility functions to deal with Pod containers",
		Run:   system.Help,
	}
)

func init() {
	Command.PersistentFlags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "If present, list the requested container(s) across all namespaces. Namespace in current\ncontext is ignored even if specified with --namespace.")
	Command.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "")
	Command.AddCommand(&cobra.Command{Use: "forwarding", Run: forwarding})
	Command.AddCommand(&cobra.Command{Use: "list", Run: list})

	logsCommand := &cobra.Command{Use: "logs", Run: logs}
	logsCommand.PersistentFlags().BoolVarP(&download, "download", "d", false, "If present, download the logs locally")
	Command.AddCommand(logsCommand)
}
