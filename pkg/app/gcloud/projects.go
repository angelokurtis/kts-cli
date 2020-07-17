package gcloud

import (
	"encoding/json"
	"github.com/pkg/errors"
	"time"
)

func ListProjects() ([]*Project, error) {
	out, err := run("projects", "list")
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
