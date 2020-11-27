package codec

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/common/constants"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/types"
	mosn "mosn.io/mosn/pkg/types"
)

type http1Codec struct {
	protocol string
}

func (src *http1Codec) Convertor() types.ProtocolConvertor {
	panic("implement me")
}

func init() {
	GetCodeFactoryInstance().Register(&http1Codec{protocol: "Http1"})
}

func (src *http1Codec) Protocol() string {
	return src.protocol
}

func (src *http1Codec) Encode(ctx types.Context, headers mosn.HeaderMap, buf mosn.IoBuffer, trailers mosn.HeaderMap) {
	headers.Set(constants.MosngServiceHeader, ctx.Router().Service().Name())
}

func (src *http1Codec) Decode(ctx types.Context, headers mosn.HeaderMap, buf mosn.IoBuffer, trailers mosn.HeaderMap) {
	// noop now
}
