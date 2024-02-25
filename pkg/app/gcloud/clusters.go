package gcloud

import (
	"encoding/json"
	"time"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"
)

func ListGKEClustersNames() ([]string, error) {
	clusters, err := ListGKEClusters()
	if err != nil {
		return nil, err
	}

	names := make([]string, 0, len(clusters))
	for _, cluster := range clusters {
		names = append(names, cluster.Name)
	}

	return names, nil
}

func ListGKEClusters() ([]*Cluster, error) {
	projects, err := ListProjects()
	if err != nil {
		return nil, err
	}

	clusters := make([]*Cluster, 0, 0)

	for _, project := range projects {
		cts, err := ListGKEClustersByProject(project)
		if err != nil {
			return nil, err
		}

		for _, ct := range cts {
			ct.Project = project
			clusters = append(clusters, ct)
		}
	}

	return clusters, nil
}

func SelectGKECluster() (*Cluster, error) {
	clusters, err := ListGKEClusters()
	if err != nil {
		return nil, err
	}

	if clusters == nil || len(clusters) == 0 {
		return nil, nil
	} else if len(clusters) == 1 {
		return clusters[0], nil
	}

	options := make([]string, 0, 0)
	m := make(map[string]*Cluster)

	for _, cluster := range clusters {
		options = append(options, cluster.Name)
		m[cluster.Name] = cluster
	}

	var k string

	prompt := &survey.Select{
		Message: "Select the Google Cluster:",
		Options: options,
	}

	err = survey.AskOne(prompt, &k, survey.WithKeepFilter(true))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return m[k], nil
}

func ListGKEClustersByProject(project *Project) ([]*Cluster, error) {
	out, err := runAndLogRead("container", "clusters", "list", "--project", project.ID)
	if err != nil {
		return nil, err
	}

	var clusters []*Cluster
	if err := json.Unmarshal(out, &clusters); err != nil {
		return nil, errors.WithStack(err)
	}

	return clusters, nil
}

type Cluster struct {
	Project                  *Project
	AddonsConfig             AddonsConfig             `json:"addonsConfig"`
	ClusterIpv4Cidr          string                   `json:"clusterIpv4Cidr"`
	CreateTime               time.Time                `json:"createTime"`
	CurrentMasterVersion     string                   `json:"currentMasterVersion"`
	CurrentNodeCount         int                      `json:"currentNodeCount"`
	CurrentNodeVersion       string                   `json:"currentNodeVersion"`
	DatabaseEncryption       DatabaseEncryption       `json:"databaseEncryption"`
	DefaultMaxPodsConstraint DefaultMaxPodsConstraint `json:"defaultMaxPodsConstraint"`
	Endpoint                 string                   `json:"endpoint"`
	InitialClusterVersion    string                   `json:"initialClusterVersion"`
	InstanceGroupUrls        []string                 `json:"instanceGroupUrls"`
	IPAllocationPolicy       IPAllocationPolicy       `json:"ipAllocationPolicy"`
	LabelFingerprint         string                   `json:"labelFingerprint"`
	Location                 string                   `json:"location"`
	Locations                []string                 `json:"locations"`
	LoggingService           string                   `json:"loggingService"`
	MaintenancePolicy        MaintenancePolicy        `json:"maintenancePolicy"`
	MasterAuth               MasterAuth               `json:"masterAuth"`
	MonitoringService        string                   `json:"monitoringService"`
	Name                     string                   `json:"name"`
	Network                  string                   `json:"network"`
	NetworkConfig            NetworkConfig            `json:"networkConfig"`
	NodeConfig               NodeConfig               `json:"nodeConfig"`
	NodePools                []NodePools              `json:"nodePools"`
	SelfLink                 string                   `json:"selfLink"`
	ServicesIpv4Cidr         string                   `json:"servicesIpv4Cidr"`
	Status                   string                   `json:"status"`
	Subnetwork               string                   `json:"subnetwork"`
	Zone                     string                   `json:"zone"`
}

type KubernetesDashboard struct {
	Disabled bool `json:"disabled"`
}

type NetworkPolicyConfig struct {
	Disabled bool `json:"disabled"`
}

type AddonsConfig struct {
	KubernetesDashboard KubernetesDashboard `json:"kubernetesDashboard"`
	NetworkPolicyConfig NetworkPolicyConfig `json:"networkPolicyConfig"`
}

type DatabaseEncryption struct {
	State string `json:"state"`
}

type DefaultMaxPodsConstraint struct {
	MaxPodsPerNode string `json:"maxPodsPerNode"`
}

type IPAllocationPolicy struct {
	ClusterIpv4Cidr            string `json:"clusterIpv4Cidr"`
	ClusterIpv4CidrBlock       string `json:"clusterIpv4CidrBlock"`
	ClusterSecondaryRangeName  string `json:"clusterSecondaryRangeName"`
	ServicesIpv4Cidr           string `json:"servicesIpv4Cidr"`
	ServicesIpv4CidrBlock      string `json:"servicesIpv4CidrBlock"`
	ServicesSecondaryRangeName string `json:"servicesSecondaryRangeName"`
	UseIPAliases               bool   `json:"useIpAliases"`
}

type MaintenancePolicy struct {
	ResourceVersion string `json:"resourceVersion"`
}

type MasterAuth struct {
	ClusterCaCertificate string `json:"clusterCaCertificate"`
}

type NetworkConfig struct {
	Network    string `json:"network"`
	Subnetwork string `json:"subnetwork"`
}

type Metadata struct {
	DisableLegacyEndpoints string `json:"disable-legacy-endpoints"`
}

type ShieldedInstanceConfig struct {
	EnableIntegrityMonitoring bool `json:"enableIntegrityMonitoring"`
}

type NodeConfig struct {
	DiskSizeGb             int                    `json:"diskSizeGb"`
	DiskType               string                 `json:"diskType"`
	ImageType              string                 `json:"imageType"`
	MachineType            string                 `json:"machineType"`
	Metadata               Metadata               `json:"metadata"`
	OauthScopes            []string               `json:"oauthScopes"`
	ServiceAccount         string                 `json:"serviceAccount"`
	ShieldedInstanceConfig ShieldedInstanceConfig `json:"shieldedInstanceConfig"`
}

type Config struct {
	DiskSizeGb             int                    `json:"diskSizeGb"`
	DiskType               string                 `json:"diskType"`
	ImageType              string                 `json:"imageType"`
	MachineType            string                 `json:"machineType"`
	Metadata               Metadata               `json:"metadata"`
	OauthScopes            []string               `json:"oauthScopes"`
	ServiceAccount         string                 `json:"serviceAccount"`
	ShieldedInstanceConfig ShieldedInstanceConfig `json:"shieldedInstanceConfig"`
}

type Management struct {
	AutoRepair  bool `json:"autoRepair"`
	AutoUpgrade bool `json:"autoUpgrade"`
}

type MaxPodsConstraint struct {
	MaxPodsPerNode string `json:"maxPodsPerNode"`
}

type UpgradeSettings struct {
	MaxSurge int `json:"maxSurge"`
}

type NodePools struct {
	Config            Config            `json:"config"`
	InitialNodeCount  int               `json:"initialNodeCount"`
	InstanceGroupUrls []string          `json:"instanceGroupUrls"`
	Locations         []string          `json:"locations"`
	Management        Management        `json:"management"`
	MaxPodsConstraint MaxPodsConstraint `json:"maxPodsConstraint"`
	Name              string            `json:"name"`
	PodIpv4CidrSize   int               `json:"podIpv4CidrSize"`
	SelfLink          string            `json:"selfLink"`
	Status            string            `json:"status"`
	UpgradeSettings   UpgradeSettings   `json:"upgradeSettings"`
	Version           string            `json:"version"`
}
