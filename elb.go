package main

import (
	ecsMdl "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/ecs/v2/model"
	elbMdl "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/elb/v3/model"
)

func AddServerToElb(client Client, subnetIpv4Id string, elbMembers []HuaweiElbMembers, servers []*ecsMdl.ServerDetail) (err error) {
	for _, member := range elbMembers {
		createMembersRequest := &elbMdl.BatchCreateMembersRequest{Body: &elbMdl.BatchCreateMembersRequestBody{Members: make([]elbMdl.BatchCreateMembersOption, 0)}}
		for _, server := range servers {
			for _, addresses := range server.Addresses {
				for _, address := range addresses {
					if *address.OSEXTIPStype == ecsMdl.GetServerAddressOSEXTIPStypeEnum().FIXED {
						createMembersRequest.PoolId = member.ID
						createMember := elbMdl.BatchCreateMembersOption{Address: address.Addr, ProtocolPort: member.Port, SubnetCidrId: &subnetIpv4Id}
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

func DeleteServerFromElb(client Client, elbMembers []HuaweiElbMembers, servers []*ecsMdl.ServerDetail) (err error) {
	for _, member := range elbMembers {
		deleteMembersRequest := &elbMdl.BatchDeleteMembersRequest{Body: &elbMdl.BatchDeleteMembersRequestBody{Members: make([]elbMdl.BatchDeleteMembersOption, 0)}}
		for _, server := range servers {
			for _, addresses := range server.Addresses {
				for _, address := range addresses {
					if *address.OSEXTIPStype == ecsMdl.GetServerAddressOSEXTIPStypeEnum().FIXED {
						deleteMembersRequest.PoolId = member.ID
						addr := address.Addr
						deleteMember := elbMdl.BatchDeleteMembersOption{Address: &addr, ProtocolPort: &member.Port}
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
