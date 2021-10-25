package dockerhub

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func UnmarshalTags(data []byte) (Tags, error) {
	var r Tags
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Tags) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Tags struct {
	Count    int64       `json:"count"`
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []*Tag      `json:"results"`
}

type Tag struct {
	Creator             int64               `json:"creator"`
	ID                  int64               `json:"id"`
	ImageID             interface{}         `json:"image_id"`
	Images              []Image             `json:"images"`
	LastUpdated         time.Time           `json:"last_updated"`
	LastUpdater         int64               `json:"last_updater"`
	LastUpdaterUsername LastUpdaterUsername `json:"last_updater_username"`
	Name                string              `json:"name"`
	Repository          int64               `json:"repository"`
	FullSize            int64               `json:"full_size"`
	V2                  bool                `json:"v2"`
	TagStatus           Status              `json:"tag_status"`
	TagLastPulled       string              `json:"tag_last_pulled"`
	TagLastPushed       string              `json:"tag_last_pushed"`
}

type Image struct {
	Architecture string      `json:"architecture"`
	Features     string      `json:"features"`
	Variant      interface{} `json:"variant"`
	Digest       string      `json:"digest"`
	OS           OS          `json:"os"`
	OSFeatures   string      `json:"os_features"`
	OSVersion    interface{} `json:"os_version"`
	Size         int64       `json:"size"`
	Status       Status      `json:"status"`
	LastPulled   string      `json:"last_pulled"`
	LastPushed   *time.Time  `json:"last_pushed"`
}

type (
	OS                  string
	Status              string
	LastUpdaterUsername string
)

func (c *Client) ListTags(repository string) ([]*Tag, int64, error) {
	if !strings.Contains(repository, "/") {
		repository = "library/" + repository
	}
	url := baseURL + "/v2/repositories/" + repository + "/tags/?page_size=1000000&page=1"
	method := "GET"

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}

	repositories, err := UnmarshalTags(body)
	if err != nil {
		return nil, 0, err
	}

	return repositories.Results, repositories.Count, nil
}
