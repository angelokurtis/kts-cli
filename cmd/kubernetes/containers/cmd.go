package containers

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
	selector      = ""
	Command       = &cobra.Command{
		Use:   "containers",
		Short: "Utility functions to deal with Pod containers",
		Run:   system.Help,
	}
)

func init() {
	Command.PersistentFlags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "If present, list the requested container(s) across all namespaces. Namespace in current\ncontext is ignored even if specified with --namespace.")
	Command.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "")
	Command.PersistentFlags().StringVarP(&selector, "selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)")
	Command.AddCommand(&cobra.Command{Use: "forwarding", Run: forwarding})
	Command.AddCommand(&cobra.Command{Use: "list", Run: list})

	logsCommand := &cobra.Command{Use: "logs", Run: logs}
	logsCommand.PersistentFlags().BoolVarP(&download, "download", "d", false, "If present, download the logs locally.")
	logsCommand.PersistentFlags().BoolVarP(&previous, "previous", "p", false, "If true, print the logs for the previous instance of the container in a pod if it exists.")
	logsCommand.PersistentFlags().StringVarP(&since, "since", "s", "0s", "Only return logs newer than a relative duration like 5s, 2m, or 3h.")
	Command.AddCommand(logsCommand)
}
