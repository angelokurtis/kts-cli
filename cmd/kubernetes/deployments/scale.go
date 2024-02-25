package deployments

import (
	"strconv"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/angelokurtis/kts-cli/internal/log"
	"github.com/angelokurtis/kts-cli/pkg/app/kubectl"
)

func scale(cmd *cobra.Command, args []string) {
	deploys, err := kubectl.ListDeployments(namespace, allNamespaces)
	if err != nil {
		log.Fatal(err)
	}

	deploys, err = deploys.SelectMany()
	if err != nil {
		log.Fatal(err)
	}

	answers := struct {
		Replicas int
	}{}

	err = survey.Ask([]*survey.Question{{
		Name:   "replicas",
		Prompt: &survey.Input{Message: "How many replicas?"},
		Validate: func(val any) error {
			switch v := val.(type) {
			case int, float64:
				return nil
			case string:
				if _, err = strconv.Atoi(v); err != nil {
					return errors.New("it should be a number")
				}
				return nil
			default:
				return errors.New("it should be a number")
			}
		},
		Transform: survey.Title,
	}}, &answers)
	if err != nil {
		log.Fatal(err)
	}

	err = deploys.Scale(answers.Replicas)
	if err != nil {
		log.Fatal(err)
	}
}
