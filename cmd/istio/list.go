package istio

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
)

// istio list
func list(cmd *cobra.Command, args []string) {
	deployments, err := kubectl.ListDeployments(namespace, allNamespaces)
	if err != nil {
		log.Fatal(err)
	}
	for _, deployment := range deployments.Items {
		if deployment.HasIstioSidecar() {
			_, err := emoji.Printf(":pager:%s/%s\n", deployment.Metadata.Namespace, deployment.Metadata.Name)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Printf("   %s/%s\n", deployment.Metadata.Namespace, deployment.Metadata.Name)
		}
	}
}
