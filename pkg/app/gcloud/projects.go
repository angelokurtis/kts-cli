package gcloud

import (
	"encoding/json"
	"github.com/angelokurtis/kts-cli/pkg/bash"
	"github.com/pkg/errors"
	"strings"
	"time"
)

func CurrentProject() (*Project, error) {
	out, err := bash.RunAndLogRead("gcloud config get-value project")
	if err != nil {
		return nil, err
	}
	project, err := DescribeProject(strings.TrimSpace(string(out)))
	if err != nil {
		return nil, err
	}
	return project, nil
}

func DescribeProject(projectId string) (*Project, error) {
	out, err := bash.RunAndLogRead("gcloud projects describe " + projectId + " --format=json")
	if err != nil {
		return nil, err
	}

	var project *Project
	if err := json.Unmarshal(out, &project); err != nil {
		return nil, errors.WithStack(err)
	}

	return project, nil
}

func ListProjects() ([]*Project, error) {
	out, err := runAndLogRead("projects", "list")
	if err != nil {
		return nil, err
	}

	var projects []*Project
	if err := json.Unmarshal(out, &projects); err != nil {
		return nil, errors.WithStack(err)
	}

	return projects, nil
}

func ListProjectNames() ([]string, error) {
	projects, err := ListProjects()
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(projects))
	for _, project := range projects {
		names = append(names, project.Name)
	}
	return names, nil
}

type Project struct {
	CreateTime     time.Time `json:"createTime"`
	LifecycleState string    `json:"lifecycleState"`
	Name           string    `json:"name"`
	Parent         struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"parent"`
	ID     string `json:"projectId"`
	Number string `json:"projectNumber"`
}
