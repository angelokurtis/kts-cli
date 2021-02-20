package istio

import (
	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/pkg/app/istioctl"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/spf13/cobra"
)

// istio uninject -n testeistio
func uninject(cmd *cobra.Command, args []string) {
	deployments, err := kubectl.ListDeployments(namespace, allNamespaces)
	if err != nil {
		log.Fatal(err)
	}
	deployments = deployments.FilterInjected()
	deployments, err = deployments.SelectMany()
	if err != nil {
		log.Fatal(err)
	}
	for _, deployment := range deployments.Items {
		err := istioctl.KubeUninject(deployment)
		if err != nil {
			log.Fatal(err)
		}
	}
}
