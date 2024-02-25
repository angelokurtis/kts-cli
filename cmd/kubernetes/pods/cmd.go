package pods

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
)

var (
	allNamespaces = false
	namespace     = ""
	selector      = ""
	sortUpdated   = false
	Command       = &cobra.Command{
		Use:   "pods",
		Short: "Utility functions to deal with Pods",
		Run:   system.Help,
	}
)

func init() {
	Command.PersistentFlags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "If present, list the requested container(s) across all namespaces. Namespace in current\ncontext is ignored even if specified with --namespace.")
	Command.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "")
	Command.AddCommand(&cobra.Command{Use: "selectors", Run: selectors})

	Command.AddCommand(&cobra.Command{Use: "variables", Run: variables})

	Command.AddCommand(&cobra.Command{Use: "restart", Run: restart})

	listCMD := &cobra.Command{Use: "list", Run: list}
	listCMD.PersistentFlags().StringVarP(&selector, "selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)")
	listCMD.PersistentFlags().BoolVar(&sortUpdated, "sort-updated", false, "")
	Command.AddCommand(listCMD)
}
