package pb_helper

import (
	"github.com/golang/protobuf/proto"
	"gitlab.alipay-inc.com/octopus/tbase-go-client/pb/base"
	"gitlab.alipay-inc.com/octopus/tbase-go-client/pb/coordinator"
)

func ConstructGetCoordinatorStatRequestMessage(requestId string) (*alibaba_proto.AliMessage) {
	aliMessage := new(alibaba_proto.AliMessage)
	aliMessage.SessionNo = proto.String(requestId)
	messageType := alibaba_coordinator_proto.MessageType_GET_COORDINATOR_STAT_REQUEST
	aliMessage.MessageType = proto.Int32(int32(messageType))
	getCoordiantorStatRequest := new(alibaba_coordinator_proto.GetCoordinatorStatRequest)
	proto.SetExtension(aliMessage, alibaba_coordinator_proto.E_GetCoordinatorStatRequest_Request, getCoordiantorStatRequest)
	return aliMessage
}

func ConstructGetClusterLayoutRequestMessage(requestId string, cluster string, tenant string, version int32) (*alibaba_proto.AliMessage) {
	aliMessage := new(alibaba_proto.AliMessage)
	aliMessage.SessionNo = proto.String(requestId)
	messageType := alibaba_coordinator_proto.MessageType_GET_LAYOUT_REQUEST
	aliMessage.MessageType = proto.Int32(int32(messageType))
	getLayoutRequest := new(alibaba_coordinator_proto.GetLayoutRequest)
	getLayoutRequest.TenantId = proto.String(tenant)
	getLayoutRequest.ClusterName = proto.String(cluster)
	layoutInfo := new(alibaba_coordinator_proto.LayoutInfo)
	layoutInfo.Version = proto.Int32(version)
	getLayoutRequest.LayoutInfo = layoutInfo
	proto.SetExtension(aliMessage, alibaba_coordinator_proto.E_GetLayoutRequest_Request, getLayoutRequest)
	return aliMessage
}
