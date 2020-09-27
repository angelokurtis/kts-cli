package terraform

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/angelokurtis/kts-cli/pkg/bash"
	"github.com/pkg/errors"
	"strings"
)

func ListResources() (Resources, error) {
	out, err := bash.RunAndLogRead("terraform state list")
	if err != nil {
		return nil, err
	}
	r := make([]string, 0, 0)
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		r = append(r, scanner.Text())
	}
	return r, nil
}

type Resources []string

func (r Resources) SelectMany() (Resources, error) {
	var selects []string
	prompt := &survey.MultiSelect{
		Message: "Select the Terraform resources:",
		Options: r,
	}

	err := survey.AskOne(prompt, &selects, survey.WithPageSize(10))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return selects, nil
}

func (r Resources) ApplyCommand() string {
	var b strings.Builder
	for _, s := range r {
		_, _ = fmt.Fprintf(&b, " -target=%s", s)
	}
	return "terraform apply" + b.String()
}

func (r Resources) DestroyCommand() string {
	var b strings.Builder
	for _, s := range r {
		_, _ = fmt.Fprintf(&b, " -target=%s", s)
	}
	return "terraform destroy" + b.String()
}
