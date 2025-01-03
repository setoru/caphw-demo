package caphw_demo

import (
	"fmt"
	"log"
	"os"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/config"
	reg "github.com/huaweicloud/huaweicloud-sdk-go-v3/core/region"
	ecs "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2"
	ecsMdl "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
	elb "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3"
	elbMdl "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3/model"
	ims "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ims/v2"
	imsMdl "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ims/v2/model"
	vpc "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3"
	vpcMdl "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vpc/v3/model"
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
	NovaAssociateSecurityGroup(request *ecsMdl.NovaAssociateSecurityGroupRequest) (*ecsMdl.NovaAssociateSecurityGroupResponse, error)
	NovaDisassociateSecurityGroup(request *ecsMdl.NovaDisassociateSecurityGroupRequest) (*ecsMdl.NovaDisassociateSecurityGroupResponse, error)
	NovaListAvailabilityZones(request *ecsMdl.NovaListAvailabilityZonesRequest) (*ecsMdl.NovaListAvailabilityZonesResponse, error)
	NovaListServerSecurityGroups(request *ecsMdl.NovaListServerSecurityGroupsRequest) (*ecsMdl.NovaListServerSecurityGroupsResponse, error)
	BatchStopServers(request *ecsMdl.BatchStopServersRequest) (*ecsMdl.BatchStopServersResponse, error)
	ListFlavors(request *ecsMdl.ListFlavorsRequest) (*ecsMdl.ListFlavorsResponse, error)
	//VPC
	CreateVpc(request *vpcMdl.CreateVpcRequest) (*vpcMdl.CreateVpcResponse, error)
	DeleteVpc(request *vpcMdl.DeleteVpcRequest) (*vpcMdl.DeleteVpcResponse, error)
	ShowVpc(request *vpcMdl.ShowVpcRequest) (*vpcMdl.ShowVpcResponse, error)
	ListSecurityGroups(request *vpcMdl.ListSecurityGroupsRequest) (*vpcMdl.ListSecurityGroupsResponse, error)
	//ELB
	ShowLoadBalancer(request *elbMdl.ShowLoadBalancerRequest) (*elbMdl.ShowLoadBalancerResponse, error)
	CreateLoadBalancer(request *elbMdl.CreateLoadBalancerRequest) (*elbMdl.CreateLoadBalancerResponse, error)
	DeleteLoadBalancer(request *elbMdl.DeleteLoadBalancerRequest) (*elbMdl.DeleteLoadBalancerResponse, error)
	BatchCreateMembers(request *elbMdl.BatchCreateMembersRequest) (*elbMdl.BatchCreateMembersResponse, error)
	BatchDeleteMembers(request *elbMdl.BatchDeleteMembersRequest) (*elbMdl.BatchDeleteMembersResponse, error)
	//IMS
	ListImages(request *imsMdl.ListImagesRequest) (*imsMdl.ListImagesResponse, error)
}

type HuaweiCloudClient struct {
	EcsClient *ecs.EcsClient
	VpcClient *vpc.VpcClient
	ElbClient *elb.ElbClient
	ImsClient *ims.ImsClient
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

func (client *HuaweiCloudClient) NovaAssociateSecurityGroup(request *ecsMdl.NovaAssociateSecurityGroupRequest) (*ecsMdl.NovaAssociateSecurityGroupResponse, error) {
	return client.EcsClient.NovaAssociateSecurityGroup(request)
}

func (client *HuaweiCloudClient) NovaDisassociateSecurityGroup(request *ecsMdl.NovaDisassociateSecurityGroupRequest) (*ecsMdl.NovaDisassociateSecurityGroupResponse, error) {
	return client.EcsClient.NovaDisassociateSecurityGroup(request)
}

func (client *HuaweiCloudClient) NovaListAvailabilityZones(request *ecsMdl.NovaListAvailabilityZonesRequest) (*ecsMdl.NovaListAvailabilityZonesResponse, error) {
	return client.EcsClient.NovaListAvailabilityZones(request)
}

func (client *HuaweiCloudClient) NovaListServerSecurityGroups(request *ecsMdl.NovaListServerSecurityGroupsRequest) (*ecsMdl.NovaListServerSecurityGroupsResponse, error) {
	return client.EcsClient.NovaListServerSecurityGroups(request)
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

func GetHuaweiClient() (Client, error) {
	ak := os.Getenv("access_key")
	sk := os.Getenv("secret_key")
	region := os.Getenv("region")
	credential, err := basic.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		SafeBuild()
	if err != nil {
		log.Print(err)
	}

	ecsRegion := reg.NewRegion(region, fmt.Sprintf("https://ecs.%s.myhuaweicloud.com", region))
	ecsBuild, err := ecs.EcsClientBuilder().
		WithRegion(ecsRegion).
		WithCredential(credential).
		WithHttpConfig(config.DefaultHttpConfig()).
		SafeBuild()
	if err != nil {
		log.Print(err)
	}
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
	imsClient := ims.NewImsClient(imsBuild)

	return &HuaweiCloudClient{
		EcsClient: ecsClient,
		VpcClient: vpcClient,
		ElbClient: elbClient,
		ImsClient: imsClient,
	}, nil
}
