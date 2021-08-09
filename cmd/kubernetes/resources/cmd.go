package resources

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
)

var (
	allNamespaces = false
	group         = ""
	namespace     = ""
	Command       = &cobra.Command{
		Use:   "resources",
		Short: "Utility functions to deal with any type of resource",
		Run:   system.Help,
	}
)

func init() {
	listCMD := &cobra.Command{Use: "list", Run: list}
	listCMD.PersistentFlags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "If present, list the requested object(s) across all namespaces. Namespace in current\ncontext is ignored even if specified with --namespace.")
	listCMD.PersistentFlags().StringVar(&group, "group", "", "")
	listCMD.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "")
	Command.AddCommand(listCMD)
}
