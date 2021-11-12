package revisions

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/pkg/app/helm"
)

var (
	allNamespaces = false
	namespace     = ""
	Command       = &cobra.Command{Use: "revisions", Run: revisions}
)

// kts helm revisions
func revisions(cmd *cobra.Command, args []string) {
	opt := []helm.OptionFunc{helm.OnNamespace(namespace)}
	if allNamespaces {
		opt = append(opt, helm.OnAnyNamespace())
	}
	release := ""
	if len(args) > 0 {
		release = args[0]
	} else {
		releases, err := helm.ListReleases(opt...)
		dieOnErr(err)

		r, err := releases.SelectOne()
		dieOnErr(err)
		release = r.Name
		opt = append(opt, helm.OnNamespace(r.Namespace))
	}
	history, err := helm.GetHistory(release, opt...)
	dieOnErr(err)

	history, err = history.SelectMany()
	dieOnErr(err)

	for _, revision := range history {
		fmt.Println(revision.Number)
	}
}

func init() {
	Command.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Set the namespace for a current request")
	Command.PersistentFlags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "If present, resources the requested object(s) across all namespaces. Namespace in current\ncontext is ignored even if specified with --namespace.")
}

func dieOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
