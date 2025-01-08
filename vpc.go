package main

import (
	"fmt"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v2/model"
)

func CreateVpc(client Client, config *Config) (string, error) {
	if config.Vpc.Id != "" {
		return config.Vpc.Id, nil
	}
	req := &model.CreateVpcRequest{Body: &model.CreateVpcRequestBody{Vpc: &model.CreateVpcOption{}}}
	req.Body.Vpc.Name = &config.Vpc.Name
	req.Body.Vpc.Cidr = &config.Vpc.Cidr
	req.Body.Vpc.EnterpriseProjectId = config.Vpc.EnterpriseProjectId
	response, err := client.CreateVpc(req)
	if err != nil {
		return "", err
	}
	fmt.Println("create vpc success")
	return response.Vpc.Id, nil
}

func CreateSecurityGroup(client Client, config *Config) ([]string, error) {
	securityGroupIds := getSecurityGroupIds(config)
	for _, group := range config.SecurityGroups {
		if group.Id != "" {
			continue
		}
		req := &model.CreateSecurityGroupRequest{Body: &model.CreateSecurityGroupRequestBody{SecurityGroup: &model.CreateSecurityGroupOption{}}}
		req.Body.SecurityGroup.Name = group.Name
		req.Body.SecurityGroup.EnterpriseProjectId = group.EnterpriseProjectId
		response, err := client.CreateSecurityGroup(req)
		if err != nil {
			return nil, err
		}
		securityGroupIds = append(securityGroupIds, response.SecurityGroup.Id)
	}
	fmt.Println("create security group success")
	return securityGroupIds, nil
}

func getSecurityGroupIds(config *Config) []string {
	securityGroupIds := make([]string, 0)

	if len(config.SecurityGroups) == 0 {
		return nil
	}

	for _, sg := range config.SecurityGroups {
		if sg.Id == "" {
			continue
		}
		securityGroupIds = append(securityGroupIds, sg.Id)
	}
	return securityGroupIds
}

func getSecurityGroupIdsSet(config *Config) map[string]struct{} {
	ids := getSecurityGroupIds(config)
	set := make(map[string]struct{})
	for _, id := range ids {
		set[id] = struct{}{}
	}
	return set
}

func CreateSubnet(client Client, config *Config, vpcId string) (string, string, error) {
	if config.Subnet.SubnetId != "" || config.Subnet.NetworkId != "" {
		return config.Subnet.SubnetId, config.Subnet.NetworkId, nil
	}
	req := &model.CreateSubnetRequest{Body: &model.CreateSubnetRequestBody{Subnet: &model.CreateSubnetOption{}}}
	req.Body.Subnet.Cidr = config.Subnet.Cidr
	req.Body.Subnet.VpcId = vpcId
	req.Body.Subnet.GatewayIp = "192.168.0.1"
	req.Body.Subnet.Name = config.Subnet.Name
	response, err := client.CreateSubnet(req)
	if err != nil {
		return "", "", err
	}
	fmt.Println("create subnet success")
	return response.Subnet.NeutronNetworkId, response.Subnet.NeutronSubnetId, nil
}

func DeleteVpc(client Client, config *Config, vpcId string) error {
	if config.Vpc.Id != "" {
		return nil
	}
	req := &model.DeleteVpcRequest{VpcId: vpcId}
	_, err := client.DeleteVpc(req)
	if err != nil {
		return err
	}
	return nil
}

func DeleteSubnet(client Client, config *Config, networkId, vpcId string) error {
	if config.Subnet.SubnetId != "" || config.Subnet.NetworkId != "" {
		return nil
	}
	req := &model.DeleteSubnetRequest{SubnetId: networkId, VpcId: vpcId}
	_, err := client.DeleteSubnet(req)
	if err != nil {
		return err
	}
	return nil
}

func DeleteSecurityGroup(client Client, config *Config, securityGroupIds []string) error {
	set := getSecurityGroupIdsSet(config)
	for _, id := range securityGroupIds {
		if _, ok := set[id]; ok {
			continue
		}
		req := &model.DeleteSecurityGroupRequest{SecurityGroupId: id}
		_, err := client.DeleteSecurityGroup(req)
		if err != nil {
			return err
		}
	}
	return nil
}
