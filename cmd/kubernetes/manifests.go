package kubernetes

import (
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
	"github.com/spf13/cobra"
)

// kube manifests -A ClusterRoleBinding
func manifests(cmd *cobra.Command, args []string) {
	resources := ""
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
		resources = rd.String()
	} else {
		resources = args[0]
	}
	results, err := kubectl.SelectResources(resources, namespace, allNamespaces)
	if err != nil {
		system.Exit(err)
	}
	err = kubectl.SaveResourcesManifests(results, status)
	if err != nil {
		system.Exit(err)
	}
}
