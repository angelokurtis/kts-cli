package gcloud

import (
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/angelokurtis/kts-cli/internal/color"
	"github.com/cheggaaa/pb/v3"
	"strings"
	"time"
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
	projects, err := ListProjects()
	if err != nil {
		return nil, err
	}

	repos := make([]string, 0)
	for _, project := range projects {
		r, err := listContainerRepositories(project)
		if err != nil {
			if !strings.Contains(err.Error(), "Bad status during token exchange: 403") {
				return nil, err
			}
			fmt.Printf(color.Error, "[WARN] You don't have permissions to list container images on project '"+project.Name+"'\n")
		}
		repos = append(repos, r...)
	}

	return repos, nil
}

func listContainerRepositories(project *Project) ([]string, error) {
	out, err := run("container", "images", "list", "--repository=gcr.io/"+project.ID)
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
		timestamp := &tag.Timestamp
		location, err := time.LoadLocation("America/Sao_Paulo")
		if err != nil {
			return nil, err
		}
		timestamp.Datetime = time.Date(timestamp.Year, time.Month(timestamp.Month), timestamp.Day, timestamp.Hour, timestamp.Minute, timestamp.Second, 0, location)
		tag.Repository = repository
		tag.FullyQualifiedDigest = repository + "@" + tag.Digest
	}
	return tags, nil
}

func ListContainerImagesWithoutTags(repository string) ([]*ContainerImage, error) {
	out, err := run("container", "images", "list-tags", repository)
	if err != nil {
		return nil, err
	}
	var tags []*ContainerImage
	if err := json.Unmarshal(out, &tags); err != nil {
		return nil, err
	}
	images := make([]*ContainerImage, 0, 0)
	for _, tag := range tags {
		if len(tag.Tags) == 0 {
			tag.Repository = repository
			tag.FullyQualifiedDigest = repository + "@" + tag.Digest
			images = append(images, tag)
		}
	}
	return images, nil
}

func DeleteContainerImage(image *ContainerImage) error {
	_, err := run("container", "images", "delete", image.FullyQualifiedDigest)
	if err != nil {
		return err
	}
	return nil
}

func SelectTags() ([]string, error) {
	fmt.Printf(color.Warning, "gcloud container images list\n")
	repositories, err := SelectContainerRepositories()
	if err != nil {
		return nil, err
	}

	tags := make([]string, 0, 0)
	if len(repositories) > 0 {
		fmt.Printf(color.Warning, "gcloud container images list-tags gcr.io/<PROJECT_ID>/<IMAGE_PATH>\n")
		bar := pb.StartNew(len(repositories))
		for _, repository := range repositories {
			images, err := ListContainerImages(repository)
			if err != nil {
				return nil, err
			}
			for _, image := range images {
				for _, tag := range image.Tags {
					tags = append(tags, image.Repository+":"+tag)
				}
			}
			bar.Increment()
		}
		bar.Finish()
	}

	var selects []string
	prompt := &survey.MultiSelect{
		Message: "Select the images tags:",
		Options: tags,
	}

	err = survey.AskOne(prompt, &selects, survey.WithPageSize(25))
	if err != nil {
		return nil, err
	}

	return selects, nil
}

func UntagImages(tags []string) error {
	bar := pb.StartNew(len(tags))
	for _, tag := range tags {
		err := UntagImage(tag)
		if err != nil {
			return err
		}
		bar.Increment()
	}
	bar.Finish()
	return nil
}

func UntagImage(tag string) error {
	_, err := run("container", "images", "untag", tag)
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
		Datetime    time.Time `json:"-"`
		Day         int       `json:"day"`
		Hour        int       `json:"hour"`
		Microsecond int       `json:"microsecond"`
		Minute      int       `json:"minute"`
		Month       int       `json:"month"`
		Second      int       `json:"second"`
		Year        int       `json:"year"`
	} `json:"timestamp"`
}
