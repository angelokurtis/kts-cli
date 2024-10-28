package daemonsets

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
)

var (
	allNamespaces = false
	namespace     = ""
	Command       = &cobra.Command{
		Use:   "daemonsets",
		Short: "Utility functions to deal with DaemonSets",
		Run:   system.Help,
	}
)

func init() {
	Command.PersistentFlags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "If present, list the requested container(s) across all namespaces. Namespace in current\ncontext is ignored even if specified with --namespace.")
	Command.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "")
	Command.AddCommand(&cobra.Command{Use: "suspend", Run: suspend})
	Command.AddCommand(&cobra.Command{Use: "activate", Run: activate})
}
