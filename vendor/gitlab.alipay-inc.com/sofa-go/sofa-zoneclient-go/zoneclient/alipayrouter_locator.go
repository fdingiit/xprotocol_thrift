package zoneclient

import (
	"gitlab.alipay-inc.com/sofa-go/sofa-zoneclient-go/zoneclient/model"
)

type Locator interface {
	GetServers() (servers []model.Server)
	GetRandomServer() (server model.Server, ok bool)
	RefreshServers() error
}
