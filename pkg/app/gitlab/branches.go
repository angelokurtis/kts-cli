package gitlab

import (
	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

func ListBranches(project *gitlab.Project) ([]*gitlab.Branch, error) {
	svc := client.Branches
	branches, _, err := svc.ListBranches(project.ID, &gitlab.ListBranchesOptions{})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return branches, nil
}
