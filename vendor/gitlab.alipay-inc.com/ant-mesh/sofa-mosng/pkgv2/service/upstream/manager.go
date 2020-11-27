package upstream

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	"sync"

	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/types"
)

type upstreamCreatorFactory struct {
	creators map[string]upstreamCreator
}

type upstreamCreator = func(conf *v1.GatewayService) types.Upstream

func (ucf *upstreamCreatorFactory) Register(protocol string, creator upstreamCreator) {
	ucf.creators[protocol] = creator
}

func (ucf *upstreamCreatorFactory) NewUpstream(conf *v1.GatewayService) types.Upstream {
	if creator := ucf.creators[conf.Protocol]; creator != nil {
		return creator(conf)
	}
	// todo err
	return nil
}

var (
	once                           sync.Once
	upstreamCreatorFactoryInstance *upstreamCreatorFactory
)

func GetUpstreamCreatorFactory() *upstreamCreatorFactory {
	once.Do(func() {
		upstreamCreatorFactoryInstance = newUpstreamCreatorFactory()
	})

	return upstreamCreatorFactoryInstance
}

func newUpstreamCreatorFactory() *upstreamCreatorFactory {
	return &upstreamCreatorFactory{
		creators: make(map[string]upstreamCreator),
	}
}
