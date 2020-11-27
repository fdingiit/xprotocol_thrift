package zoneclient

import (
	"gitlab.alipay-inc.com/sofa-go/sofa-zoneclient-go/zoneclient/model"
)

var _ Locator = (*RegistryAntVipLocator)(nil)

type RegistryAntVipLocator struct {
	registry     *RegistryLocator
	antvip       *AntVipLocator
	skipRegistry bool
}

func NewRegistryAntVipLocator(r *RegistryLocator, a *AntVipLocator) *RegistryAntVipLocator {
	return &RegistryAntVipLocator{
		registry: r,
		antvip:   a,
	}
}

func (r *RegistryAntVipLocator) SetSkipRegistry(skipRegistry bool) {
	r.skipRegistry = skipRegistry
}

func (r *RegistryAntVipLocator) GetServers() []model.Server {
	var servers []model.Server
	if !r.skipRegistry {
		servers = r.registry.GetServers()
	}

	if len(servers) == 0 {
		return r.antvip.GetServers()
	}
	return servers
}

func (r *RegistryAntVipLocator) GetRandomServer() (model.Server, bool) {
	if !r.skipRegistry {
		server, ok := r.registry.GetRandomServer()
		if ok {
			return server, true
		}
	}

	return r.antvip.GetRandomServer()
}

func (r *RegistryAntVipLocator) RefreshServers() error {
	var rerr error
	if !r.skipRegistry {
		rerr = r.registry.RefreshServers()
	}
	aerr := r.antvip.RefreshServers()
	return merror(rerr, aerr)
}
