package terraformer

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/angelokurtis/kts-cli/pkg/app/gcloud"
	"github.com/angelokurtis/kts-cli/pkg/bash"
	"github.com/pkg/errors"
	"strings"
)

func ListResources(provider string) (*Resources, error) {
	cmd := fmt.Sprintf("terraformer import %s list", provider)
	if provider == "google" {
		project, err := gcloud.CurrentProject()
		if err != nil {
			return nil, err
		}
		region, err := gcloud.CurrentRegion()
		if err != nil {
			return nil, err
		}
		cmd = fmt.Sprintf("%s --projects=%s --regions=%s", cmd, project.ID, region)
	}
	out, err := bash.RunAndLogRead(cmd)
	if err != nil {
		return nil, err
	}
	_ = out
	scanner := bufio.NewScanner(bytes.NewReader(out))
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	resources := make([]string, 0, 0)
	for scanner.Scan() {
		resources = append(resources, scanner.Text())
	}
	return &Resources{
		Provider: provider,
		Items:    resources,
	}, nil
}

type Resources struct {
	Provider string
	Items    []string
}

func (r *Resources) SelectMany() (*Resources, error) {
	items := r.Items
	var selects []string
	d := make([]int, 0, len(items))
	for i, _ := range items {
		d = append(d, i)
	}
	prompt := &survey.MultiSelect{
		Renderer: survey.Renderer{},
		Message:  "Select the Terraformer Resources:",
		Options:  items,
		//Default:  d,
	}

	err := survey.AskOne(prompt, &selects, survey.WithPageSize(10))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Resources{
		Provider: r.Provider,
		Items:    selects,
	}, nil
}

func (r *Resources) Import() error {
	provider := r.Provider
	cmd := fmt.Sprintf("terraformer import %s --resources=%s", provider, strings.Join(r.Items, ","))
	if provider == "google" {
		project, err := gcloud.CurrentProject()
		if err != nil {
			return err
		}
		region, err := gcloud.CurrentRegion()
		if err != nil {
			return err
		}
		cmd = fmt.Sprintf("%s --projects=%s --regions=%s", cmd, project.ID, region)
	}
	_, err := bash.RunAndLogWrite(cmd)
	if err != nil {
		return err
	}
	return nil
}
