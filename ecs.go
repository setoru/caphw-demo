package main

import (
	"fmt"
	"log"
	"time"

	ecsMdl "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
)

const (
	InstanceDefaultTimeout = 900

	DefaultWaitForInterval = 5

	ECSInstanceStatusRunning = "ACTIVE"
)

func CreateECS(client Client, config *Config, vpcId, networkId string, securityGroupIds []string) ([]*ecsMdl.ServerDetail, error) {
	request := &ecsMdl.CreateServersRequest{Body: &ecsMdl.CreateServersRequestBody{Server: &ecsMdl.PrePaidServer{}}}
	request.Body.Server.Name = config.ECSName

	imageID, minDisk, err := GetImageIDAndDisk(config, client)
	if err != nil {
		return nil, err
	}
	request.Body.Server.ImageRef = imageID

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

	if config.PublicIP {
		var size int32 = 100
		mode := "traffic"
		publicip := ecsMdl.PrePaidServerPublicip{Eip: &ecsMdl.PrePaidServerEip{}}
		publicip.Eip.Iptype = "5_bgp"
		publicip.Eip.Bandwidth = &ecsMdl.PrePaidServerEipBandwidth{Size: &size, Sharetype: ecsMdl.GetPrePaidServerEipBandwidthSharetypeEnum().PER, Chargemode: &mode}

		request.Body.Server.Publicip = &publicip
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
	request.Body.Server.Extendparam.RegionID = &config.RegionID

	request.Body.Server.OsschedulerHints = getServerSchedulerHints(config)

	response, err := client.CreateServers(request)
	if err != nil {
		return nil, err
	}
	// Sleep
	time.Sleep(5 * time.Second)

	instance, err := waitForInstancesStatus(client, config.PublicIP, *response.ServerIds, ECSInstanceStatusRunning, InstanceDefaultTimeout)
	if err != nil {
		return nil, err
	}
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
	rootVolume.Extendparam = &ecsMdl.PrePaidServerRootVolumeExtendParam{SnapshotId: &config.RootVolume.SnapshotID}
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
		dataVolume.Extendparam = &ecsMdl.PrePaidServerDataVolumeExtendParam{SnapshotId: &volume.SnapshotID}
		dataVolume.DataImageId = &volume.DataImageId
		if volume.ClusterID != "" {
			dataVolume.ClusterId = &volume.ClusterID
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

// Scan machine tags, and return a deduped tags list. The first found value gets precedence.
func removeDuplicatedTags(tags []ecsMdl.PrePaidServerTag) *[]ecsMdl.PrePaidServerTag {
	m := make(map[string]bool)
	result := make([]ecsMdl.PrePaidServerTag, 0)

	// look for duplicates
	for _, entry := range tags {
		if _, value := m[entry.Key]; !value {
			m[entry.Key] = true
			result = append(result, entry)
		}
	}
	return &result
}

// waitForInstancesStatus waits for instances to given status when instance.NotFound wait until timeout
func waitForInstancesStatus(client Client, publicIP bool, instanceIds []string, instanceStatus string, timeout int) ([]*ecsMdl.ServerDetail, error) {
	if timeout <= 0 {
		timeout = InstanceDefaultTimeout
	}
	result, err := WaitForResult(fmt.Sprintf("Wait for the instances %v state to change to %s ", instanceIds, instanceStatus), func() (stop bool, result interface{}, err error) {
		showServerRequest := ecsMdl.ListServersDetailsRequest{}
		var serverIds string
		for _, id := range instanceIds {
			if serverIds == "" {
				serverIds = id
				continue
			}
			serverIds += ","
			serverIds += id
		}
		showServerRequest.ServerId = &serverIds
		listServersDetailsResponse, err := client.ListServersDetails(&showServerRequest)
		if err != nil {
			return false, nil, err
		}

		if len(*listServersDetailsResponse.Servers) <= 0 {
			return true, nil, fmt.Errorf("the instances %v not found. ", instanceIds)
		}

		idsLen := len(instanceIds)
		servers := make([]*ecsMdl.ServerDetail, 0)

		for _, server := range *listServersDetailsResponse.Servers {
			if server.Status == instanceStatus {
				servers = append(servers, &server)
			}
			if publicIP {
				wait := true
				for _, addresses := range server.Addresses {
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

		if len(servers) == idsLen {
			return true, servers, nil
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

// WaitForResult wait func
func WaitForResult(name string, predicate func() (bool, interface{}, error), returnWhenError bool, delay int, timeout int) (interface{}, error) {
	endTime := time.Now().Add(time.Duration(timeout) * time.Second)
	delaySecond := time.Duration(delay) * time.Second
	for {
		// Execute the function
		satisfied, result, err := predicate()
		if err != nil {
			log.Print(fmt.Sprintf("%s\nInvoke func: %++s\n error: %++v", name, "predicate func() (bool, error)", err))
			if returnWhenError {
				return result, err
			}
		}
		if satisfied {
			return result, nil
		}
		// Sleep
		time.Sleep(delaySecond)
		// If a timeout is set, and that's been exceeded, shut it down
		if timeout >= 0 && time.Now().After(endTime) {
			return nil, fmt.Errorf("wait for %s timeout", name)
		}
	}
}
