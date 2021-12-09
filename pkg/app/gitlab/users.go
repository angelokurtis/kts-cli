package gitlab

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
	"github.com/xanzy/go-gitlab"
)

func GetMyUser() (*gitlab.User, error) {
	svc := client.Users
	u, _, err := svc.CurrentUser()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return u, nil
}

func SearchUsers() ([]*gitlab.User, error) {
	u := make([]*gitlab.User, 0, 0)
	lastPage := 0
	currentPage := 0
	svc := client.Users

	for currentPage <= lastPage {
		users, res, err := svc.ListUsers(&gitlab.ListUsersOptions{
			ListOptions: gitlab.ListOptions{
				Page:    currentPage,
				PerPage: 100,
			},
		})
		if err != nil {
			return nil, errors.WithStack(err)
		}
		u = append(u, users...)
		lastPage = res.TotalPages
		currentPage++
	}

	return u, nil
}

func SearchUser(username string) (*gitlab.User, error) {
	svc := client.Users
	users, _, err := svc.ListUsers(&gitlab.ListUsersOptions{
		Username: &username,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if len(users) == 0 {
		return nil, nil
	}
	return users[0], nil
}

func SelectOneUser() (*gitlab.User, error) {
	u, err := SearchUsers()
	if err != nil {
		return nil, err
	}

	list := userList{items: u}
	return list.selectOne()
}

type userList struct {
	items []*gitlab.User
}

func (u *userList) get(name string) *gitlab.User {
	for _, item := range u.items {
		if item.Username == name {
			return item
		}
	}
	return nil
}

func (u *userList) names() []string {
	names := make([]string, 0, 0)
	for _, item := range u.items {
		names = append(names, item.Username)
	}
	return names
}

func (u *userList) selectOne() (*gitlab.User, error) {
	if len(u.items) == 1 {
		return u.items[0], nil
	} else if len(u.items) == 0 {
		return nil, errors.New("no project found")
	}

	prompt := &survey.Select{
		Message: "Select the GitLab user:",
		Options: u.names(),
	}

	answer := ""
	err := survey.AskOne(prompt, &answer, survey.WithKeepFilter(true))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return u.get(answer), nil
}
