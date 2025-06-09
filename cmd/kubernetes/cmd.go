package kubernetes

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/cmd/kubernetes/certificates"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/clusters"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/containers"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/daemonsets"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/deployments"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/events"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/ingresses"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/labels"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/nodes"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/pods"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/resources"
	"github.com/angelokurtis/kts-cli/cmd/kubernetes/services"
	"github.com/angelokurtis/kts-cli/internal/system"
)

var (
	sanitize      = false
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
	Command.AddCommand(events.Command)
	Command.AddCommand(ingresses.Command)
	Command.AddCommand(services.Command)
	Command.AddCommand(labels.Command)
	Command.AddCommand(nodes.Command)
	Command.AddCommand(daemonsets.Command)
	Command.AddCommand(deployments.Command)
	Command.AddCommand(pods.Command)
	Command.AddCommand(certificates.Command)
	Command.AddCommand(resources.Command)

	manifestsCommand := &cobra.Command{Use: "manifests", Run: manifests}
	manifestsCommand.PersistentFlags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "If set, retrieves the specified resources from all namespaces. Overrides any namespace set with --namespace.")
	manifestsCommand.PersistentFlags().StringVar(&group, "group", "", "Filter the resources by API group (e.g., 'apps', 'batch'). Leave empty to include all groups.")
	manifestsCommand.PersistentFlags().BoolVar(&status, "status", false, "Include status fields in the output YAML manifests, if available.")
	manifestsCommand.PersistentFlags().BoolVar(&decodeSecrets, "decode-secrets", false, "If true, decodes Secret data fields from base64 to plain text in the output.")
	manifestsCommand.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Specify the namespace to query. Ignored if --all-namespaces is set.")
	manifestsCommand.PersistentFlags().BoolVar(&sanitize, "sanitize", false, "Remove auto-generated fields from the output YAML (e.g., status, creationTimestamp, managedFields). This is useful for producing clean manifests suitable for version control or reuse.")
	Command.AddCommand(manifestsCommand)
}
