package revisions

import (
	"fmt"
	"io/ioutil"
	log "log/slog"

	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/pkg/app/helm"
	"github.com/angelokurtis/kts-cli/pkg/bash"
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
		_, err = bash.Run(fmt.Sprintf("mkdir -p helm-releases/%s/%d", release, revision.Number))
		dieOnErr(err)

		manifests, err := helm.GetManifests(release, revision.Number, opt...)
		dieOnErr(err)

		err = ioutil.WriteFile(fmt.Sprintf("helm-releases/%s/%d/manifests.yaml", release, revision.Number), manifests, 0o644)
		dieOnErr(err)

		svalues, err := helm.GetSuppliedValues(release, revision.Number, opt...)
		dieOnErr(err)

		err = ioutil.WriteFile(fmt.Sprintf("helm-releases/%s/%d/values.supplied.yaml", release, revision.Number), svalues, 0o644)
		dieOnErr(err)

		cvalues, err := helm.GetComputedValues(release, revision.Number, opt...)
		dieOnErr(err)

		err = ioutil.WriteFile(fmt.Sprintf("helm-releases/%s/%d/values.computed.yaml", release, revision.Number), cvalues, 0o644)
		dieOnErr(err)

		notes, err := helm.GetNotes(release, revision.Number, opt...)
		dieOnErr(err)

		err = ioutil.WriteFile(fmt.Sprintf("helm-releases/%s/%d/notes.txt", release, revision.Number), notes, 0o644)
		dieOnErr(err)
	}
}

func init() {
	Command.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Set the namespace for a current request")
	Command.PersistentFlags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "If present, resources the requested object(s) across all namespaces. Namespace in current\ncontext is ignored even if specified with --namespace.")
}

func dieOnErr(err error) {
	if err != nil {
		log.Error(err.Error())
		return
	}
}
