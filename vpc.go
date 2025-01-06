package main

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v2/model"
)

func CreateVpc(client Client, config *Config) (string, error) {
	if config.Vpc.ID != "" {
		return config.Vpc.ID, nil
	}
	req := &model.CreateVpcRequest{Body: &model.CreateVpcRequestBody{Vpc: &model.CreateVpcOption{}}}
	req.Body.Vpc.Name = &config.Vpc.Name
	req.Body.Vpc.Cidr = &config.Vpc.Cidr
	req.Body.Vpc.EnterpriseProjectId = config.Vpc.EnterpriseProjectId
	response, err := client.CreateVpc(req)
	if err != nil {
		return "", err
	}
	return response.Vpc.Id, nil
}

func CreateSecurityGroup(client Client, config *Config) ([]string, error) {
	securityGroupIDs, _ := getSecurityGroupIDs(config)
	for _, group := range config.SecurityGroups {
		if group.ID != "" {
			continue
		}
		req := &model.CreateSecurityGroupRequest{Body: &model.CreateSecurityGroupRequestBody{SecurityGroup: &model.CreateSecurityGroupOption{}}}
		req.Body.SecurityGroup.Name = group.Name
		req.Body.SecurityGroup.EnterpriseProjectId = group.EnterpriseProjectId
		response, err := client.CreateSecurityGroup(req)
		if err != nil {
			return nil, err
		}
		securityGroupIDs = append(securityGroupIDs, response.SecurityGroup.Id)
	}
	return securityGroupIDs, nil
}

func getSecurityGroupIDs(config *Config) ([]string, error) {
	var securityGroupIDs []string

	if len(config.SecurityGroups) == 0 {
		return []string{}, nil
	}

	for _, sg := range config.SecurityGroups {
		if sg.ID == "" {
			continue
		}
		securityGroupIDs = append(securityGroupIDs, sg.ID)
	}
	return securityGroupIDs, nil
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
	return response.Subnet.NeutronNetworkId, response.Subnet.NeutronSubnetId, nil
}
