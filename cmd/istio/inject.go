package istio

import (
	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/pkg/app/istioctl"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/spf13/cobra"
)

// istio inject
func inject(cmd *cobra.Command, args []string) {
	deployments, err := kubectl.ListDeployments(namespace, allNamespaces)
	if err != nil {
		log.Fatal(err)
	}
	deployments = deployments.FilterUninjected()
	//deployments = deployments.FilterInjectable()
	deployments, err = deployments.SelectMany()
	if err != nil {
		log.Fatal(err)
	}
	for _, deployment := range deployments.Items {
		err := istioctl.KubeInject(deployment)
		if err != nil {
			log.Fatal(err)
		}
	}
}
