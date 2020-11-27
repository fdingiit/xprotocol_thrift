package filter

import (
	v1 "gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	"sync"

	"strings"

	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/metadata"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/types"
)

var (
	filterFactoryMap sync.Map
)

func Register(name string, factory types.FilterFactory) {
	filterFactoryMap.Store(name, factory)
}

func New(conf *v1.Filter) (types.GatewayFilter, error) {
	// Create Filter with FilterFactory
	if factory, ok := filterFactoryMap.Load(conf.Name); ok {
		// todo: now only support string
		if ref, ok := conf.Metadata.(string); ok {
			if strings.Index(ref, "$") == 0 {
				if val, node := metadata.Get(ref); node != nil {
					filterFactory := createConfigFilterFactory(factory.(types.FilterFactory), conf, ref)
					filter := filterFactory.Create(val)
					node.Register(filter, filterFactory.Create)
					return filter.(types.GatewayFilter), nil
				}
			}
		}

		return factory.(types.FilterFactory).CreateFilter(conf), nil
	}

	return nil, nil
}

type ConfigFilterFactory struct {
	factory types.FilterFactory
	conf    *v1.Filter
	ref     string
}

func (cff *ConfigFilterFactory) Create(conf interface{}) interface{} {
	cff.conf.Metadata = conf
	return cff.factory.CreateFilter(cff.conf)
}

func createConfigFilterFactory(factory types.FilterFactory, conf *v1.Filter, ref string) *ConfigFilterFactory {
	return &ConfigFilterFactory{
		factory: factory,
		conf:    conf,
		ref:     ref,
	}
}
