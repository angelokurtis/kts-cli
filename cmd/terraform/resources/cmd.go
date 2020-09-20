package resources

import (
	"fmt"
	"github.com/angelokurtis/kts-cli/internal/colors"
	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/internal/system"
	"github.com/angelokurtis/kts-cli/pkg/app/terraform"
	"github.com/angelokurtis/kts-cli/pkg/bash"
	changeCase "github.com/ku/go-change-case"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
				log.Fatal(err)
			}
			filename := strings.Replace(file, ".yaml", ".tf", -1)
			filename = strings.Replace(filename, ".yml", ".tf", -1)
			if err = ioutil.WriteFile(filename, out, 0644); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func importCmd(cmd *cobra.Command, args []string) {
	if provider == "" {
		p, err := terraform.SelectProvider()
		if err != nil {
			log.Fatal(err)
		}
		provider = p.Name
	}

	resource, err := terraform.SelectResource(provider)
	if err != nil {
		log.Fatal(err)
	}
	out, err := resource.Encode()
	if err != nil {
		log.Fatal(err)
	}
	filename := fmt.Sprintf("%s.tf", changeCase.Param(resource.Name))
	err = ioutil.WriteFile(filename, out, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	err = resource.Import()
	if err != nil {
		log.Fatal(err)
	}
	state, err := bash.Run("terraform state show " + resource.GetID())
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(filename, colors.Strip(state), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}
