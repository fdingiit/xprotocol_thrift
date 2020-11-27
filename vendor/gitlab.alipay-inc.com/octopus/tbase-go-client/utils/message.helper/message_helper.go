package message_helper

import (
	"bufio"
	"fmt"
	"io"
	"net"

	tbase_log "gitlab.alipay-inc.com/octopus/tbase-go-client/tloggger"

	"github.com/golang/protobuf/proto"
	error2 "gitlab.alipay-inc.com/octopus/tbase-go-client/model/error"
	alibaba_proto "gitlab.alipay-inc.com/octopus/tbase-go-client/pb/base"
)

func Encode(message proto.Message) ([]byte, error) {
	size := proto.Size(message)
	sizeBytes := proto.EncodeVarint(uint64(size))
	data, err := proto.Marshal(message)
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[MESSAGE_HELPER] marshal proto message error. proto message: %v, error: %v", message, err)
		return nil, err
	}
	data = append(sizeBytes, data...)
	return data, nil
}

func Decode(conn net.Conn) (*alibaba_proto.AliMessage, error) {
	respLenBuf := make([]byte, 0)
	respLenByte := make([]byte, 1)
	var respLen uint64
	for i := 0; i < 5; i++ {
		l, err := conn.Read(respLenByte)
		if err != nil {
			tbase_log.TBaseLogger.Errorf("[MESSAGE_HELPER] read response data from Conn error. error: %v", err)
			return nil, err
		}
		if l <= 0 {
			tbase_log.TBaseLogger.Errorf("[MESSAGE_HELPER] read response data from Conn error. expect data length > 0")
			return nil, error2.NewTBaseClientInternalError("read response data from Conn error. expect data length > 0")
		}

		respLenBuf = append(respLenBuf, respLenByte[0])
		if respLenByte[0] <= 127 {
			respLen, _ = proto.DecodeVarint(respLenBuf[:i+1])
			break
		}
	}

	if respLen <= 0 {
		tbase_log.TBaseLogger.Errorf("[MESSAGE_HELPER] response body len <= 0")
		return nil, error2.NewTBaseClientInternalError("response body len <= 0")
	}

	respBodyBuf := make([]byte, respLen)
	reader := bufio.NewReader(conn)
	respBodyBufLen, err := io.ReadFull(reader, respBodyBuf)
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[MESSAGE_HELPER] read response body from Conn error. error: %v", err)
		return nil, err
	}

	if uint64(respBodyBufLen) != respLen {
		tbase_log.TBaseLogger.Errorf("[MESSAGE_HELPER] expect length of reponse body buffer %v, actually read %v length buffer", respLen, respBodyBufLen)
		return nil, error2.NewTBaseClientInternalError(fmt.Sprintf("expect length of reponse body buffer %v, actually read %v length buffer", respLen, respBodyBufLen))
	}

	resp := &alibaba_proto.AliMessage{}
	err = proto.Unmarshal(respBodyBuf, resp)
	return resp, err
}
