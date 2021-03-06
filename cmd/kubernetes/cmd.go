package kubernetes

import (
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/certificates"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/clusters"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/containers"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/deployments"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/ingresses"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/labels"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/pods"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/services"
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/spf13/cobra"
)

var (
	status        = false
	allNamespaces = false
	group         = ""
	namespace     = ""
	Command       = &cobra.Command{
		Use:   "kube",
		Short: "Kubernetes container-orchestration utilities",
		Run:   system.Help,
	}
)

func init() {
	Command.AddCommand(clusters.Command)
	Command.AddCommand(containers.Command)
	Command.AddCommand(ingresses.Command)
	Command.AddCommand(services.Command)
	Command.AddCommand(labels.Command)
	Command.AddCommand(deployments.Command)
	Command.AddCommand(pods.Command)
	Command.AddCommand(certificates.Command)

	listCommand := &cobra.Command{Use: "resources", Run: resources}
	listCommand.PersistentFlags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "If present, list the requested object(s) across all namespaces. Namespace in current\ncontext is ignored even if specified with --namespace.")
	listCommand.PersistentFlags().StringVar(&group, "group", "", "")
	listCommand.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "")
	Command.AddCommand(listCommand)

	manifestsCommand := &cobra.Command{Use: "manifests", Run: manifests}
	manifestsCommand.PersistentFlags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "If present, resources the requested object(s) across all namespaces. Namespace in current\ncontext is ignored even if specified with --namespace.")
	manifestsCommand.PersistentFlags().StringVar(&group, "group", "", "")
	manifestsCommand.PersistentFlags().BoolVar(&status, "status", false, "")
	manifestsCommand.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "")
	Command.AddCommand(manifestsCommand)
}
