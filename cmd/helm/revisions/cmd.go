package revisions

import (
	"bytes"
	"compress/gzip"
	b64 "encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	log "log/slog"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/itchyny/json2yaml"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/pkg/app/helm"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
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

	releases, err := helm.ListReleases(opt...)
	dieOnErr(err)

	release, err := releases.SelectOne()
	dieOnErr(err)

	opt = append(opt, helm.OnNamespace(release.Namespace))

	history, err := helm.GetHistory(release.Name, opt...)
	dieOnErr(err)

	history, err = history.SelectMany()
	dieOnErr(err)

	chartMetadata, chartValues := getChartMetadata(release)
	saveChart(release, chartMetadata, chartValues)

	for _, revision := range history {
		_, err = bash.Run(fmt.Sprintf("mkdir -p helm-releases/%s/%d", release.Name, revision.Number))
		dieOnErr(err)

		manifests, err := helm.GetManifests(release.Name, revision.Number, opt...)
		dieOnErr(err)

		err = ioutil.WriteFile(fmt.Sprintf("helm-releases/%s/%d/manifests.yaml", release.Name, revision.Number), manifests, 0o644)
		dieOnErr(err)

		svalues, err := helm.GetSuppliedValues(release.Name, revision.Number, opt...)
		dieOnErr(err)

		err = ioutil.WriteFile(fmt.Sprintf("helm-releases/%s/%d/values.supplied.yaml", release.Name, revision.Number), svalues, 0o644)
		dieOnErr(err)

		cvalues, err := helm.GetComputedValues(release.Name, revision.Number, opt...)
		dieOnErr(err)

		err = ioutil.WriteFile(fmt.Sprintf("helm-releases/%s/%d/values.computed.yaml", release.Name, revision.Number), cvalues, 0o644)
		dieOnErr(err)

		notes, err := helm.GetNotes(release.Name, revision.Number, opt...)
		dieOnErr(err)

		err = ioutil.WriteFile(fmt.Sprintf("helm-releases/%s/%d/notes.txt", release.Name, revision.Number), notes, 0o644)
		dieOnErr(err)
	}
}

func getChartMetadata(release *helm.Release) ([]byte, []byte) {
	secretVal, err := kubectl.GetSecretKeyValue(&kubectl.KeyRef{Name: fmt.Sprintf("sh.helm.release.v1.%s.v1", release.Name), Key: "release"}, release.Namespace)
	dieOnErr(err)

	compressed, err := b64.StdEncoding.DecodeString(secretVal)
	dieOnErr(err)

	gzipReader, err := gzip.NewReader(bytes.NewReader(compressed))
	dieOnErr(err)

	decompressed, err := io.ReadAll(gzipReader)
	dieOnErr(err)

	chartMeta, _, _, err := jsonparser.Get(decompressed, "chart", "metadata")
	dieOnErr(err)

	chartValues, _, _, err := jsonparser.Get(decompressed, "chart", "values")
	dieOnErr(err)

	return chartMeta, chartValues
}

func saveChart(release *helm.Release, metadata, values []byte) {
	metaReader := strings.NewReader(string(metadata))
	var metaYAML strings.Builder
	err := json2yaml.Convert(&metaYAML, metaReader)
	dieOnErr(err)

	_, err = bash.Run(fmt.Sprintf("mkdir -p helm-releases/%s", release.Name))
	dieOnErr(err)

	err = ioutil.WriteFile(fmt.Sprintf("helm-releases/%s/Chart.yaml", release.Name), []byte(metaYAML.String()), 0o644)
	dieOnErr(err)

	valuesReader := strings.NewReader(string(values))
	var valuesYAML strings.Builder
	err = json2yaml.Convert(&valuesYAML, valuesReader)
	dieOnErr(err)

	err = ioutil.WriteFile(fmt.Sprintf("helm-releases/%s/values.yaml", release.Name), []byte(valuesYAML.String()), 0o644)
	dieOnErr(err)
}

func init() {
	Command.PersistentFlags().StringVarP(&namespace, "namespace", "n", "", "Set the namespace for a current request")
	Command.PersistentFlags().BoolVarP(&allNamespaces, "all-namespaces", "A", false, "If present, resources the requested object(s) across all namespaces. Namespace in current\ncontext is ignored even if specified with --namespace.")
}

func dieOnErr(err error) {
	if err != nil {
		log.Error(err.Error())
		panic(nil)
	}
}
