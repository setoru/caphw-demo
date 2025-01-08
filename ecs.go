package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	ecsMdl "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
)

const (
	InstanceDefaultTimeout = 900

	DefaultWaitForInterval = 5

	ECSInstanceStatusRunning = "ACTIVE"
)

func CreateEcs(client Client, config *Config, vpcId, networkId string, securityGroupIds []string) ([]*ecsMdl.ServerDetail, error) {
	request := &ecsMdl.CreateServersRequest{Body: &ecsMdl.CreateServersRequestBody{Server: &ecsMdl.PrePaidServer{}}}
	request.Body.Server.Name = config.ECSName

	imageId, minDisk, err := GetImageIdAndDisk(config, client)
	if err != nil {
		return nil, err
	}
	request.Body.Server.ImageRef = imageId

	request.Body.Server.FlavorRef = config.Flavor
	request.Body.Server.Vpcid = vpcId

	nics := make([]ecsMdl.PrePaidServerNic, 0)
	nics = append(nics, ecsMdl.PrePaidServerNic{SubnetId: networkId})
	request.Body.Server.Nics = nics

	securityGroups := make([]ecsMdl.PrePaidServerSecurityGroup, 0)
	for _, id := range securityGroupIds {
		securityGroups = append(securityGroups, ecsMdl.PrePaidServerSecurityGroup{Id: &id})
	}
	request.Body.Server.SecurityGroups = &securityGroups

	if config.PublicIp {
		publicIp := ecsMdl.PrePaidServerPublicip{}
		publicIp.Eip = &ecsMdl.PrePaidServerEip{
			Iptype: config.PublicIpSpec.IpType,
			Bandwidth: &ecsMdl.PrePaidServerEipBandwidth{
				Size:       &config.PublicIpSpec.Size,
				Sharetype:  getPublicIpShareType(config.PublicIpSpec.ShareType),
				Chargemode: &config.PublicIpSpec.ChargeMode,
			},
		}
		request.Body.Server.Publicip = &publicIp
	}

	request.Body.Server.ServerTags = buildTagList(config)

	if config.Userdata != "" {
		request.Body.Server.UserData = &config.Userdata
	}

	request.Body.Server.RootVolume = getRootVolumeProperties(minDisk, config)

	request.Body.Server.DataVolumes = getDataVolumeProperties(minDisk, config)

	request.Body.Server.AvailabilityZone = &config.AvailabilityZone

	request.Body.Server.BatchCreateInMultiAz = &config.BatchCreateInMultiAz

	//use iso image
	//metedata := make(map[string]string)
	//metedata["virtual_env_type"] = "IsoImage"
	//request.Body.Server.Metadata = metedata

	request.Body.Server.Extendparam = getCharging(config)

	diskPrior := "true"
	request.Body.Server.Extendparam.DiskPrior = &diskPrior
	request.Body.Server.Extendparam.RegionID = &config.RegionId

	request.Body.Server.OsschedulerHints = getServerSchedulerHints(config)

	response, err := client.CreateServers(request)
	if err != nil {
		return nil, err
	}
	// Sleep
	time.Sleep(5 * time.Second)

	instance, err := waitForInstancesStatus(client, config.PublicIp, *response.ServerIds, ECSInstanceStatusRunning, InstanceDefaultTimeout)
	if err != nil {
		return nil, err
	}
	fmt.Println("create ecs success")
	return instance, nil
}

func getServerSchedulerHints(config *Config) *ecsMdl.PrePaidServerSchedulerHints {
	if config.ServerSchedulerHints.Group == "" {
		return nil
	}
	var serverSchedulerHints ecsMdl.PrePaidServerSchedulerHints
	switch config.ServerSchedulerHints.Tenancy {
	case "shared":
		shared := ecsMdl.GetPrePaidServerSchedulerHintsTenancyEnum().SHARED
		serverSchedulerHints.Tenancy = &shared
	case "dedicated":
		dedicated := ecsMdl.GetPrePaidServerSchedulerHintsTenancyEnum().DEDICATED
		serverSchedulerHints.Tenancy = &dedicated
	}
	serverSchedulerHints.Group = &config.ServerSchedulerHints.Group
	serverSchedulerHints.DedicatedHostId = &config.ServerSchedulerHints.DedicatedHostId
	return &serverSchedulerHints
}

func getCharging(config *Config) *ecsMdl.PrePaidServerExtendParam {
	var charging ecsMdl.PrePaidServerExtendParam
	var chargingMode ecsMdl.PrePaidServerExtendParamChargingMode
	var isAutoPay ecsMdl.PrePaidServerExtendParamIsAutoPay
	var isAutoRenew ecsMdl.PrePaidServerExtendParamIsAutoRenew

	switch config.Charging.ChargingMode {
	case "prePaid":
		chargingMode = ecsMdl.GetPrePaidServerExtendParamChargingModeEnum().PRE_PAID
		periodType := config.Charging.PeriodType
		periodNum := config.Charging.PeriodNum
		switch periodType {
		case "month":
			month := ecsMdl.GetPrePaidServerExtendParamPeriodTypeEnum().MONTH
			charging.PeriodType = &month
		case "year":
			year := ecsMdl.GetPrePaidServerExtendParamPeriodTypeEnum().YEAR
			charging.PeriodType = &year
		}
		charging.PeriodNum = &periodNum
		switch config.Charging.IsAutoPay {
		case true:
			isAutoPay = ecsMdl.GetPrePaidServerExtendParamIsAutoPayEnum().TRUE
		case false:
			isAutoPay = ecsMdl.GetPrePaidServerExtendParamIsAutoPayEnum().FALSE
		}

		switch config.Charging.IsAutoRenew {
		case true:
			isAutoRenew = ecsMdl.GetPrePaidServerExtendParamIsAutoRenewEnum().TRUE
		case false:
			isAutoRenew = ecsMdl.GetPrePaidServerExtendParamIsAutoRenewEnum().FALSE
		}
	case "PostPaid":
		chargingMode = ecsMdl.GetPrePaidServerExtendParamChargingModeEnum().POST_PAID
	default:
		chargingMode = ecsMdl.GetPrePaidServerExtendParamChargingModeEnum().POST_PAID
	}

	charging.ChargingMode = &chargingMode
	charging.IsAutoPay = &isAutoPay
	charging.IsAutoRenew = &isAutoRenew
	charging.EnterpriseProjectId = config.Charging.EnterpriseProjectId
	return &charging
}

func getRootVolumeProperties(minDisk int32, config *Config) *ecsMdl.PrePaidServerRootVolume {
	rootVolume := ecsMdl.PrePaidServerRootVolume{}
	switch config.RootVolume.VolumeType {
	case "SSD":
		rootVolume.Volumetype = ecsMdl.GetPrePaidServerRootVolumeVolumetypeEnum().SSD
	case "GPSSD":
		rootVolume.Volumetype = ecsMdl.GetPrePaidServerRootVolumeVolumetypeEnum().GPSSD
	case "SATA":
		rootVolume.Volumetype = ecsMdl.GetPrePaidServerRootVolumeVolumetypeEnum().SATA
	case "SAS":
		rootVolume.Volumetype = ecsMdl.GetPrePaidServerRootVolumeVolumetypeEnum().SAS
	case "GPSSD2":
		rootVolume.Volumetype = ecsMdl.GetPrePaidServerRootVolumeVolumetypeEnum().GPSSD2
		rootVolume.Iops = &config.RootVolume.Iops
		rootVolume.Throughput = &config.RootVolume.Throughput
	case "ESSD":
		rootVolume.Volumetype = ecsMdl.GetPrePaidServerRootVolumeVolumetypeEnum().ESSD
	case "ESSD2":
		rootVolume.Volumetype = ecsMdl.GetPrePaidServerRootVolumeVolumetypeEnum().ESSD2
		rootVolume.Iops = &config.RootVolume.Iops
	default:
		rootVolume.Volumetype = ecsMdl.GetPrePaidServerRootVolumeVolumetypeEnum().GPSSD
	}
	if config.RootVolume.Size > minDisk {
		rootVolume.Size = &config.RootVolume.Size
	}
	rootVolume.Extendparam = &ecsMdl.PrePaidServerRootVolumeExtendParam{SnapshotId: &config.RootVolume.SnapshotId}
	return &rootVolume
}

func getDataVolumeProperties(minDisk int32, config *Config) *[]ecsMdl.PrePaidServerDataVolume {
	dataVolumes := make([]ecsMdl.PrePaidServerDataVolume, 0)
	if len(config.DataVolumes) < 0 {
		dataVolumes = append(dataVolumes, ecsMdl.PrePaidServerDataVolume{Volumetype: ecsMdl.GetPrePaidServerDataVolumeVolumetypeEnum().GPSSD, Size: 40})
		return &dataVolumes
	}
	for _, volume := range config.DataVolumes {
		dataVolume := ecsMdl.PrePaidServerDataVolume{}
		switch volume.VolumeType {
		case "SSD":
			dataVolume.Volumetype = ecsMdl.GetPrePaidServerDataVolumeVolumetypeEnum().SSD
		case "GPSSD":
			dataVolume.Volumetype = ecsMdl.GetPrePaidServerDataVolumeVolumetypeEnum().GPSSD
		case "SATA":
			dataVolume.Volumetype = ecsMdl.GetPrePaidServerDataVolumeVolumetypeEnum().SATA
		case "SAS":
			dataVolume.Volumetype = ecsMdl.GetPrePaidServerDataVolumeVolumetypeEnum().SAS
		case "GPSSD2":
			dataVolume.Volumetype = ecsMdl.GetPrePaidServerDataVolumeVolumetypeEnum().GPSSD2
			dataVolume.Iops = &volume.Iops
			dataVolume.Throughput = &volume.Throughput
		case "ESSD":
			dataVolume.Volumetype = ecsMdl.GetPrePaidServerDataVolumeVolumetypeEnum().ESSD
		case "ESSD2":
			dataVolume.Volumetype = ecsMdl.GetPrePaidServerDataVolumeVolumetypeEnum().ESSD2
			dataVolume.Iops = &volume.Iops
		default:
			dataVolume.Volumetype = ecsMdl.GetPrePaidServerDataVolumeVolumetypeEnum().GPSSD
		}
		if volume.Size > minDisk {
			dataVolume.Size = volume.Size
		} else {
			dataVolume.Size = minDisk
		}
		dataVolume.Extendparam = &ecsMdl.PrePaidServerDataVolumeExtendParam{SnapshotId: &volume.SnapshotId}
		dataVolume.DataImageId = &volume.DataImageId
		if volume.ClusterId != "" {
			dataVolume.ClusterId = &volume.ClusterId
			clusterType := ecsMdl.GetPrePaidServerDataVolumeClusterTypeEnum().DSS
			dataVolume.ClusterType = &clusterType
		}
		dataVolume.Multiattach = &volume.Multiattach
		dataVolume.Hwpassthrough = &volume.Passthrough
		dataVolumes = append(dataVolumes, dataVolume)
	}
	return &dataVolumes
}

func buildTagList(config *Config) *[]ecsMdl.PrePaidServerTag {
	rawTagList := make([]ecsMdl.PrePaidServerTag, 0)
	for _, tag := range config.Tags {
		rawTagList = append(rawTagList, ecsMdl.PrePaidServerTag{Key: tag.Name, Value: tag.Value})
	}
	return removeDuplicatedTags(rawTagList)
}

func removeDuplicatedTags(tags []ecsMdl.PrePaidServerTag) *[]ecsMdl.PrePaidServerTag {
	m := make(map[string]bool)
	result := make([]ecsMdl.PrePaidServerTag, 0)

	for _, entry := range tags {
		if _, value := m[entry.Key]; !value {
			m[entry.Key] = true
			result = append(result, entry)
		}
	}
	return &result
}

func waitForInstancesStatus(client Client, publicIP bool, instanceIds []string, instanceStatus string, timeout int) ([]*ecsMdl.ServerDetail, error) {
	if timeout <= 0 {
		timeout = InstanceDefaultTimeout
	}
	result, err := WaitForResult(fmt.Sprintf("Wait for the instances %v state to change to %s ", instanceIds, instanceStatus), func() (stop bool, result interface{}, err error) {
		instances, err := describeInstances(client, instanceIds)
		if err != nil {
			return false, nil, err
		}

		if len(instances) <= 0 {
			return true, nil, fmt.Errorf("the instances %v not found. ", instanceIds)
		}

		idsLen := len(instanceIds)
		needInstances := make([]*ecsMdl.ServerDetail, 0)

		for _, instance := range instances {
			if instance.Status == instanceStatus {
				needInstances = append(needInstances, &instance)
			}
			if publicIP {
				wait := true
				for _, addresses := range instance.Addresses {
					for _, address := range addresses {
						if *address.OSEXTIPStype == ecsMdl.GetServerAddressOSEXTIPStypeEnum().FLOATING {
							wait = false
						}
					}
				}
				if wait {
					return false, nil, fmt.Errorf("wait for public ip ")
				}
			}
		}

		if len(needInstances) == idsLen {
			return true, needInstances, nil
		}

		return false, nil, fmt.Errorf("the instances  %v state are not  the expected state  %s ", instanceIds, instanceStatus)
	}, false, DefaultWaitForInterval, timeout)

	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.([]*ecsMdl.ServerDetail), nil
}

func WaitForResult(name string, predicate func() (bool, interface{}, error), returnWhenError bool, delay int, timeout int) (interface{}, error) {
	endTime := time.Now().Add(time.Duration(timeout) * time.Second)
	delaySecond := time.Duration(delay) * time.Second
	for {
		satisfied, result, err := predicate()
		if err != nil {
			log.Print(fmt.Sprintf("%s Invoke func: %++s error: %++v", name, "predicate func() (bool, error)", err))
			if returnWhenError {
				return result, err
			}
		}
		if satisfied {
			return result, nil
		}
		time.Sleep(delaySecond)

		if timeout >= 0 && time.Now().After(endTime) {
			return nil, fmt.Errorf("wait for %s timeout", name)
		}
	}
}

func StopEcs(client Client, instanceIds []string) error {
	request := ecsMdl.BatchStopServersRequest{Body: &ecsMdl.BatchStopServersRequestBody{OsStop: &ecsMdl.BatchStopServersOption{}}}
	ids := make([]ecsMdl.ServerId, 0)
	for _, instanceId := range instanceIds {
		ids = append(ids, ecsMdl.ServerId{Id: instanceId})
	}
	request.Body.OsStop.Servers = ids
	_, err := client.BatchStopServers(&request)
	if err != nil {
		return err
	}
	fmt.Println("stop ecs success")
	return nil
}

func DeleteEcs(client Client, instanceIds []string) error {
	request := ecsMdl.DeleteServersRequest{}
	ids := make([]ecsMdl.ServerId, 0)
	for _, instanceId := range instanceIds {
		ids = append(ids, ecsMdl.ServerId{Id: instanceId})
	}
	deletePublicIp := true
	deleteVolume := true
	request.Body = &ecsMdl.DeleteServersRequestBody{
		Servers:        ids,
		DeletePublicip: &deletePublicIp,
		DeleteVolume:   &deleteVolume,
	}
	_, err := client.DeleteServers(&request)
	if err != nil {
		return err
	}
	fmt.Println("stop ecs success")
	return nil
}

func GetInstanceById(client Client, instanceId string) (*ecsMdl.ServerDetail, error) {
	if instanceId == "" {
		return nil, fmt.Errorf("instanceId not specified")
	}
	request := &ecsMdl.ShowServerRequest{}
	request.ServerId = instanceId
	response, err := client.ShowServer(request)
	if err != nil {
		return nil, err
	}
	return response.Server, nil
}

func GetExistingInstances(client Client, huaweiTag HuaweiTag) ([]ecsMdl.ServerDetail, error) {
	request := ecsMdl.ListServersDetailsRequest{}
	tag := fmt.Sprintf("%s,%s", huaweiTag.Name, huaweiTag.Value)
	request.Tags = &tag
	response, err := client.ListServersDetails(&request)
	if err != nil {
		return nil, err
	}
	return *response.Servers, nil
}

func describeInstances(client Client, instanceIds []string) ([]ecsMdl.ServerDetail, error) {
	if len(instanceIds) < 1 {
		return nil, fmt.Errorf("instance-ids not specified")
	}
	request := ecsMdl.ListServersDetailsRequest{}
	serverIds := strings.Join(instanceIds, ",")
	request.ServerId = &serverIds
	response, err := client.ListServersDetails(&request)
	if err != nil {
		return nil, err
	}
	return *response.Servers, nil
}

func getPublicIpShareType(shareType string) ecsMdl.PrePaidServerEipBandwidthSharetype {
	var shareTypeEnum ecsMdl.PrePaidServerEipBandwidthSharetype
	switch shareType {
	case "per":
		shareTypeEnum = ecsMdl.GetPrePaidServerEipBandwidthSharetypeEnum().PER
	case "whole":
		shareTypeEnum = ecsMdl.GetPrePaidServerEipBandwidthSharetypeEnum().WHOLE
	default:
		shareTypeEnum = ecsMdl.GetPrePaidServerEipBandwidthSharetypeEnum().PER
	}
	return shareTypeEnum
}
