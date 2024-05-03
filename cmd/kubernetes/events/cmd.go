package events

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
)

var (
	since         time.Duration
	allNamespaces = false
	namespace     = ""
	Command       = &cobra.Command{
		Use:   "events",
		Short: "Display events from the Kubernetes cluster",
		Run:   system.Help,
	}
)

func init() {
	Command.PersistentFlags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "If present, list the requested container(s) across all namespaces. Namespace in current\ncontext is ignored even if specified with --namespace.")
	Command.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "")
	Command.PersistentFlags().DurationVarP(&since, "since", "s", time.Second*0, "Only return events newer than a relative duration like 5s, 2m, or 3h.")
	Command.AddCommand(&cobra.Command{Use: "list", Run: list})
}
