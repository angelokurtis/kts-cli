package gitlab

import (
	"net/url"
	"os"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

func SearchMergeRequestsByUser(username string) (MergeRequests, error) {
	var u *gitlab.User
	var err error
	if username == "" {
		u, err = SelectOneUser()
	} else {
		u, err = SearchUser(username)
	}
	if err != nil {
		return nil, err
	}
	mr1, err := searchMergeRequestsByAssigneeID(u.ID)
	if err != nil {
		return nil, err
	}
	mr2, err := searchMergeRequestsByAuthorID(u.ID)
	if err != nil {
		return nil, err
	}
	return mr1.Join(mr2), nil
}

func searchMergeRequestsByAuthorID(authorID int) (MergeRequests, error) {
	svc := client.MergeRequests
	orderBy := "created_at"
	scope := "all"
	state := "opened"
	mr, _, err := svc.ListMergeRequests(&gitlab.ListMergeRequestsOptions{
		ListOptions: gitlab.ListOptions{
			Page:    0,
			PerPage: 100000,
		},
		OrderBy:  &orderBy,
		Scope:    &scope,
		AuthorID: &authorID,
		State:    &state,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return mr, nil
}

func searchMergeRequestsByAssigneeID(assigneeID int) (MergeRequests, error) {
	svc := client.MergeRequests
	orderBy := "created_at"
	scope := "all"
	state := "opened"
	mr, _, err := svc.ListMergeRequests(&gitlab.ListMergeRequestsOptions{
		ListOptions: gitlab.ListOptions{
			Page:    0,
			PerPage: 100000,
		},
		State:      &state,
		OrderBy:    &orderBy,
		View:       nil,
		Scope:      &scope,
		AssigneeID: &assigneeID,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return mr, nil
}

type MergeRequest gitlab.MergeRequest

type MergeRequests []*gitlab.MergeRequest

func (m MergeRequest) Path() (string, error) {
	baseURL, err := url.Parse(os.Getenv("CB_GITLAB_BASE_URL"))
	if err != nil {
		return "", err
	}
	return strings.Split(m.WebURL, baseURL.Host)[1], nil
}

func (m MergeRequest) Assign() error {
	user, err := SelectOneUser()
	if err != nil {
		return err
	}
	svc := client.MergeRequests
	_, _, err = svc.UpdateMergeRequest(m.ProjectID, m.IID, &gitlab.UpdateMergeRequestOptions{AssigneeID: &user.ID})
	return errors.WithStack(err)
}

func (m *MergeRequests) SelectOne() (*MergeRequest, error) {
	paths, err := m.Paths()
	if err != nil {
		return nil, err
	}

	if len(paths) == 1 {
		return m.Get(paths[0])
	}

	var selected string
	prompt := &survey.Select{
		Message: "Select the MergeRequest:",
		Options: paths,
	}

	err = survey.AskOne(prompt, &selected, survey.WithPageSize(10), survey.WithKeepFilter(true))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return m.Get(selected)
}

func (m MergeRequests) Get(path string) (*MergeRequest, error) {
	for _, mr := range m {
		mr := (*MergeRequest)(mr)
		p, err := mr.Path()
		if err != nil {
			return nil, err
		}
		if p == path {
			return mr, nil
		}
	}
	return nil, nil
}

func (m MergeRequests) Paths() ([]string, error) {
	paths := make([]string, 0, len(m))
	for _, mr := range m {
		mr := (*MergeRequest)(mr)
		path, err := mr.Path()
		if err != nil {
			return nil, err
		}
		paths = append(paths, path)
	}
	return paths, nil
}

func (m MergeRequests) Join(mrs MergeRequests) MergeRequests {
	res := make([]*gitlab.MergeRequest, 0, 0)
	ids := make([]int, 0, 0)
	for _, mr := range m {
		res = append(res, mr)
		ids = append(ids, mr.ID)
	}
	for _, mr := range mrs {
		if !contains(ids, mr.ID) {
			res = append(res, mr)
		}
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].CreatedAt.After(*res[j].CreatedAt)
	})
	return res
}

func contains(slice []int, item int) bool {
	set := make(map[int]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}
