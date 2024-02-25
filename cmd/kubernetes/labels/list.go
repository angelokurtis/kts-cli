package labels

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
)

// kube labels list
func list(cmd *cobra.Command, args []string) {
	kinds := ""

	if len(args) == 0 {
		rd, err := kubectl.ListResourceDefinitions()
		if err != nil {
			system.Exit(err)
		}

		rd = rd.FilterVerbs("list")
		if !allNamespaces {
			rd = rd.FilterNamespaced()
		}

		if group != "" {
			rd = rd.FilterAPIGroup(group)
		}

		kinds = rd.String()
	} else {
		kinds = args[0]
	}

	labels, err := kubectl.ListLabels(kinds, namespace, allNamespaces)
	if err != nil {
		system.Exit(err)
	}

	for _, label := range labels {
		fmt.Println(label)
	}
}
