package main

import (
	"fmt"

	ecsMdl "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
	elbMdl "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3/model"
)

func AddServerToElb(client Client, subnetId string, pools []HuaweiElbPool, servers []*ecsMdl.ServerDetail) (err error) {
	for _, pool := range pools {
		createMembersRequest := &elbMdl.BatchCreateMembersRequest{Body: &elbMdl.BatchCreateMembersRequestBody{Members: make([]elbMdl.BatchCreateMembersOption, 0)}}
		for _, server := range servers {
			for _, addresses := range server.Addresses {
				for _, address := range addresses {
					if *address.OSEXTIPStype == ecsMdl.GetServerAddressOSEXTIPStypeEnum().FIXED {
						createMembersRequest.PoolId = pool.Id
						createMember := elbMdl.BatchCreateMembersOption{Address: address.Addr, ProtocolPort: pool.Port, SubnetCidrId: &subnetId}
						createMembersRequest.Body.Members = append(createMembersRequest.Body.Members, createMember)
					}
				}
			}
		}
		_, err := client.BatchCreateMembers(createMembersRequest)
		if err != nil {
			return err
		}
	}
	return
}

func DeleteServerFromElb(client Client, pools []HuaweiElbPool, servers []*ecsMdl.ServerDetail) (err error) {
	for _, pool := range pools {
		deleteMembersRequest := &elbMdl.BatchDeleteMembersRequest{Body: &elbMdl.BatchDeleteMembersRequestBody{Members: make([]elbMdl.BatchDeleteMembersOption, 0)}}
		for _, server := range servers {
			for _, addresses := range server.Addresses {
				for _, address := range addresses {
					if *address.OSEXTIPStype == ecsMdl.GetServerAddressOSEXTIPStypeEnum().FIXED {
						deleteMembersRequest.PoolId = pool.Id
						addr := address.Addr
						deleteMember := elbMdl.BatchDeleteMembersOption{Address: &addr, ProtocolPort: &pool.Port}
						deleteMembersRequest.Body.Members = append(deleteMembersRequest.Body.Members, deleteMember)
					}
				}
			}
		}
		_, err := client.BatchDeleteMembers(deleteMembersRequest)
		if err != nil {
			return err
		}
	}
	return
}

func CreateElb(client Client, subnetId, vpcId string, config *Config) (string, error) {
	request := &elbMdl.CreateLoadBalancerRequest{}
	nameBandwidth := fmt.Sprintf("eip-%s", GenerateRandomString(4))
	chargeMode := GetLoadBalancerChargeMode(config.LoadBalancer.ChargingMode)
	shareType := GetLoadBalancerShareType(config.LoadBalancer.ShareType)
	bandwidthPublicip := &elbMdl.CreateLoadBalancerBandwidthOption{
		Name:       &nameBandwidth,
		Size:       &config.LoadBalancer.Size,
		ChargeMode: &chargeMode,
		ShareType:  &shareType,
	}
	publicipLoadbalancer := &elbMdl.CreateLoadBalancerPublicIpOption{
		NetworkType: "5_bgp",
		Bandwidth:   bandwidthPublicip,
	}
	zones, err := GetAvailabilityZones(client)
	if err != nil {
		return "", err
	}
	name := fmt.Sprintf("elb-%s", GenerateRandomString(4))
	loadbalancerbody := &elbMdl.CreateLoadBalancerOption{
		Name:                 &name,
		VipSubnetCidrId:      &subnetId,
		VpcId:                &vpcId,
		AvailabilityZoneList: zones,
		Publicip:             publicipLoadbalancer,
	}
	request.Body = &elbMdl.CreateLoadBalancerRequestBody{
		Loadbalancer: loadbalancerbody,
	}
	response, err := client.CreateLoadBalancer(request)
	if err != nil {
		return "", err
	}
	fmt.Println("create elb success")
	return response.Loadbalancer.Id, nil
}

func GetAvailabilityZones(client Client) ([]string, error) {
	availabilityZones := make([]string, 0)
	request := &elbMdl.ListAvailabilityZonesRequest{}
	response, err := client.ListAvailabilityZones(request)
	if err != nil {
		return nil, err
	}
	for _, zones := range *response.AvailabilityZones {
		for _, zone := range zones {
			availabilityZones = append(availabilityZones, zone.Code)
			if len(availabilityZones) == 2 {
				break
			}
		}
	}
	return availabilityZones, nil
}

func CreateListener(client Client, LoadbalancerId string, port int32) (string, error) {
	request := &elbMdl.CreateListenerRequest{}
	name := fmt.Sprintf("caphw_tcp-%s", GenerateRandomString(4))
	listenerbody := &elbMdl.CreateListenerOption{
		LoadbalancerId: LoadbalancerId,
		Name:           &name,
		Protocol:       "TCP",
		ProtocolPort:   &port,
	}
	request.Body = &elbMdl.CreateListenerRequestBody{
		Listener: listenerbody,
	}
	response, err := client.CreateListener(request)
	if err != nil {
		return "", err
	}
	fmt.Println("create listener success")
	return response.Listener.Id, nil
}

func CreatePool(client Client, listenerId, vpcId string) (string, error) {
	request := &elbMdl.CreatePoolRequest{}
	name := fmt.Sprintf("server_group-%s", GenerateRandomString(4))
	typePool := "instance"
	poolbody := &elbMdl.CreatePoolOption{
		LbAlgorithm: "ROUND_ROBIN",
		ListenerId:  &listenerId,
		Name:        &name,
		Protocol:    "TCP",
		VpcId:       &vpcId,
		Type:        &typePool,
	}
	request.Body = &elbMdl.CreatePoolRequestBody{
		Pool: poolbody,
	}
	response, err := client.CreatePool(request)
	if err != nil {
		return "", err
	}
	fmt.Println("create pool success")
	return response.Pool.Id, nil
}

func GetLoadBalancerChargeMode(chargeMode string) elbMdl.CreateLoadBalancerBandwidthOptionChargeMode {
	var chargeModeEnum elbMdl.CreateLoadBalancerBandwidthOptionChargeMode
	switch chargeMode {
	case "traffic":
		chargeModeEnum = elbMdl.GetCreateLoadBalancerBandwidthOptionChargeModeEnum().TRAFFIC
	case "bandwidth":
		chargeModeEnum = elbMdl.GetCreateLoadBalancerBandwidthOptionChargeModeEnum().BANDWIDTH
	default:
		chargeModeEnum = elbMdl.GetCreateLoadBalancerBandwidthOptionChargeModeEnum().TRAFFIC
	}
	return chargeModeEnum
}

func GetLoadBalancerShareType(shareType string) elbMdl.CreateLoadBalancerBandwidthOptionShareType {
	var shareTypeEnum elbMdl.CreateLoadBalancerBandwidthOptionShareType
	switch shareType {
	case "per":
		shareTypeEnum = elbMdl.GetCreateLoadBalancerBandwidthOptionShareTypeEnum().PER
	case "whole":
		shareTypeEnum = elbMdl.GetCreateLoadBalancerBandwidthOptionShareTypeEnum().WHOLE
	default:
		shareTypeEnum = elbMdl.GetCreateLoadBalancerBandwidthOptionShareTypeEnum().PER
	}
	return shareTypeEnum
}

func DeleteElb(client Client, loadBalancerId string) error {
	req := &elbMdl.DeleteLoadBalancerRequest{LoadbalancerId: loadBalancerId}
	_, err := client.DeleteLoadBalancer(req)
	if err != nil {
		return err
	}
	return nil
}

func DeleteListener(client Client, listenerId string) error {
	req := &elbMdl.DeleteListenerRequest{ListenerId: listenerId}
	_, err := client.DeleteListener(req)
	if err != nil {
		return err
	}
	return nil
}

func DeletePool(client Client, poolId string) error {
	req := &elbMdl.DeletePoolRequest{PoolId: poolId}
	_, err := client.DeletePool(req)
	if err != nil {
		return err
	}
	return nil
}
