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
	release := ""
	if len(args) > 0 {
		release = args[0]
	} else {
		opt := []helm.ListReleasesOptionFunc{helm.OnNamespace(namespace)}
		if allNamespaces {
			opt = append(opt, helm.OnAnyNamespace())
		}
		releases, err := helm.ListReleases(opt...)
		dieOnErr(err)

		r, err := releases.SelectOne()
		dieOnErr(err)
		release = r.Name
		namespace = r.Namespace
	}
	fmt.Println(release + " -n " + namespace)
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
