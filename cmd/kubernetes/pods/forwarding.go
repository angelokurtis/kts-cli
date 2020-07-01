package pods

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/cmd/common"
	"github.com/angelokurtis/kts-cli/internal/color"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl/portforward"
	"github.com/spf13/cobra"
)

func forwarding(cmd *cobra.Command, args []string) {
	pods, err := kubectl.ListAllPods()
	if err != nil {
		common.Exit(err)
	}
	labels, err := pods.SelectLabels()
	if err != nil {
		common.Exit(err)
	}
	namespace, err := pods.SelectNamespace(labels)
	if err != nil {
		common.Exit(err)
	}
	port, err := pods.SelectContainerPort(namespace, labels)
	if err != nil {
		common.Exit(err)
	}
	command := portforward.NewCommand(namespace, labels, port)
	fmt.Printf(color.Warning, command.String()+"\n")
}
