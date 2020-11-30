package sofadrm

import (
	"gitlab.alipay-inc.com/sofa-go/sofa-drm-client-go/sofadrm/model"
)

var _ Locator = (*RegistryAntvipLocator)(nil)

type RegistryAntvipLocator struct {
	registry *RegistryLocator
	antvip   *AntvipLocator
}

func NewRegistryAntvipLocator(r *RegistryLocator, a *AntvipLocator) *RegistryAntvipLocator {
	return &RegistryAntvipLocator{
		registry: r,
		antvip:   a,
	}
}

func (r *RegistryAntvipLocator) GetServers() []model.Server {
	servers := r.registry.GetServers()
	if len(servers) == 0 {
		return r.antvip.GetServers()
	}
	return servers
}

func (r *RegistryAntvipLocator) GetRandomServer() (model.Server, bool) {
	server, ok := r.registry.GetRandomServer()
	if ok {
		return server, true
	}

	return r.antvip.GetRandomServer()
}

func (r *RegistryAntvipLocator) RefreshServers() error {
	rerr := r.registry.RefreshServers()
	aerr := r.antvip.RefreshServers()
	return merror(rerr, aerr)
}
