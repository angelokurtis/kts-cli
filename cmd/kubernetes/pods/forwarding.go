package pods

import (
	"github.com/angelokurtis/kts-cli/cmd/common"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/gookit/color"
	"github.com/spf13/cobra"
	"strconv"
)

func forwarding(cmd *cobra.Command, args []string) {
	pods, err := kubectl.ListAllPods()
	if err != nil {
		common.Exit(err)
	}
	pod, err := pods.SelectOne()
	if err != nil {
		common.Exit(err)
	}
	port, err := pod.SelectContainerPort()
	if err != nil {
		common.Exit(err)
	}
	n := pod.Metadata.Name
	ns := pod.Metadata.Namespace
	p := strconv.Itoa(port)
	color.Comment.Println("kubectl port-forward " + n + " " + p + ":" + p + " -n " + ns)
}
