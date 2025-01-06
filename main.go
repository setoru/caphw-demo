package main

import (
	"fmt"
)

func main() {
	client, err := GetHuaweiClient()
	if err != nil {
		fmt.Println(err)
		return
	}

	config, err := parseConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	vpcId, err := CreateVpc(client, config)
	if err != nil {
		fmt.Println(err)
		return
	}

	networkId, subnetId, err := CreateSubnet(client, config, vpcId)
	if err != nil {
		fmt.Println(err)
		return
	}

	securityGroupIds, err := CreateSecurityGroup(client, config)
	if err != nil {
		fmt.Println(err)
		return
	}

	instance, err := CreateECS(client, config, vpcId, networkId, securityGroupIds)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = AddServerToElb(client, subnetId, config.ElbMembers, instance)
	if err != nil {
		fmt.Println(err)
		return
	}

}
