package istio

import (
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/spf13/cobra"
)

var (
	sortUpdated   = false
	allNamespaces = false
	namespace     = ""
	filename      = ""
	Command       = &cobra.Command{
		Use:   "istio",
		Short: "Utilities to deal with Istio's service mesh platform",
		Run:   system.Help,
	}
)

func init() {
	injectCmd := &cobra.Command{Use: "inject", Run: inject}
	injectCmd.PersistentFlags().StringVarP(&filename, "filename", "f", "", "that contains the configuration to apply")
	injectCmd.PersistentFlags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "If present, list the requested object(s) across all namespaces. Namespace in current\ncontext is ignored even if specified with --namespace.")
	injectCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "")
	Command.AddCommand(injectCmd)

	uninjectCmd := &cobra.Command{Use: "uninject", Run: uninject}
	uninjectCmd.PersistentFlags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "If present, list the requested object(s) across all namespaces. Namespace in current\ncontext is ignored even if specified with --namespace.")
	uninjectCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "")
	Command.AddCommand(uninjectCmd)

	listCmd := &cobra.Command{Use: "list", Run: list}
	listCmd.PersistentFlags().BoolVar(&sortUpdated, "sort-updated", false, "")
	listCmd.PersistentFlags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "If present, list the requested object(s) across all namespaces. Namespace in current\ncontext is ignored even if specified with --namespace.")
	listCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "")
	Command.AddCommand(listCmd)

	statusCmd := &cobra.Command{Use: "status", Run: status}
	statusCmd.PersistentFlags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "If present, list the requested object(s) across all namespaces. Namespace in current\ncontext is ignored even if specified with --namespace.")
	statusCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "")
	Command.AddCommand(statusCmd)
}
