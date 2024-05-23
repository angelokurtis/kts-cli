package pods

import (
	log "log/slog"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
)

// kube pods restart
func restart(cmd *cobra.Command, args []string) {
	pods, err := kubectl.ListPods(namespace, allNamespaces, selector)
	if err != nil {
		system.Exit(err)
	}

	pods, err = pods.SelectMany()
	if err != nil {
		system.Exit(err)
	}

	log.Info("")
}
