package gitlab

import (
	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

func PipelineByMergeRequest(mr *gitlab.MergeRequest) (*gitlab.PipelineInfo, error) {
	svc := client.MergeRequests
	res, _, err := svc.GetMergeRequest(mr.ProjectID, mr.IID, &gitlab.GetMergeRequestsOptions{})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return res.Pipeline, nil
}
