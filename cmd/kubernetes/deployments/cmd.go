package deployments

import (
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/spf13/cobra"
)

var (
	previous      = false
	download      = false
	since         = ""
	allNamespaces = false
	namespace     = ""
	Command       = &cobra.Command{
		Use:   "deployments",
		Short: "Utility functions to deal with Deployments",
		Run:   system.Help,
	}
)

func init() {
	Command.PersistentFlags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "If present, list the requested container(s) across all namespaces. Namespace in current\ncontext is ignored even if specified with --namespace.")
	Command.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "")
	Command.AddCommand(&cobra.Command{Use: "rollout", Run: rollout})
	Command.AddCommand(&cobra.Command{Use: "update-images", Run: updateImages})

	logsCommand := &cobra.Command{Use: "logs", Run: logs}
	logsCommand.PersistentFlags().BoolVarP(&download, "download", "d", false, "If present, download the logs locally.")
	logsCommand.PersistentFlags().BoolVarP(&previous, "previous", "p", false, "If true, print the logs for the previous instance of the container in a pod if it exists.")
	logsCommand.PersistentFlags().StringVarP(&since, "since", "s", "0s", "Only return logs newer than a relative duration like 5s, 2m, or 3h.")
	Command.AddCommand(logsCommand)

}
