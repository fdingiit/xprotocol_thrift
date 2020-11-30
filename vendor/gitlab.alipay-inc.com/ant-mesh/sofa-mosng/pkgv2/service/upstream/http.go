package upstream

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/service/upstream/codec"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/types"
)

func init() {
	GetUpstreamCreatorFactory().Register("Http1", createHttpUpstream)
}

func createHttpUpstream(conf *v1.GatewayService) types.Upstream {
	return &baseUpstream{
		protocol:    "Http1",
		serviceName: conf.Name,
		codec:       codec.GetCodeFactoryInstance().GetCodec("Http1"),
	}
}

type baseUpstream struct {
	protocol    string
	serviceName string
	codec       types.Codec
}

func (bu *baseUpstream) Invoke(ctx types.Context) error {

	// protocol convert
	if ctx.DownStreamProtocol() != ctx.UpstreamProtocol() {
		headers, buf, trailers := bu.Codec().Convertor().EncodeConv(ctx, ctx.Request().GetHeaders(), ctx.Request().GetDataBuf(), ctx.Request().GetTrailers())
		ctx.Request().SetHeaders(headers)
		ctx.Request().SetDataBuf(buf)
		ctx.Request().SetTrailers(trailers)
	}

	// default invoke codec
	bu.Codec().Encode(ctx, ctx.Request().GetHeaders(), ctx.Request().GetDataBuf(), ctx.Request().GetTrailers())

	return nil
}

func (bu *baseUpstream) Receive(ctx types.Context) error {
	// default invoke codec
	bu.Codec().Decode(ctx, ctx.Response().GetHeaders(), ctx.Response().GetDataBuf(), ctx.Response().GetTrailers())

	if ctx.DownStreamProtocol() != ctx.UpstreamProtocol() {
		headers, buf, trailers := bu.Codec().Convertor().DecodeConv(ctx, ctx.Response().GetHeaders(), ctx.Response().GetDataBuf(), ctx.Response().GetTrailers())
		ctx.Response().SetHeaders(headers)
		ctx.Response().SetDataBuf(buf)
		ctx.Response().SetTrailers(trailers)
	}

	return nil
}

func (bu *baseUpstream) Codec() types.Codec {
	return bu.codec
}
