/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"net/http"

	"mosn.io/mosn/pkg/log"
	"mosn.io/mosn/pkg/protocol/xprotocol"
	"mosn.io/mosn/pkg/types"
	"mosn.io/pkg/buffer"
)

/**
 * Request command protocol for v1
 * 0     1     2           4           6           8          10           12          14         16
 * +-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+
 * |proto| type| cmdcode   |ver2 |   requestID           |codec|        timeout        |  classLen |
 * +-----------+-----------+-----------+-----------+-----------+-----------+-----------+-----------+
 * |headerLen  | contentLen            |                             ... ...                       |
 * +-----------+-----------+-----------+                                                                                               +
 * |               className + header  + content  bytes                                            |
 * +                                                                                               +
 * |                               ... ...                                                         |
 * +-----------------------------------------------------------------------------------------------+
 *
 * proto: code for protocol
 * type: request/response/request oneway
 * cmdcode: code for remoting command
 * ver2:version for remoting command
 * requestID: id of request
 * codec: code for codec
 * headerLen: length of header
 * contentLen: length of content
 *
 * Response command protocol for v1
 * 0     1     2     3     4           6           8          10           12          14         16
 * +-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+-----+
 * |proto| type| cmdcode   |ver2 |   requestID           |codec|respstatus |  classLen |headerLen  |
 * +-----------+-----------+-----------+-----------+-----------+-----------+-----------+-----------+
 * | contentLen            |                  ... ...                                              |
 * +-----------------------+                                                                       +
 * |                         className + header  + content  bytes                                  |
 * +                                                                                               +
 * |                               ... ...                                                         |
 * +-----------------------------------------------------------------------------------------------+
 * respstatus: response status
 */

func init() {
	xprotocol.RegisterProtocol(ProtocolName, &ThriftProtocol{})
}

const (
	MessageUnknown       byte = 0
	MessageTypeCall           = 1
	MessageTypeReply          = 2
	MessageTypeException      = 3
	MessageTypeOneway         = 4
)

type Message struct {
	Version string
	Type    byte
	NameLen int32
	Name    string
	SeqId   int32
}

type ThriftProtocol struct{}

// types.Protocol
func (proto *ThriftProtocol) Name() types.ProtocolName {
	return ProtocolName
}

func (proto *ThriftProtocol) Encode(ctx context.Context, model interface{}) (types.IoBuffer, error) {
	ProtocolLogger.Infof("Encode: %+v\n", model)

	switch frame := model.(type) {
	case *Request:
		return encodeRequest(ctx, frame)
	case *Response:
		return encodeResponse(ctx, frame)
	default:
		log.Proxy.Errorf(ctx, "[protocol][thrift] encode with unknown command : %+v", model)
		return nil, xprotocol.ErrUnknownType
	}
}

func messageType(b byte) byte {
	switch b {
	case 1:
		return MessageTypeCall
	case 2:
		return MessageTypeReply
	case 3:
		return MessageTypeException
	case 4:
		return MessageTypeOneway
	default:
		return MessageUnknown
	}
}

func (proto *ThriftProtocol) decodeMessage(ctx context.Context, data types.IoBuffer) (Message, error) {
	bytesLen := data.Len()
	bytes := data.Bytes()

	// 1. least bytes to decode header is RequestHeaderLen(22)
	if bytesLen < MessageMinLen {
		// todo
		panic("implement me")
	}

	// 2. least bytes to decode whole frame
	messageType := messageType(bytes[3])
	nameLen := binary.BigEndian.Uint32(bytes[4:8])
	name := string(bytes[8 : 8+nameLen])
	seqId := binary.BigEndian.Uint32(bytes[8+nameLen : 8+nameLen+4])

	return Message{
		Version: "", // todo
		Type:    messageType,
		NameLen: int32(nameLen),
		Name:    name,
		SeqId:   int32(seqId),
	}, nil
}

func (proto *ThriftProtocol) Decode(ctx context.Context, data types.IoBuffer) (interface{}, error) {
	ProtocolLogger.Infof("Decode: %+v\n", data)

	if data.Len() >= MessageMinLen {
		message, err := proto.decodeMessage(ctx, data)
		if err != nil {
			return nil, err
		}

		cmdType := data.Bytes()[0]

		switch message.Type {
		case MessageTypeCall:
			return decodeRequest(ctx, data, false, message)
		case MessageTypeOneway:
			return decodeRequest(ctx, data, true, message)
		case MessageTypeReply:
			return decodeResponse(ctx, data, message)
		case MessageTypeException:
			// todo
		default:
			// unknown cmd type
			return nil, fmt.Errorf("Decode Error, type = %s, value = %d", UnKnownCmdType, cmdType)
		}
	}

	return nil, nil
}

// Heartbeater
func (proto *ThriftProtocol) Trigger(requestId uint64) xprotocol.XFrame {
	request := &Request{
		RequestHeader: RequestHeader{
			Protocol: ProtocolCode,
			CmdType:  CmdTypeRequest,
			CmdCode:  CmdCodeHeartbeat,
			//Version:   1,
			RequestId: uint32(requestId),
			//Codec:     Hessian2Serialize,
			Timeout: -1,
		},
	}

	request.Data = buffer.GetIoBuffer(len(pingCmd))

	//5. copy data for io multiplexing
	request.Data.Write(pingCmd)
	request.rawData = request.Data.Bytes()
	request.rawContent = request.rawData[:]
	request.Content = buffer.NewIoBufferBytes(request.rawContent)

	return request
}

func (proto *ThriftProtocol) Reply(request xprotocol.XFrame) xprotocol.XRespFrame {
	response := &Response{
		ResponseHeader: ResponseHeader{
			Protocol: ProtocolCode,
			CmdType:  CmdTypeResponse,
			CmdCode:  CmdCodeHeartbeat,
			//Version:        ProtocolVersion,
			RequestId: uint32(request.GetRequestId()),
			//Codec:          Hessian2Serialize,
			ResponseStatus: ResponseStatusSuccess,
		},
	}

	response.Data = buffer.GetIoBuffer(len(pongCmd))
	//5. copy data for io multiplexing
	response.Data.Write(pongCmd)
	response.rawData = response.Data.Bytes()
	response.rawContent = response.rawData[:]
	response.Content = buffer.NewIoBufferBytes(response.rawContent)

	return response
}

// Hijacker
func (proto *ThriftProtocol) Hijack(request xprotocol.XFrame, statusCode uint32) xprotocol.XRespFrame {
	return &Response{
		ResponseHeader: ResponseHeader{
			Protocol: ProtocolCode,
			CmdType:  CmdTypeResponse,
			CmdCode:  CmdCodeResponse,
			//Version:        ProtocolVersion,
			RequestId:      0,                 // this would be overwrite by stream layer
			Codec:          Hessian2Serialize, //todo: read default codec from config
			ResponseStatus: uint16(statusCode),
		},
	}
}

func (proto *ThriftProtocol) Mapping(httpStatusCode uint32) uint32 {
	switch httpStatusCode {
	case http.StatusOK:
		return uint32(ResponseStatusSuccess)
	case types.RouterUnavailableCode:
		return uint32(ResponseStatusNoProcessor)
	case types.NoHealthUpstreamCode:
		return uint32(ResponseStatusConnectionClosed)
	case types.UpstreamOverFlowCode:
		return uint32(ResponseStatusServerThreadpoolBusy)
	case types.CodecExceptionCode:
		//Decode or Encode Error
		return uint32(ResponseStatusCodecException)
	case types.DeserialExceptionCode:
		//Hessian Exception
		return uint32(ResponseStatusServerDeserialException)
	case types.TimeoutExceptionCode:
		//Response Timeout
		return uint32(ResponseStatusTimeout)
	default:
		return uint32(ResponseStatusUnknown)
	}
}
