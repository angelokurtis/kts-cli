package deployments

import (
	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/spf13/cobra"
)

func rollout(cmd *cobra.Command, args []string) {
	deploys, err := kubectl.ListDeployments(namespace, allNamespaces)
	if err != nil {
		log.Fatal(err)
	}
	deploys, err = deploys.SelectMany()
	if err != nil {
		log.Fatal(err)
	}
	err = deploys.Rollout()
	if err != nil {
		log.Fatal(err)
	}
}
