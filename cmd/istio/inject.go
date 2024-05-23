package istio

import (
	log "log/slog"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/pkg/app/istioctl"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
)

// istio inject
func inject(cmd *cobra.Command, args []string) {
	deployments, err := kubectl.ListDeployments(namespace, allNamespaces)
	if err != nil {
		log.Error(err.Error())
		return
	}

	deployments = deployments.FilterUninjected()
	// deployments = deployments.FilterInjectable()
	deployments, err = deployments.SelectMany()
	if err != nil {
		log.Error(err.Error())
		return
	}

	for _, deployment := range deployments.Items {
		err := istioctl.KubeInject(deployment)
		if err != nil {
			log.Error(err.Error())
			return
		}
	}
}
