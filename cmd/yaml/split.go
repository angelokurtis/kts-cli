package yaml

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"strings"

	changeCase "github.com/ku/go-change-case"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"

	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func split(cmd *cobra.Command, args []string) {
	filename := args[0]

	out, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(errors.WithStack(err))
	}

	file := make([]string, 0, 0)
	files := make([]string, 0, 0)

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" && len(file) > 0 {
			files = append(files, strings.Join(file, "\n"))
			file = make([]string, 0, 0)
		} else {
			file = append(file, line)
		}
	}

	if len(file) > 0 {
		files = append(files, strings.Join(file, "\n"))
		file = make([]string, 0, 0)
	}

	resources := make([]*Resource, 0, 0)

	for _, f := range files {
		r := &Resource{}

		err = yaml.Unmarshal([]byte(f), r)
		if err != nil {
			log.Fatal(errors.WithStack(err))
		}

		resources = append(resources, r)
	}

	directory := strings.Replace(filename, ".yaml", "", -1)
	directory = strings.Replace(directory, ".yml", "", -1)

	err = SaveMany(files, directory)
	if err != nil {
		log.Fatal(err)
	}
}

func SaveMany(files []string, directory string) error {
	for _, file := range files {
		err := SaveOne(file, directory)
		if err != nil {
			return err
		}
	}

	return nil
}

func SaveOne(file, directory string) error {
	_, err := bash.Run("mkdir -p " + directory)
	if err != nil {
		return err
	}

	r := &Resource{}
	out := []byte(file)

	err = yaml.Unmarshal(out, r)
	if err != nil {
		log.Fatal(errors.WithStack(err))
	}

	filename := directory + "/" + r.Metadata.Name + "." + changeCase.Param(r.Kind) + ".yaml"
	if err = ioutil.WriteFile(filename, out, 0o644); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

type Resource struct {
	Kind       string `yaml:"kind"`
	APIVersion string `yaml:"apiVersion"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
}
