package yaml

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func join(cmd *cobra.Command, args []string) {
	var b strings.Builder
	err := filepath.Walk(args[0], func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}
		if filepath.Ext(path) == ".yml" || filepath.Ext(path) == ".yaml" {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatal(err)
			}
			_, err = fmt.Fprintf(&b, "---\n%s\n", content)
			if err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	if err = ioutil.WriteFile("manifests.yaml", []byte(b.String()), 0o644); err != nil {
		log.Fatal(err)
	}
}
