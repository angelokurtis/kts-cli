package pods

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/cmd/common"
	"github.com/angelokurtis/kts-cli/internal/color"
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
		pn := pod.Metadata.Name
		ns := pod.Metadata.Namespace

		if len(pod.Spec.Containers) == 1 {
			fmt.Printf(color.Notice, fmt.Sprintf("kubectl logs %s -n %s > ./%s.log\n", pn, ns, pn))
			continue
		}
		for _, container := range pod.Spec.Containers {
			cn := container.Name
			fmt.Printf(color.Notice, fmt.Sprintf("kubectl logs %s -c %s -n %s> ./%s.%s.log\n", pn, cn, ns, pn, cn))
		}
	}
}
