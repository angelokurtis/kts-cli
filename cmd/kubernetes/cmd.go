package kubernetes

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/cmd/kubernetes/certificates"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/clusters"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/containers"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/deployments"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/ingresses"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/labels"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/pods"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/resources"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/services"
	"github.com/angelokurtis/kts-cli/internal/system"
)

var (
	status        = false
	decodeSecrets = false
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
	Command.AddCommand(resources.Command)

	manifestsCommand := &cobra.Command{Use: "manifests", Run: manifests}
	manifestsCommand.PersistentFlags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "If present, resources the requested object(s) across all namespaces. Namespace in current\ncontext is ignored even if specified with --namespace.")
	manifestsCommand.PersistentFlags().StringVar(&group, "group", "", "")
	manifestsCommand.PersistentFlags().BoolVar(&status, "status", false, "")
	manifestsCommand.PersistentFlags().BoolVar(&decodeSecrets, "decode-secrets", false, "")
	manifestsCommand.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "")
	Command.AddCommand(manifestsCommand)
}
