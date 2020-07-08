package gcloud

import (
	"encoding/json"
	"time"
)

func ListAllClusters() ([]*Cluster, error) {
	projects, err := ListProjects()
	if err != nil {
		return nil, err
	}
	clusters := make([]*Cluster, 0, 0)
	for _, project := range projects {
		cts, err := ListClusters(project)
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

func ListClusters(project *Project) ([]*Cluster, error) {
	out, err := run("container", "clusters", "list", "--project", project.ID)
	if err != nil {
		return nil, err
	}

	var clusters []*Cluster
	if err := json.Unmarshal(out, &clusters); err != nil {
		return nil, err
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
