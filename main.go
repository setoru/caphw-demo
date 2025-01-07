package main

import (
	"log"
)

func main() {
	client, err := GetHuaweiClient()
	if err != nil {
		log.Print(err)
		return
	}

	config, err := parseConfig()
	if err != nil {
		log.Print(err)
		return
	}

	//create network
	vpcId, err := CreateVpc(client, config)
	if err != nil {
		log.Print(err)
		return
	}

	networkId, subnetId, err := CreateSubnet(client, config, vpcId)
	if err != nil {
		log.Print(err)
		return
	}

	securityGroupIds, err := CreateSecurityGroup(client, config)
	if err != nil {
		log.Print(err)
		return
	}

	//create server
	instance, err := CreateECS(client, config, vpcId, networkId, securityGroupIds)
	if err != nil {
		log.Print(err)
		return
	}

	//create loadBalancer
	elbId, err := CreateELB(client, subnetId, vpcId, config)
	if err != nil {
		log.Print(err)
		return
	}

	listener80, err := CreateListener(client, elbId, 80)
	if err != nil {
		log.Print(err)
		return
	}
	listener443, err := CreateListener(client, elbId, 443)
	if err != nil {
		log.Print(err)
		return
	}

	pool80, err := CreatePool(client, listener80, vpcId)
	if err != nil {
		log.Print(err)
		return
	}

	pool443, err := CreatePool(client, listener443, vpcId)
	if err != nil {
		log.Print(err)
		return
	}

	elbMembers := make([]HuaweiElbMembers, 0)
	elbMembers = append(elbMembers, HuaweiElbMembers{ID: pool80, Port: 80})
	elbMembers = append(elbMembers, HuaweiElbMembers{ID: pool443, Port: 443})
	err = AddServerToElb(client, subnetId, elbMembers, instance)
	if err != nil {
		log.Print(err)
		return
	}
}
