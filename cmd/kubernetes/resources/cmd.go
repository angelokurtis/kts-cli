package resources

import (
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/spf13/cobra"
)

var (
	allNamespaces = false
	group         = ""
	namespace     = ""
	Command       = &cobra.Command{
		Use:   "resources",
		Short: "Utility function to deal with Kubernetes API resources available on the server",
		Run:   system.Help,
	}
)

func init() {
	Command.PersistentFlags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "If present, list the requested object(s) across all namespaces. Namespace in current\ncontext is ignored even if specified with --namespace.")
	Command.PersistentFlags().StringVar(&group, "group", "", "")
	Command.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "")
	Command.AddCommand(&cobra.Command{Use: "list", Run: list})
	Command.AddCommand(&cobra.Command{Use: "manifests", Run: manifests})
}
