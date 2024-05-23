package deployments

import (
	log "log/slog"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
)

func rollout(cmd *cobra.Command, args []string) {
	deploys, err := kubectl.ListDeployments(namespace, allNamespaces)
	if err != nil {
		log.Error(err.Error())
		return
	}

	deploys, err = deploys.SelectMany()
	if err != nil {
		log.Error(err.Error())
		return
	}

	err = deploys.Rollout()
	if err != nil {
		log.Error(err.Error())
		return
	}
}
