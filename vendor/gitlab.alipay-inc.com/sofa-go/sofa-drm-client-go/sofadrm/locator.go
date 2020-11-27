package sofadrm

import "gitlab.alipay-inc.com/sofa-go/sofa-drm-client-go/sofadrm/model"

type Locator interface {
	GetServers() (servers []model.Server)
	GetRandomServer() (server model.Server, ok bool)
	RefreshServers() error
}
