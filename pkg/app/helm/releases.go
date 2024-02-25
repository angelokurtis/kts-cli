package helm

import (
	"encoding/json"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"

	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func ListReleases(options ...OptionFunc) (Releases, error) {
	o := new(Option)
	if err := o.apply(options...); err != nil {
		return nil, err
	}

	cmd := "helm list -a -o=json"
	if o.AllNamespaces {
		cmd += " -A"
	} else if o.Namespace != "" {
		cmd += " -n " + o.Namespace
	}

	out, err := bash.RunAndLogRead(cmd)
	if err != nil {
		return nil, err
	}

	var releases Releases
	if err := json.Unmarshal(out, &releases); err != nil {
		return nil, errors.WithStack(err)
	}

	return releases, nil
}

type Releases []*Release

func (r Releases) IDs() []string {
	ids := make([]string, 0, len(r))
	for _, release := range r {
		ids = append(ids, release.Namespace+"/"+release.Name)
	}

	return ids
}

func (r Releases) Get(id string) *Release {
	for _, release := range r {
		if release.Namespace+"/"+release.Name == id {
			return release
		}
	}

	return nil
}

func (r Releases) SelectOne() (*Release, error) {
	ids := r.IDs()

	if len(ids) == 1 {
		return r.Get(ids[0]), nil
	}

	var selected string

	prompt := &survey.Select{
		Message: "Select the release:",
		Options: ids,
	}

	err := survey.AskOne(prompt, &selected, survey.WithPageSize(10), survey.WithKeepFilter(true))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return r.Get(selected), nil
}

type Release struct {
	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Revision   string `json:"revision"`
	Updated    string `json:"updated"`
	Status     string `json:"status"`
	Chart      string `json:"chart"`
	AppVersion string `json:"app_version"`
}
