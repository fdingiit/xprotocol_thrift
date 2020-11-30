package resource

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config"
)

type baseResource struct {
	codec      config.Codec
	dispatcher config.Dispatcher
}

func (b *baseResource) SetCodec(c config.Codec) {
	b.codec = c
}

func (b *baseResource) SetDispatcher(d config.Dispatcher) {
	b.dispatcher = d
}
