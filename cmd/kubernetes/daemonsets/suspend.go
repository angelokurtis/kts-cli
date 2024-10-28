package daemonsets

import (
	log "log/slog"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
)

func suspend(cmd *cobra.Command, args []string) {
	daemonSets, err := kubectl.ListDaemonSets(namespace, allNamespaces)
	if err != nil {
		log.Error(err.Error())
		return
	}

	daemonSets, err = daemonSets.SelectMany()
	if err != nil {
		log.Error(err.Error())
		return
	}

	if err = daemonSets.Suspend(); err != nil {
		log.Error(err.Error())
		return
	}
}
