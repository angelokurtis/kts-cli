package resources

import (
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/system"
)

var (
	allNamespaces = false
	owners        = false
	noOwners      = false
	group         = ""
	namespace     = ""
	Command       = &cobra.Command{
		Use:   "resources",
		Short: "Utility functions to deal with any type of resource",
		Run:   system.Help,
	}
)

func init() {
	listCMD := &cobra.Command{Use: "list", Run: run(list)}
	listCMD.PersistentFlags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "List the requested object(s) across all namespaces, ignoring the current context's namespace.")
	listCMD.PersistentFlags().BoolVarP(&owners, "owners", "O", false, "Filter and list only the object(s) that own other objects.")
	listCMD.PersistentFlags().BoolVar(&noOwners, "no-owners", false, "Filter and list only the object(s) that do not have an owner.")
	listCMD.PersistentFlags().StringVar(&group, "group", "", "Specify the API group of the resource(s) to list.")
	listCMD.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Specify the namespace to list the resource(s) from.")
	Command.AddCommand(listCMD)
}

func run(fn func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := fn(cmd, args); err != nil {
			system.Exit(err)
		}
	}
}
