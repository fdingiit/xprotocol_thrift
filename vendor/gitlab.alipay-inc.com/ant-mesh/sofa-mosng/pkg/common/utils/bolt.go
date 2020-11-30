package utils

import (
	"mosn.io/mosn/pkg/protocol/xprotocol/bolt"
)

const (
	SofaRequestClassName    = "com.alipay.sofa.rpc.core.request.SofaRequest"
	SofaResponseClassName   = "com.alipay.sofa.rpc.core.response.SofaResponse"
	SofaHeaderTargetService = "sofa_head_target_service"
	SofaHeaderService       = "service"
	SofaHeaderMethodName    = "sofa_head_method_name"
	SofaHeaderUid           = "uid"
	SofaHeaderVipAddr       = "mesh_vip_address"
	SofaHeaderTargetLocal   = "sofa_head_target_local"
)

func BuildSofaRequest(headerMap map[string]string, contentLen int) bolt.Request {
	return BuildBoltRequest(headerMap, contentLen, SofaRequestClassName, 0)
}

func BuildSofaRequestWithTimeout(headerMap map[string]string, contentLen int, timeout int) bolt.Request {
	return BuildBoltRequest(headerMap, contentLen, SofaRequestClassName, timeout)
}

func BuildBoltRequest(headerMap map[string]string, contentLen int, className string, timeout int) bolt.Request {
	command := buildBasicRequestCommand()
	command.Class = className
	for key, value := range headerMap {
		command.Set(key, value)
	}
	command.ContentLen = uint32(contentLen) // TODO: 这个应该可以考虑不要了, 可以核对一下
	if timeout > 0 {
		command.Timeout = int32(timeout)
	}
	return command
}

func buildBasicRequestCommand() bolt.Request {
	command := bolt.Request{
		RequestHeader: bolt.RequestHeader{
			Protocol: bolt.ProtocolCode,
			CmdType:  bolt.CmdTypeRequest,
			CmdCode:  bolt.CmdCodeRpcRequest,
			Version:  1,
			Codec:    11,
			Timeout:  3000,
		},
	}
	return command
}
