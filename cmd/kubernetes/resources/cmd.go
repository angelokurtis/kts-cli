package resources

import (
	"github.com/angelokurtis/kts-cli/cmd/common"
	"github.com/spf13/cobra"
)

var (
	allNamespaces = false
	group         = ""
	Command       = &cobra.Command{
		Use:   "resources",
		Short: "Utility function to deal with Kubernetes API resources available on the server",
		Run:   common.Help,
	}
)

func init() {
	Command.PersistentFlags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "If present, list the requested object(s) across all namespaces. Namespace in current\ncontext is ignored even if specified with --namespace.")
	Command.PersistentFlags().StringVar(&group, "group", "", "")
	Command.AddCommand(&cobra.Command{Use: "list", Run: list})
}
