package sofadrm

import (
	"math/rand"

	"gitlab.alipay-inc.com/sofa-go/sofa-drm-client-go/sofadrm/model"
)

type DirectLocator struct {
	servers []model.Server
}

func NewDirectLocator(servers []model.Server) *DirectLocator {
	return &DirectLocator{
		servers: servers,
	}
}

func (dl *DirectLocator) GetServers() (servers []model.Server) {
	return dl.servers
}

func (dl *DirectLocator) GetRandomServer() (server model.Server, ok bool) {
	servers := dl.GetServers()
	if len(servers) == 0 {
		return model.Server{}, false
	}

	return servers[rand.Intn(len(servers))], true
}

func (dl *DirectLocator) RefreshServers() error {
	return nil
}
