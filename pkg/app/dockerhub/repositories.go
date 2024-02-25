package dockerhub

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

func UnmarshalRepositories(data []byte) (*Repositories, error) {
	var r *Repositories
	err := json.Unmarshal(data, &r)

	return r, errors.WithStack(err)
}

func (r *Repositories) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Repositories struct {
	Count    int64         `json:"count"`
	Next     string        `json:"next"`
	Previous interface{}   `json:"previous"`
	Results  []*Repository `json:"results"`
}

type Repository struct {
	User              string         `json:"user"`
	Name              string         `json:"name"`
	Namespace         string         `json:"namespace"`
	RepositoryType    RepositoryType `json:"repository_type"`
	Status            int64          `json:"status"`
	Description       string         `json:"description"`
	IsPrivate         bool           `json:"is_private"`
	IsAutomated       bool           `json:"is_automated"`
	CanEdit           bool           `json:"can_edit"`
	StarCount         int64          `json:"star_count"`
	PullCount         int64          `json:"pull_count"`
	LastUpdated       *time.Time     `json:"last_updated"`
	IsMigrated        bool           `json:"is_migrated"`
	CollaboratorCount int64          `json:"collaborator_count"`
	Affiliation       Affiliation    `json:"affiliation"`
	HubUser           string         `json:"hub_user"`
}

type Affiliation string

type RepositoryType string

func (c *Client) ListRepositories(hubuser string) ([]*Repository, error) {
	url := baseURL + "/v2/repositories/" + hubuser + "?page=1&page_size=1000000"
	method := "GET"

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	repositories, err := UnmarshalRepositories(body)
	if err != nil {
		return nil, err
	}

	return repositories.Results, nil
}
