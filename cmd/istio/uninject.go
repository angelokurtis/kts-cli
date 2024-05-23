package istio

import (
	log "log/slog"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/pkg/app/istioctl"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
)

// istio uninject -n testeistio
func uninject(cmd *cobra.Command, args []string) {
	deployments, err := kubectl.ListDeployments(namespace, allNamespaces)
	if err != nil {
		log.Error(err.Error())
		return
	}

	deployments = deployments.FilterInjected()

	deployments, err = deployments.SelectMany()
	if err != nil {
		log.Error(err.Error())
		return
	}

	for _, deployment := range deployments.Items {
		err := istioctl.KubeUninject(deployment)
		if err != nil {
			log.Error(err.Error())
			return
		}
	}
}
