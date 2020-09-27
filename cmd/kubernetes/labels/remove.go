package labels

import (
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/spf13/cobra"
)

// kube labels remove
func remove(cmd *cobra.Command, args []string) {
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
	label, err := kubectl.SelectLabel(kinds, namespace, allNamespaces)
	if err != nil {
		system.Exit(err)
	}
	err = kubectl.RemoveLabels(kinds, label, namespace, allNamespaces)
	if err != nil {
		system.Exit(err)
	}
}
