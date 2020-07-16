package resources

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/cmd/common"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/spf13/cobra"
)

func list(cmd *cobra.Command, args []string) {
	resources := ""
	if len(args) == 0 {
		rd, err := kubectl.ListResourceDefinitions()
		if err != nil {
			common.Exit(err)
		}
		rd = rd.FilterVerbs("get")
		if !allNamespaces {
			rd = rd.FilterNamespaced()
		}
		if group != "" {
			rd = rd.FilterAPIGroup(group)
		}
		resources = rd.String()
	} else {
		resources = args[0]
	}
	results, err := kubectl.ListResources(resources, allNamespaces)
	if err != nil {
		common.Exit(err)
	}
	for _, result := range results {
		fmt.Println(result)
	}
}
