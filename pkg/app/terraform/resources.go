package terraform

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
)

func ListResources(provider string) []string {
	switch provider {
	case "google":
		return []string{
			"google_container_cluster",
			"google_container_node_pool",
		}
	case "aws":
		return nil
	}
	return nil
}

func SelectResource(provider string) (*Resource, error) {
	r := ListResources(provider)

	var selected string
	if len(r) == 0 {
		return nil, nil
	} else if len(r) > 1 {
		prompt := &survey.Select{
			Message: "Select the Terraform Resource:",
			Options: r,
		}

		err := survey.AskOne(prompt, &selected, survey.WithPageSize(10))
		if err != nil {
			return nil, errors.WithStack(err)
		}
	} else {
		selected = r[0]
	}

	return NewResource(selected)
}
