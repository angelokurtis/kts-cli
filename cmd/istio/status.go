package istio

import (
	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/pkg/app/istioctl"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/spf13/cobra"
)

// istio status
func status(cmd *cobra.Command, args []string) {
	pods, err := kubectl.ListPods(namespace, allNamespaces)
	if err != nil {
		log.Fatal(err)
	}
	pods, err = pods.SelectMany()
	if err != nil {
		log.Fatal(err)
	}
	for _, pod := range pods.Items {
		err := istioctl.ProxyStatus(pod)
		if err != nil {
			log.Fatal(err)
		}
	}
}
