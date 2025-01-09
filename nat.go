package main

import (
	"fmt"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/nat/v2/model"
)

func CreateNat(client Client, config *Config, vpcId, networkId string) error {
	request := &model.CreateNatGatewayRequest{}
	request.Body = &model.CreateNatGatewayRequestBody{
		NatGateway: &model.CreateNatGatewayOption{
			Name:              fmt.Sprintf("nat-%s", GenerateRandomString(4)),
			RouterId:          vpcId,
			InternalNetworkId: networkId,
			Spec:              getNatSpec(config.Nat.Spec),
		},
	}
	response, err := client.CreateNatGateway(request)
	if err != nil {
		return err
	}
	publicIpId, err := CreatePublicIp(client, config)
	if err != nil {
		return err
	}
	snatRequest := &model.CreateNatGatewaySnatRuleRequest{}
	snatRequest.Body = &model.CreateNatGatewaySnatRuleRequestOption{
		SnatRule: &model.CreateNatGatewaySnatRuleOption{
			NatGatewayId: response.NatGateway.Id,
			FloatingIpId: publicIpId,
			NetworkId:    &networkId,
		},
	}
	_, err = client.CreateNatGatewaySnatRule(snatRequest)
	if err != nil {
		return err
	}
	return nil
}

func getNatSpec(spec int32) model.CreateNatGatewayOptionSpec {
	var specEnum model.CreateNatGatewayOptionSpec
	switch spec {
	case 1:
		specEnum = model.GetCreateNatGatewayOptionSpecEnum().E_1
	case 2:
		specEnum = model.GetCreateNatGatewayOptionSpecEnum().E_2
	case 3:
		specEnum = model.GetCreateNatGatewayOptionSpecEnum().E_3
	case 4:
		specEnum = model.GetCreateNatGatewayOptionSpecEnum().E_4
	default:
		specEnum = model.GetCreateNatGatewayOptionSpecEnum().E_1
	}
	return specEnum
}
