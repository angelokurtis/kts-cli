package aws

import (
	"encoding/json"
	"fmt"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"

	"github.com/angelokurtis/kts-cli/pkg/bash"
)

func ConnectToEKSCluster(cluster string) error {
	_, err := bash.RunAndLogWrite(fmt.Sprintf("aws eks update-kubeconfig --name %s", cluster))
	if err != nil {
		return err
	}

	return nil
}

func SelectEKSCluster() (string, error) {
	clusters, err := ListEKSClusters()
	if err != nil {
		return "", err
	}

	var selected string

	if len(clusters) == 0 {
		return "", nil
	} else if len(clusters) > 1 {
		prompt := &survey.Select{
			Message: "Select the EKS Cluster:",
			Options: clusters,
		}

		err = survey.AskOne(prompt, &selected, survey.WithPageSize(10), survey.WithKeepFilter(true))
		if err != nil {
			return "", errors.WithStack(err)
		}
	} else {
		selected = clusters[0]
	}

	return selected, nil
}

func ListEKSClusters(region ...string) ([]string, error) {
	if len(region) == 0 {
		out, err := bash.RunAndLogRead("aws eks list-clusters")
		if err != nil {
			return nil, err
		}

		var eks *eksClusterList
		if err := json.Unmarshal(out, &eks); err != nil {
			return nil, errors.WithStack(err)
		}

		return eks.Clusters, nil
	}

	clusters := make([]string, 0, 0)

	for _, r := range region {
		out, err := bash.RunAndLogRead("aws eks list-clusters --region " + r)
		if err != nil {
			return nil, err
		}

		var eks *eksClusterList
		if err := json.Unmarshal(out, &eks); err != nil {
			return nil, errors.WithStack(err)
		}

		clusters = append(clusters, eks.Clusters...)
	}

	return clusters, nil
}

func DescribeEKSCluster(cluster string) (*Cluster, error) {
	out, err := bash.RunAndLogRead("aws eks describe-cluster --name " + cluster)
	if err != nil {
		return nil, err
	}

	var eks *eksClusterDetails
	if err := json.Unmarshal(out, &eks); err != nil {
		return nil, errors.WithStack(err)
	}

	return eks.Cluster, nil
}

type (
	eksClusterList struct {
		Clusters []string `json:"clusters"`
	}
	eksClusterDetails struct {
		Cluster *Cluster `json:"cluster"`
	}
)

type Cluster struct {
	Name      string `json:"name"`
	Arn       string `json:"arn"`
	CreatedAt string `json:"createdAt"`
	Version   string `json:"version"`
	Endpoint  string `json:"endpoint"`
	RoleArn   string `json:"roleArn"`
	Vpc       struct {
		SubnetIds        []string `json:"subnetIds"`
		SecurityGroupIds []string `json:"securityGroupIds"`
		ID               string   `json:"vpcId"`
	} `json:"resourcesVpcConfig"`
	Status               string `json:"status"`
	CertificateAuthority struct {
		Data string `json:"data"`
	} `json:"certificateAuthority"`
}
