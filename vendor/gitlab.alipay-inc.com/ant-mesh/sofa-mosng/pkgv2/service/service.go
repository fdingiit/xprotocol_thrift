package service

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/common/log"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api"
	v1 "gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/event"
	"sync"

	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/service/upstream"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/types"
)

func init() {
	event.EventListenerManagerInstance().Register(event.ResourceEventListenerFuncs{
		Type:       api.SERVICE,
		AddFunc:    GetServiceManagerInstance().AddOrUpdateService,
		UpdateFunc: GetServiceManagerInstance().AddOrUpdateService,
	})
}

type service struct {
	name     string
	conf     *v1.GatewayService
	upstream types.Upstream
}

func (si *service) Name() string {
	return si.name
}

func (si *service) Upstream() types.Upstream {
	return si.upstream
}

func (si *service) Conf() *v1.GatewayService {
	return si.conf
}

type manager struct {
	mux sync.Mutex
	// name : service
	services map[string]types.Service
}

func (m *manager) AddOrUpdateService(o api.Object) (error, bool) {
	conf := o.(*v1.GatewayService)

	log.ConfigLogger().Infof("[gateway][service][AddOrUpdateService] start add service %s", conf.Name)

	svc := newService(conf)

	m.mux.Lock()
	defer m.mux.Unlock()
	m.services[conf.Name] = svc

	return nil, true
}

func (m *manager) GetService(serviceName, protocol string) types.Service {
	service := m.services[serviceName]
	// todo protocol
	return service
}

var (
	serviceManagerInstance *manager
	once                   sync.Once
)

func GetServiceManagerInstance() *manager {
	once.Do(func() {
		serviceManagerInstance = buildServiceManager()
	})
	return serviceManagerInstance
}

func buildServiceManager() *manager {
	return &manager{
		mux:      sync.Mutex{},
		services: make(map[string]types.Service),
	}
}

func newService(conf *v1.GatewayService) types.Service {
	return &service{
		name:     conf.Name,
		conf:     conf,
		upstream: upstream.GetUpstreamCreatorFactory().NewUpstream(conf),
	}
}
