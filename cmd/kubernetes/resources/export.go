package resources

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/cmd/common"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/spf13/cobra"
)

func export(cmd *cobra.Command, args []string) {
	resources, err := kubectl.ListResources()
	if err != nil {
		common.Exit(err)
	}
	for _, resource := range resources {
		fmt.Printf(resource)
	}
}
