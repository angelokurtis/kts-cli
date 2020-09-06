package aws

import (
	"encoding/json"
	"github.com/AlecAivazis/survey/v2"
	"github.com/angelokurtis/kts-cli/pkg/bash"
	"github.com/pkg/errors"
)

func ListECRImages(repo *ECRRepository) (*ECRImages, error) {
	out, err := bash.RunAndLogRead("aws ecr describe-images --repository-name " + repo.Name)
	if err != nil {
		return nil, err
	}

	var images *ECRImages
	if err := json.Unmarshal(out, &images); err != nil {
		return nil, errors.WithStack(err)
	}

	return images, nil
}

func ListECRRepositories() (*ECRRepositories, error) {
	out, err := bash.RunAndLogRead("aws ecr describe-repositories")
	if err != nil {
		return nil, err
	}

	var repos *ECRRepositories
	if err := json.Unmarshal(out, &repos); err != nil {
		return nil, errors.WithStack(err)
	}

	return repos, nil
}

type ECRRepositories struct {
	Items []*ECRRepository `json:"repositories"`
}

func (e *ECRRepositories) SelectMany() (*ECRRepositories, error) {
	repos := make(map[string]*ECRRepository, 0)
	uris := make([]string, 0, 0)
	for _, repo := range e.Items {
		repos[repo.URI] = repo
		uris = append(uris, repo.URI)
	}

	var selects []string
	prompt := &survey.MultiSelect{
		Message: "Select the ECR repositories:",
		Options: uris,
	}

	err := survey.AskOne(prompt, &selects, survey.WithPageSize(10))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := make([]*ECRRepository, 0, 0)
	for _, s := range selects {
		result = append(result, repos[s])
	}

	return &ECRRepositories{Items: result}, nil
}

type ECRRepository struct {
	Arn                        string `json:"repositoryArn"`
	RegistryID                 string `json:"registryId"`
	Name                       string `json:"repositoryName"`
	URI                        string `json:"repositoryUri"`
	CreatedAt                  string `json:"createdAt"`
	ImageTagMutability         string `json:"imageTagMutability"`
	ImageScanningConfiguration struct {
		ScanOnPush bool `json:"scanOnPush"`
	} `json:"imageScanningConfiguration"`
	EncryptionConfiguration struct {
		EncryptionType string `json:"encryptionType"`
	} `json:"encryptionConfiguration"`
}

type ECRImages struct {
	Items []*ECRImage `json:"imageDetails"`
}

type ECRImage struct {
	RegistryID     string   `json:"registryId"`
	RepositoryName string   `json:"repositoryName"`
	Digest         string   `json:"imageDigest"`
	Tags           []string `json:"imageTags"`
	SizeInBytes    int      `json:"imageSizeInBytes"`
	PushedAt       string   `json:"imagePushedAt"`
}
