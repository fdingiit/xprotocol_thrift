package types

import (
	"mosn.io/mosn/pkg/types"
)

type Codec interface {
	Protocol() string
	Encode(ctx Context, headers types.HeaderMap, buf types.IoBuffer, trailers types.HeaderMap)
	Decode(ctx Context, headers types.HeaderMap, buf types.IoBuffer, trailers types.HeaderMap)
	Convertor() ProtocolConvertor
}

type ProtocolConvertor interface {
	EncodeConv(ctx Context, headers types.HeaderMap, buf types.IoBuffer, trailers types.HeaderMap) (convHeaders types.HeaderMap, convBuf types.IoBuffer, convTrailers types.HeaderMap)
	DecodeConv(ctx Context, headers types.HeaderMap, buf types.IoBuffer, trailers types.HeaderMap) (convHeaders types.HeaderMap, convBuf types.IoBuffer, convTrailers types.HeaderMap)
}

type CodecFactory interface {
	Register(codec Codec)
	GetCodec(protocol string) Codec
}

type Upstream interface {
	Invoke(ctx Context) error
	Receive(ctx Context) error
	Codec() Codec
}
