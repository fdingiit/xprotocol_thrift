package model

func ConvertToHeartbeatRequestPb(request *HeartbeatRequest) *HeartbeatRequestPb {
	requestPb := &HeartbeatRequestPb{
		Zone:          request.Zone,
		ClientIp:      request.ClientIp,
		InstanceId:    request.InstanceId,
		VersionMap:    request.VersionMap,
		Profile:       request.Profile,
		AckVersionMap: request.AckVersionMap,
	}

	return requestPb
}

func ConvertToHeartbeatResponse(responsePb *HeartbeatResponsePb) *HeartbeatResponse {
	response := &HeartbeatResponse{
		WaitTime: responsePb.WaitTime,
		DiffMap:  responsePb.DiffMap,
	}
	return response
}

func ConvertToSubscriberRegReqPb(request *SubscriberRegReq) *SubscriberRegReqPb {
	baseInfoPb := &BaseInfoPb{
		Zone:       request.Zone,
		DataId:     request.DataId,
		Uuid:       request.Uuid,
		InstanceId: request.InstanceId,
		Attributes: request.Attributes,
		Profile:    request.Profile,
	}
	requestPb := &SubscriberRegReqPb{
		BaseInfo: baseInfoPb,
	}
	return requestPb
}

func ConvertToSubscriberRegResult(responsePb *SubscriberRegResultPb) *SubscriberRegResult {
	response := &SubscriberRegResult{
		Result:  responsePb.Result,
		Message: responsePb.Message,
	}
	if responsePb.BaseInfo != nil {
		baseInfo := responsePb.BaseInfo
		response.Zone = baseInfo.Zone
		response.DataId = baseInfo.DataId
		response.Uuid = baseInfo.Uuid
		response.InstanceId = baseInfo.InstanceId
		response.Attributes = baseInfo.Attributes
		response.Profile = baseInfo.Profile
	}
	return response
}

func ConvertToAttributeSetRequest(requestPb *AttributeSetRequestPb) *AttributeSetRequest {
	request := &AttributeSetRequest{
		DataId: requestPb.DataId,
		Value:  requestPb.Value,
	}
	return request
}
