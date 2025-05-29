package daemonsets

import (
	log "log/slog"

	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/spf13/cobra"
)

func activate(cmd *cobra.Command, args []string) {
	daemonSets, err := kubectl.ListDaemonSets(namespace, allNamespaces)
	if err != nil {
		log.Error(err.Error())
		return
	}

	daemonSets, err = daemonSets.FilterSuspended()
	if err != nil {
		log.Error(err.Error())
		return
	}

	daemonSets, err = daemonSets.SelectMany()
	if err != nil {
		log.Error(err.Error())
		return
	}

	if err = daemonSets.Activate(); err != nil {
		log.Error(err.Error())
		return
	}
}
