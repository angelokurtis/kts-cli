package dockerhub

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
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
	Count    int64  `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []*Tag `json:"results"`
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

	var (
		allTags []*Tag
		total   int64
		url     = baseURL + "/v2/repositories/" + repository + "/tags/?page_size=100"
	)

	for url != "" {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to create HTTP request for URL %q: %w", url, err)
		}

		req.Header.Set("Content-Type", "application/json")

		res, err := c.client.Do(req)
		if err != nil {
			return nil, 0, fmt.Errorf("HTTP request to %q failed: %w", url, err)
		}

		body, err := io.ReadAll(res.Body)
		_ = res.Body.Close()

		if err != nil {
			return nil, 0, fmt.Errorf("failed to read HTTP response body from %q: %w", url, err)
		}

		if res.StatusCode != http.StatusOK {
			return nil, 0, fmt.Errorf("unexpected HTTP status from %q: %d - %s", url, res.StatusCode, string(body))
		}

		var tagsResp Tags
		if err := json.Unmarshal(body, &tagsResp); err != nil {
			return nil, 0, fmt.Errorf("failed to parse JSON response from %q: %w", url, err)
		}

		allTags = append(allTags, tagsResp.Results...)
		total = tagsResp.Count
		url = tagsResp.Next
	}

	return allTags, total, nil
}
