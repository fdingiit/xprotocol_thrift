package upstream

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/service/upstream/codec"
	"time"

	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/types"
)

func init() {
	GetUpstreamCreatorFactory().Register("SofaRpc", createSofaRpcUpstream)
}

type sofaRpcUpstream struct {
	*baseUpstream
	timeout   time.Duration
	targetUrl string
}

func createSofaRpcUpstream(conf *v1.GatewayService) types.Upstream {
	return &sofaRpcUpstream{
		&baseUpstream{
			protocol:    "SofaRpc",
			serviceName: conf.Name,
			codec:       codec.GetCodeFactoryInstance().GetCodec("SofaRpc"),
		},
		time.Duration(conf.Timeout),
		"todo",
	}
}
