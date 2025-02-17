package main

import (
	"fmt"
	"os"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/config"
	reg "github.com/huaweicloud/huaweicloud-sdk-go-v3/core/region"
	ecs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2"
	ecsMdl "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
	eip "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/eip/v2"
	eipMdl "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/eip/v2/model"
	elb "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3"
	elbMdl "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3/model"
	ims "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ims/v2"
	imsMdl "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ims/v2/model"
	nat "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/nat/v2"
	natMdl "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/nat/v2/model"
	vpc "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v2"
	vpcMdl "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v2/model"
)

type Client interface {
	//ECS
	CreateServers(request *ecsMdl.CreateServersRequest) (*ecsMdl.CreateServersResponse, error)
	ShowServer(request *ecsMdl.ShowServerRequest) (*ecsMdl.ShowServerResponse, error)
	ListServersDetails(request *ecsMdl.ListServersDetailsRequest) (*ecsMdl.ListServersDetailsResponse, error)
	DeleteServers(request *ecsMdl.DeleteServersRequest) (*ecsMdl.DeleteServersResponse, error)
	ShowServerBlockDevice(request *ecsMdl.ShowServerBlockDeviceRequest) (*ecsMdl.ShowServerBlockDeviceResponse, error)
	CreateServerGroup(request *ecsMdl.CreateServerGroupRequest) (*ecsMdl.CreateServerGroupResponse, error)
	DeleteServerGroup(request *ecsMdl.DeleteServerGroupRequest) (*ecsMdl.DeleteServerGroupResponse, error)
	BatchCreateServerTags(request *ecsMdl.BatchCreateServerTagsRequest) (*ecsMdl.BatchCreateServerTagsResponse, error)
	BatchStopServers(request *ecsMdl.BatchStopServersRequest) (*ecsMdl.BatchStopServersResponse, error)
	ListFlavors(request *ecsMdl.ListFlavorsRequest) (*ecsMdl.ListFlavorsResponse, error)
	//VPC
	CreateVpc(request *vpcMdl.CreateVpcRequest) (*vpcMdl.CreateVpcResponse, error)
	DeleteVpc(request *vpcMdl.DeleteVpcRequest) (*vpcMdl.DeleteVpcResponse, error)
	ShowVpc(request *vpcMdl.ShowVpcRequest) (*vpcMdl.ShowVpcResponse, error)
	ListSecurityGroups(request *vpcMdl.ListSecurityGroupsRequest) (*vpcMdl.ListSecurityGroupsResponse, error)
	CreateSecurityGroup(request *vpcMdl.CreateSecurityGroupRequest) (*vpcMdl.CreateSecurityGroupResponse, error)
	DeleteSecurityGroup(request *vpcMdl.DeleteSecurityGroupRequest) (*vpcMdl.DeleteSecurityGroupResponse, error)
	CreateSubnet(request *vpcMdl.CreateSubnetRequest) (*vpcMdl.CreateSubnetResponse, error)
	DeleteSubnet(request *vpcMdl.DeleteSubnetRequest) (*vpcMdl.DeleteSubnetResponse, error)
	ListSecurityGroupRules(request *vpcMdl.ListSecurityGroupRulesRequest) (*vpcMdl.ListSecurityGroupRulesResponse, error)
	DeleteSecurityGroupRule(request *vpcMdl.DeleteSecurityGroupRuleRequest) (*vpcMdl.DeleteSecurityGroupRuleResponse, error)
	//ELB
	ShowLoadBalancer(request *elbMdl.ShowLoadBalancerRequest) (*elbMdl.ShowLoadBalancerResponse, error)
	CreateLoadBalancer(request *elbMdl.CreateLoadBalancerRequest) (*elbMdl.CreateLoadBalancerResponse, error)
	DeleteLoadBalancer(request *elbMdl.DeleteLoadBalancerRequest) (*elbMdl.DeleteLoadBalancerResponse, error)
	BatchCreateMembers(request *elbMdl.BatchCreateMembersRequest) (*elbMdl.BatchCreateMembersResponse, error)
	BatchDeleteMembers(request *elbMdl.BatchDeleteMembersRequest) (*elbMdl.BatchDeleteMembersResponse, error)
	ListAvailabilityZones(request *elbMdl.ListAvailabilityZonesRequest) (*elbMdl.ListAvailabilityZonesResponse, error)
	CreateListener(request *elbMdl.CreateListenerRequest) (*elbMdl.CreateListenerResponse, error)
	DeleteListener(request *elbMdl.DeleteListenerRequest) (*elbMdl.DeleteListenerResponse, error)
	CreatePool(request *elbMdl.CreatePoolRequest) (*elbMdl.CreatePoolResponse, error)
	DeletePool(request *elbMdl.DeletePoolRequest) (*elbMdl.DeletePoolResponse, error)
	//IMS
	ListImages(request *imsMdl.ListImagesRequest) (*imsMdl.ListImagesResponse, error)
	//NAT
	CreateNatGateway(request *natMdl.CreateNatGatewayRequest) (*natMdl.CreateNatGatewayResponse, error)
	CreateNatGatewaySnatRule(request *natMdl.CreateNatGatewaySnatRuleRequest) (*natMdl.CreateNatGatewaySnatRuleResponse, error)
	//EIP
	CreatePublicip(request *eipMdl.CreatePublicipRequest) (*eipMdl.CreatePublicipResponse, error)
}

type HuaweiCloudClient struct {
	EcsClient *ecs.EcsClient
	VpcClient *vpc.VpcClient
	ElbClient *elb.ElbClient
	ImsClient *ims.ImsClient
	NatClient *nat.NatClient
	EipClient *eip.EipClient
}

func (client *HuaweiCloudClient) ShowServerBlockDevice(request *ecsMdl.ShowServerBlockDeviceRequest) (*ecsMdl.ShowServerBlockDeviceResponse, error) {
	return client.EcsClient.ShowServerBlockDevice(request)
}

func (client *HuaweiCloudClient) CreateServerGroup(request *ecsMdl.CreateServerGroupRequest) (*ecsMdl.CreateServerGroupResponse, error) {
	return client.EcsClient.CreateServerGroup(request)
}

func (client *HuaweiCloudClient) DeleteServerGroup(request *ecsMdl.DeleteServerGroupRequest) (*ecsMdl.DeleteServerGroupResponse, error) {
	return client.EcsClient.DeleteServerGroup(request)
}

func (client *HuaweiCloudClient) BatchCreateServerTags(request *ecsMdl.BatchCreateServerTagsRequest) (*ecsMdl.BatchCreateServerTagsResponse, error) {
	return client.EcsClient.BatchCreateServerTags(request)
}

func (client *HuaweiCloudClient) CreateVpc(request *vpcMdl.CreateVpcRequest) (*vpcMdl.CreateVpcResponse, error) {
	return client.VpcClient.CreateVpc(request)
}

func (client *HuaweiCloudClient) DeleteVpc(request *vpcMdl.DeleteVpcRequest) (*vpcMdl.DeleteVpcResponse, error) {
	return client.VpcClient.DeleteVpc(request)
}

func (client *HuaweiCloudClient) ShowVpc(request *vpcMdl.ShowVpcRequest) (*vpcMdl.ShowVpcResponse, error) {
	return client.VpcClient.ShowVpc(request)
}

func (client *HuaweiCloudClient) CreateSecurityGroup(request *vpcMdl.CreateSecurityGroupRequest) (*vpcMdl.CreateSecurityGroupResponse, error) {
	return client.VpcClient.CreateSecurityGroup(request)
}

func (client *HuaweiCloudClient) DeleteSecurityGroup(request *vpcMdl.DeleteSecurityGroupRequest) (*vpcMdl.DeleteSecurityGroupResponse, error) {
	return client.VpcClient.DeleteSecurityGroup(request)
}

func (client *HuaweiCloudClient) DeleteSecurityGroupRule(request *vpcMdl.DeleteSecurityGroupRuleRequest) (*vpcMdl.DeleteSecurityGroupRuleResponse, error) {
	return client.VpcClient.DeleteSecurityGroupRule(request)
}

func (client *HuaweiCloudClient) ListSecurityGroupRules(request *vpcMdl.ListSecurityGroupRulesRequest) (*vpcMdl.ListSecurityGroupRulesResponse, error) {
	return client.VpcClient.ListSecurityGroupRules(request)
}

func (client *HuaweiCloudClient) CreateSubnet(request *vpcMdl.CreateSubnetRequest) (*vpcMdl.CreateSubnetResponse, error) {
	return client.VpcClient.CreateSubnet(request)
}

func (client *HuaweiCloudClient) DeleteSubnet(request *vpcMdl.DeleteSubnetRequest) (*vpcMdl.DeleteSubnetResponse, error) {
	return client.VpcClient.DeleteSubnet(request)
}

func (client *HuaweiCloudClient) ShowLoadBalancer(request *elbMdl.ShowLoadBalancerRequest) (*elbMdl.ShowLoadBalancerResponse, error) {
	return client.ElbClient.ShowLoadBalancer(request)
}

func (client *HuaweiCloudClient) CreateLoadBalancer(request *elbMdl.CreateLoadBalancerRequest) (*elbMdl.CreateLoadBalancerResponse, error) {
	return client.ElbClient.CreateLoadBalancer(request)
}

func (client *HuaweiCloudClient) DeleteLoadBalancer(request *elbMdl.DeleteLoadBalancerRequest) (*elbMdl.DeleteLoadBalancerResponse, error) {
	return client.ElbClient.DeleteLoadBalancer(request)
}

func (client *HuaweiCloudClient) BatchCreateMembers(request *elbMdl.BatchCreateMembersRequest) (*elbMdl.BatchCreateMembersResponse, error) {
	return client.ElbClient.BatchCreateMembers(request)
}

func (client HuaweiCloudClient) BatchDeleteMembers(request *elbMdl.BatchDeleteMembersRequest) (*elbMdl.BatchDeleteMembersResponse, error) {
	return client.ElbClient.BatchDeleteMembers(request)
}

func (client HuaweiCloudClient) ListAvailabilityZones(request *elbMdl.ListAvailabilityZonesRequest) (*elbMdl.ListAvailabilityZonesResponse, error) {
	return client.ElbClient.ListAvailabilityZones(request)
}

func (client HuaweiCloudClient) CreateListener(request *elbMdl.CreateListenerRequest) (*elbMdl.CreateListenerResponse, error) {
	return client.ElbClient.CreateListener(request)
}

func (client HuaweiCloudClient) DeleteListener(request *elbMdl.DeleteListenerRequest) (*elbMdl.DeleteListenerResponse, error) {
	return client.ElbClient.DeleteListener(request)
}

func (client HuaweiCloudClient) DeletePool(request *elbMdl.DeletePoolRequest) (*elbMdl.DeletePoolResponse, error) {
	return client.ElbClient.DeletePool(request)
}

func (client HuaweiCloudClient) CreatePool(request *elbMdl.CreatePoolRequest) (*elbMdl.CreatePoolResponse, error) {
	return client.ElbClient.CreatePool(request)
}

func (client *HuaweiCloudClient) CreateServers(request *ecsMdl.CreateServersRequest) (*ecsMdl.CreateServersResponse, error) {
	return client.EcsClient.CreateServers(request)
}

func (client *HuaweiCloudClient) ShowServer(request *ecsMdl.ShowServerRequest) (*ecsMdl.ShowServerResponse, error) {
	return client.EcsClient.ShowServer(request)
}

func (client *HuaweiCloudClient) ListServersDetails(request *ecsMdl.ListServersDetailsRequest) (*ecsMdl.ListServersDetailsResponse, error) {
	return client.EcsClient.ListServersDetails(request)
}

func (client *HuaweiCloudClient) DeleteServers(request *ecsMdl.DeleteServersRequest) (*ecsMdl.DeleteServersResponse, error) {
	return client.EcsClient.DeleteServers(request)
}

func (client *HuaweiCloudClient) ListImages(request *imsMdl.ListImagesRequest) (*imsMdl.ListImagesResponse, error) {
	return client.ImsClient.ListImages(request)
}

func (client *HuaweiCloudClient) ListSecurityGroups(request *vpcMdl.ListSecurityGroupsRequest) (*vpcMdl.ListSecurityGroupsResponse, error) {
	return client.VpcClient.ListSecurityGroups(request)
}

func (client *HuaweiCloudClient) BatchStopServers(request *ecsMdl.BatchStopServersRequest) (*ecsMdl.BatchStopServersResponse, error) {
	return client.EcsClient.BatchStopServers(request)
}

func (client *HuaweiCloudClient) ListFlavors(request *ecsMdl.ListFlavorsRequest) (*ecsMdl.ListFlavorsResponse, error) {
	return client.EcsClient.ListFlavors(request)
}

func (client *HuaweiCloudClient) CreateNatGateway(request *natMdl.CreateNatGatewayRequest) (*natMdl.CreateNatGatewayResponse, error) {
	return client.NatClient.CreateNatGateway(request)
}

func (client *HuaweiCloudClient) CreateNatGatewaySnatRule(request *natMdl.CreateNatGatewaySnatRuleRequest) (*natMdl.CreateNatGatewaySnatRuleResponse, error) {
	return client.NatClient.CreateNatGatewaySnatRule(request)
}

func (client *HuaweiCloudClient) CreatePublicip(request *eipMdl.CreatePublicipRequest) (*eipMdl.CreatePublicipResponse, error) {
	return client.EipClient.CreatePublicip(request)
}

func GetHuaweiClient() (Client, error) {
	ak := os.Getenv("access_key")
	sk := os.Getenv("secret_key")
	region := os.Getenv("region")
	credential, err := basic.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		SafeBuild()
	if err != nil {
		return nil, err
	}

	ecsRegion := reg.NewRegion(region, fmt.Sprintf("https://ecs.%s.myhuaweicloud.com", region))
	ecsBuild, err := ecs.EcsClientBuilder().
		WithRegion(ecsRegion).
		WithCredential(credential).
		WithHttpConfig(config.DefaultHttpConfig()).
		SafeBuild()
	ecsClient := ecs.NewEcsClient(ecsBuild)

	vpcRegion := reg.NewRegion(region, fmt.Sprintf("https://vpc.%s.myhuaweicloud.com", region))
	vpcBuild, err := vpc.VpcClientBuilder().
		WithRegion(vpcRegion).
		WithCredential(credential).
		WithHttpConfig(config.DefaultHttpConfig()).
		SafeBuild()
	vpcClient := vpc.NewVpcClient(vpcBuild)

	elbRegion := reg.NewRegion(region, fmt.Sprintf("https://elb.%s.myhuaweicloud.com", region))
	elbBuild, err := elb.ElbClientBuilder().
		WithRegion(elbRegion).
		WithCredential(credential).
		WithHttpConfig(config.DefaultHttpConfig()).
		SafeBuild()
	elbClient := elb.NewElbClient(elbBuild)

	imsRegion := reg.NewRegion(region, fmt.Sprintf("https://ims.%s.myhuaweicloud.com", region))
	imsBuild, err := ims.ImsClientBuilder().
		WithRegion(imsRegion).
		WithCredential(credential).
		WithHttpConfig(config.DefaultHttpConfig()).
		SafeBuild()
	if err != nil {
		return nil, err
	}
	imsClient := ims.NewImsClient(imsBuild)

	natRegion := reg.NewRegion(region, fmt.Sprintf("https://nat.%s.myhuaweicloud.com", region))
	natBuild, err := nat.NatClientBuilder().
		WithRegion(natRegion).
		WithCredential(credential).
		WithHttpConfig(config.DefaultHttpConfig()).
		SafeBuild()
	if err != nil {
		return nil, err
	}
	natClient := nat.NewNatClient(natBuild)

	eipRegion := reg.NewRegion(region, fmt.Sprintf("https://eip.%s.myhuaweicloud.com", region))
	eipBuild, err := eip.EipClientBuilder().
		WithRegion(eipRegion).
		WithCredential(credential).
		WithHttpConfig(config.DefaultHttpConfig()).
		SafeBuild()
	if err != nil {
		return nil, err
	}
	eipClient := eip.NewEipClient(eipBuild)
	return &HuaweiCloudClient{
		EcsClient: ecsClient,
		VpcClient: vpcClient,
		ElbClient: elbClient,
		ImsClient: imsClient,
		NatClient: natClient,
		EipClient: eipClient,
	}, nil
}
