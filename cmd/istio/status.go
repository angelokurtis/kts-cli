package istio

import (
	log "log/slog"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/pkg/app/istioctl"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
)

// istio status
func status(cmd *cobra.Command, args []string) {
	pods, err := kubectl.ListPods(namespace, allNamespaces, "")
	if err != nil {
		log.Error(err.Error())
		return
	}

	pods, err = pods.SelectMany()
	if err != nil {
		log.Error(err.Error())
		return
	}

	for _, pod := range pods.Items {
		err := istioctl.ProxyStatus(pod)
		if err != nil {
			log.Error(err.Error())
			return
		}
	}
}
