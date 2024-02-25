package m3u

import (
	"path/filepath"
	"strings"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
)

type Groups struct {
	Items    []string
	filename string
}

func (g *Groups) SelectMany() (*Groups, error) {
	if len(g.Items) == 0 {
		return nil, nil
	}

	prompt := &survey.MultiSelect{
		Message: "Select Groups:",
		Options: g.IDs(),
		Default: g.defaults(),
	}

	var selects []string

	err := survey.AskOne(prompt, &selects, survey.WithPageSize(20), survey.WithKeepFilter(true))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	groups := make([]string, 0, len(selects))

	for _, name := range selects {
		group := g.Get(name)
		if group != "" {
			groups = append(groups, group)
		}
	}

	return &Groups{
		Items:    groups,
		filename: g.filename,
	}, nil
}

func (g *Groups) IDs() []string {
	n := make([]string, 0, 0)
	for _, group := range g.Items {
		n = append(n, strings.ReplaceAll(group, " ", "."))
	}

	return n
}

func (g Groups) defaults() []string {
	filename := g.filename
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]

	channels, err := ListChannels(name + "[edited]" + ext)
	if err != nil {
		channels, err = ListChannels(filename)
		if err != nil {
			return nil
		}
	}

	return channels.Groups().IDs()
}

func (g *Groups) Get(id string) string {
	for _, group := range g.Items {
		if strings.ReplaceAll(group, " ", ".") == id {
			return group
		}
	}

	return ""
}
