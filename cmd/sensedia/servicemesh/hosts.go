package servicemesh

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/cmd/common"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/spf13/cobra"
)

func hosts(cmd *cobra.Command, args []string) {
	ingress, err := kubectl.IstioIngress()
	if err != nil {
		common.Exit(err)
	}

	hosts, err := kubectl.MeshesHosts()
	if err != nil {
		common.Exit(err)
	}

	for _, host := range hosts {
		fmt.Printf("%s\t\t%s\n", ingress, host)
	}
}
