package pods

import (
	"github.com/angelokurtis/kts-cli/cmd/common"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/spf13/cobra"
)

func logs(cmd *cobra.Command, args []string) {
	pods, err := kubectl.ListAllPods()
	if err != nil {
		common.Exit(err)
	}
	selects, err := pods.SelectMany()
	if err != nil {
		common.Exit(err)
	}
	for _, pod := range selects {
		err := kubectl.SaveLogs(pod)
		if err != nil {
			common.Exit(err)
		}
	}
}
