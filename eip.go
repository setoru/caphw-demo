package main

import (
	"fmt"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/eip/v2/model"
)

func CreatePublicIp(client Client, config *Config) (id string, err error) {
	request := &model.CreatePublicipRequest{}
	publicIpBody := &model.CreatePublicipOption{
		Type: getIpType(config.Nat.PublicIpSpec.IpType),
	}
	chargeMode := getPublicIpChargeMode(config.Nat.PublicIpSpec.ChargeMode)
	name := fmt.Sprintf("eip-%s", GenerateRandomString(4))
	bandwidthBody := &model.CreatePublicipBandwidthOption{
		ChargeMode: &chargeMode,
		Name:       &name,
		ShareType:  getNatPublicIpShareType(config.Nat.PublicIpSpec.ShareType),
		Size:       &config.Nat.PublicIpSpec.Size,
	}
	request.Body = &model.CreatePublicipRequestBody{
		Publicip:  publicIpBody,
		Bandwidth: bandwidthBody,
	}
	response, err := client.CreatePublicip(request)
	if err != nil {
		return "", err
	}
	return *response.Publicip.Id, nil
}

func getPublicIpChargeMode(chargeMode string) model.CreatePublicipBandwidthOptionChargeMode {
	var chargeModeEnum model.CreatePublicipBandwidthOptionChargeMode
	switch chargeMode {
	case "traffic":
		chargeModeEnum = model.GetCreatePublicipBandwidthOptionChargeModeEnum().TRAFFIC
	case "bandwidth":
		chargeModeEnum = model.GetCreatePublicipBandwidthOptionChargeModeEnum().BANDWIDTH
	default:
		chargeModeEnum = model.GetCreatePublicipBandwidthOptionChargeModeEnum().TRAFFIC
	}
	return chargeModeEnum
}

func getNatPublicIpShareType(shareType string) model.CreatePublicipBandwidthOptionShareType {
	var shareTypeEnum model.CreatePublicipBandwidthOptionShareType
	switch shareType {
	case "per":
		shareTypeEnum = model.GetCreatePublicipBandwidthOptionShareTypeEnum().PER
	case "whole":
		shareTypeEnum = model.GetCreatePublicipBandwidthOptionShareTypeEnum().WHOLE
	default:
		shareTypeEnum = model.GetCreatePublicipBandwidthOptionShareTypeEnum().PER
	}
	return shareTypeEnum
}
