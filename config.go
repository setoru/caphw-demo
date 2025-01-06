package main

import (
	"flag"
	"os"

	"k8s.io/apimachinery/pkg/util/yaml"
)

type Config struct {
	ECSName              string                 `json:"ecsName,omitempty"`
	Flavor               string                 `json:"flavor,omitempty"`
	Vpc                  HuaweiVpc              `json:"vpc,omitempty"`
	Subnet               Subnet                 `json:"subnet,omitempty"`
	RegionID             string                 `json:"regionId,omitempty"`
	ImageID              string                 `json:"imageId,omitempty"`
	SecurityGroups       []SecurityGroup        `json:"securityGroups,omitempty"`
	PublicIP             bool                   `json:"publicIp,omitempty"`
	Tags                 []HuaweiTag            `json:"tags,omitempty"`
	RootVolume           RootVolumeProperties   `json:"rootVolume,omitempty"`
	DataVolumes          []DataVolumeProperties `json:"dataVolumes,omitempty"`
	ElbMembers           []HuaweiElbMembers     `json:"elbMembers,omitempty"`
	Charging             Charging               `json:"charging,omitempty"`
	AvailabilityZone     string                 `json:"availabilityZone,omitempty"`
	BatchCreateInMultiAz bool                   `json:"batchCreateInMultiAz,omitempty"`
	ServerSchedulerHints ServerSchedulerHints   `json:"serverSchedulerHints,omitempty"`
	Userdata             string                 `json:"userdata,omitempty"`
}

type RootVolumeProperties struct {
	VolumeType string `json:"volumeType,omitempty"`
	Size       int32  `json:"size,omitempty"`
	Iops       int32  `json:"iops,omitempty"`
	Throughput int32  `json:"throughput,omitempty"`
	SnapshotID string `json:"snapshotId,omitempty"`
}

type Charging struct {
	ChargingMode        string  `json:"chargingMode,omitempty"`
	PeriodType          string  `json:"periodType,omitempty"`
	PeriodNum           int32   `json:"periodNum,omitempty"`
	IsAutoPay           bool    `json:"isAutoPay,omitempty"`
	IsAutoRenew         bool    `json:"isAutoRenew,omitempty"`
	EnterpriseProjectId *string `json:"enterpriseProjectId,omitempty"`
}

type ServerSchedulerHints struct {
	Group           string `json:"group,omitempty"`
	Tenancy         string `json:"tenancy,omitempty"`
	DedicatedHostId string `json:"dedicatedHostId,omitempty"`
}

type DataVolumeProperties struct {
	VolumeType  string `json:"volumeType,omitempty"`
	Size        int32  `json:"size,omitempty"`
	Iops        int32  `json:"iops,omitempty"`
	Throughput  int32  `json:"throughput,omitempty"`
	SnapshotID  string `json:"snapshotId,omitempty"`
	Multiattach bool   `json:"multiattach,omitempty"`
	Passthrough bool   `json:"passthrough,omitempty"`
	ClusterID   string `json:"clusterId,omitempty"`
	ClusterType string `json:"clusterType,omitempty"`
	DataImageId string `json:"dataImageId,omitempty"`
}

type SecurityGroup struct {
	ID                  string  `json:"id,omitempty"`
	Name                string  `json:"name,omitempty"`
	EnterpriseProjectId *string `json:"enterpriseProjectId,omitempty"`
}

type Subnet struct {
	NetworkId string `json:"networkId,omitempty"`
	Name      string `json:"name,omitempty"`
	SubnetId  string `json:"subnetId,omitempty"`
	Cidr      string `json:"cidr,omitempty"`
}

type HuaweiTag struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type HuaweiElbMembers struct {
	ID   string `json:"id,omitempty"`
	Port int32  `json:"port,omitempty"`
}

type HuaweiVpc struct {
	ID                  string  `json:"id,omitempty"`
	Name                string  `json:"name,omitempty"`
	Cidr                string  `json:"cidr,omitempty"`
	EnterpriseProjectId *string `json:"enterpriseProjectId,omitempty"`
}

var configFile = flag.String("f", "config.yaml", "the config file")

func parseConfig() (*Config, error) {
	file, err := os.ReadFile(*configFile)
	if err != nil {
		return nil, err
	}

	config := Config{}
	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
