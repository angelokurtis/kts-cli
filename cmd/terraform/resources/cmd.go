package resources

import (
	"fmt"
	"io/ioutil"
	log "log/slog"
	"os"
	"path/filepath"
	"strings"

	changeCase "github.com/ku/go-change-case"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/colors"
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/terraform"
	"github.com/angelokurtis/kts-cli/pkg/bash"
)

var (
	provider = ""
	filename = ""
	Command  = &cobra.Command{
		Use:   "resources",
		Short: "Utility functions to deal with Resources in the Terraform Registry",
		Run:   system.Help,
	}
)

func init() {
	Command.PersistentFlags().StringVarP(&provider, "provider", "p", "", "")
	Command.AddCommand(&cobra.Command{Use: "import", Run: importCmd})
	convertCMD := &cobra.Command{Use: "convert", Run: convert}
	convertCMD.PersistentFlags().StringVarP(&filename, "filename", "f", "", "")
	Command.AddCommand(convertCMD)
}

func convert(cmd *cobra.Command, args []string) {
	var files []string

	err := filepath.Walk(filename, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml") {
			out, err := terraform.YAMLDecode(file)
			if err != nil {
				log.Error(err.Error())
				return
			}

			filename := strings.Replace(file, ".yaml", ".tf", -1)
			filename = strings.Replace(filename, ".yml", ".tf", -1)

			if err = ioutil.WriteFile(filename, out, 0o644); err != nil {
				log.Error(err.Error())
				return
			}
		}
	}
}

func importCmd(cmd *cobra.Command, args []string) {
	if provider == "" {
		p, err := terraform.SelectProvider()
		if err != nil {
			log.Error(err.Error())
			return
		}

		provider = p.Name
	}

	resource, err := terraform.SelectResource(provider)
	if err != nil {
		log.Error(err.Error())
		return
	}

	out, err := resource.Encode()
	if err != nil {
		log.Error(err.Error())
		return
	}

	filename := fmt.Sprintf("%s.tf", changeCase.Param(resource.Name))

	err = ioutil.WriteFile(filename, out, os.ModePerm)
	if err != nil {
		log.Error(err.Error())
		return
	}

	err = resource.Import()
	if err != nil {
		log.Error(err.Error())
		return
	}

	state, err := bash.Run("terraform state show " + resource.GetID())
	if err != nil {
		log.Error(err.Error())
		return
	}

	err = ioutil.WriteFile(filename, colors.Strip(state), os.ModePerm)
	if err != nil {
		log.Error(err.Error())
		return
	}
}
