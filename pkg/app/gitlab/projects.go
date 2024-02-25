package gitlab

import (
	"fmt"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	gitlab "github.com/xanzy/go-gitlab"
)

func SearchProjects(search string) (projects []*gitlab.Project, total int, err error) {
	svc := client.Projects
	orderBy := "last_activity_at"

	projects, res, err := svc.ListProjects(&gitlab.ListProjectsOptions{
		OrderBy: &orderBy,
		Search:  &search,
	})
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}

	return projects, res.TotalItems, nil
}

func SelectOneProject(search string) (projects *gitlab.Project, err error) {
	p, total, err := SearchProjects(search)
	if err != nil {
		return nil, err
	}

	plus := total - len(p)
	if plus > 0 {
		return nil, errors.New(fmt.Sprintf("be more specific. %d projects were found", total))
	}

	list := projectList{items: p}

	return list.selectOne()
}

type projectList struct {
	items []*gitlab.Project
}

func (p *projectList) get(name string) *gitlab.Project {
	for _, item := range p.items {
		if item.PathWithNamespace == name {
			return item
		}
	}

	return nil
}

func (p *projectList) names() []string {
	names := make([]string, 0, 0)
	for _, item := range p.items {
		names = append(names, item.PathWithNamespace)
	}

	return names
}

func (p *projectList) selectOne() (*gitlab.Project, error) {
	if len(p.items) == 1 {
		return p.items[0], nil
	} else if len(p.items) == 0 {
		return nil, errors.New("no project found")
	}

	prompt := &survey.Select{
		Message: "Select the GitLab project:",
		Options: p.names(),
	}

	answer := ""

	err := survey.AskOne(prompt, &answer, survey.WithKeepFilter(true))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return p.get(answer), nil
}
