package gcloud

import (
	"encoding/json"
	"github.com/AlecAivazis/survey/v2"
)

func SelectContainerRepositories() ([]string, error) {
	repositories, err := ListContainerRepositories()
	if err != nil {
		return nil, err
	}

	var selects []string
	prompt := &survey.MultiSelect{
		Message: "Select the images repositories:",
		Options: repositories,
	}

	err = survey.AskOne(prompt, &selects, survey.WithPageSize(25))
	if err != nil {
		return nil, err
	}

	return selects, nil
}

func ListContainerRepositories() ([]string, error) {
	out, err := run("container", "images", "list")
	if err != nil {
		return nil, err
	}

	var decoded []map[string]string
	if err := json.Unmarshal(out, &decoded); err != nil {
		return nil, err
	}

	images := make([]string, 0, len(decoded))
	for _, v := range decoded {
		images = append(images, v["name"])
	}
	return images, nil
}

func ListContainerImages(repository string) ([]*ContainerImage, error) {
	out, err := run("container", "images", "list-tags", repository)
	if err != nil {
		return nil, err
	}
	var tags []*ContainerImage
	if err := json.Unmarshal(out, &tags); err != nil {
		return nil, err
	}
	for _, tag := range tags {
		tag.Repository = repository
		tag.FullyQualifiedDigest = repository + "@" + tag.Digest
	}
	return tags, nil
}

func ListContainerImagesWithoutTags(repository string) ([]*ContainerImage, error) {
	out, err := run("container", "images", "list-tags", repository, "--filter=\"NOT tags:*\"")
	if err != nil {
		return nil, err
	}
	var tags []*ContainerImage
	if err := json.Unmarshal(out, &tags); err != nil {
		return nil, err
	}
	for _, tag := range tags {
		tag.Repository = repository
		tag.FullyQualifiedDigest = repository + "@" + tag.Digest
	}
	return tags, nil
}

func DeleteContainerImage(image *ContainerImage) error {
	_, err := run("container", "images", "delete", image.FullyQualifiedDigest)
	if err != nil {
		return err
	}
	return nil
}

type ContainerImage struct {
	Repository           string   `json:"image"`
	Digest               string   `json:"digest"`
	FullyQualifiedDigest string   `json:"fully_qualified_digest"`
	Tags                 []string `json:"tags"`
	Timestamp            struct {
		Datetime    string `json:"datetime"`
		Day         int    `json:"day"`
		Hour        int    `json:"hour"`
		Microsecond int    `json:"microsecond"`
		Minute      int    `json:"minute"`
		Month       int    `json:"month"`
		Second      int    `json:"second"`
		Year        int    `json:"year"`
	} `json:"timestamp"`
}
