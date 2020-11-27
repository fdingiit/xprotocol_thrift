package rpc_helper

import (
	"net"
	"time"

	tbase_log "gitlab.alipay-inc.com/octopus/tbase-go-client/tloggger"

	"github.com/golang/protobuf/proto"

	error2 "gitlab.alipay-inc.com/octopus/tbase-go-client/model/error"

	message_helper "gitlab.alipay-inc.com/octopus/tbase-go-client/utils/message.helper"
)

func CreateConnectionAndSendRequest(requestId string, server string, timeout time.Duration, data []byte) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", server, timeout)
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[RPC_HELPER] dial to %v with timeout %v error. error: %v", server, timeout, err)
		return nil, err
	}

	err = conn.SetDeadline(time.Now().Add(timeout))
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[RPC_HELPER] set deadline for connection error. dest server: %v, timeout: %v, error: %v", server, timeout, err)
		return nil, err
	}

	err = SendRequest(conn, data)
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[RPC_HELPER] send request error. dest server: %v, error: %v", server, err)
		return nil, err
	}

	return conn, nil
}

func SendRequest(conn net.Conn, b []byte) error {
	var writeLen = 0

	// timeout will ensure no dead loop
	for writeLen < len(b) {
		tmpWriteLen, err := conn.Write(b[writeLen:])
		if err != nil {
			tbase_log.TBaseLogger.Errorf("[RPC_HELPER] write request buffer error. error: %v", err)
			return err
		}
		writeLen += tmpWriteLen
	}

	return nil
}

func GetResponse(conn net.Conn, requestId string, extension *proto.ExtensionDesc) (interface{}, error) {
	resp, err := message_helper.Decode(conn)
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[RPC_HELPER] read and decode message error. requestId: %v, error: %v", requestId, err)
		return nil, err
	}
	if resp.GetSessionNo() != requestId {
		tbase_log.TBaseLogger.Errorf("[RPC_HELPER] sent requestId and received requestId mismatch. sent requestId: %v, received requestId: %v", requestId, resp.GetSessionNo())
		return nil, error2.NewTBaseClientInternalError("sent requestId and received requestId mismatch")
	}

	result, err := proto.GetExtension(resp, extension)
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[RPC_HELPER] get \"%v\" extension from \"AliMessage\" error. error: %v", extension, err)
		return nil, err
	}
	if result == nil {
		tbase_log.TBaseLogger.Errorf("[RPC_HELPER] can't get extension from response. response: %v", resp)
		return nil, error2.NewTBaseClientInternalError("can't get extension from response")
	}

	return result, nil
}
