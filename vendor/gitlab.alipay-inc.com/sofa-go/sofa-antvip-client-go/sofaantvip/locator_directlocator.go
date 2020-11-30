package sofaantvip

type DirectLocator struct {
	servers vipServers
}

func NewDirectLocator(servers ...VipServer) *DirectLocator {
	dl := &DirectLocator{
		servers: newVipServers(servers),
	}
	return dl
}

func (dl *DirectLocator) GetServers() (servers []VipServer) {
	return dl.servers.getServers()
}

func (dl *DirectLocator) GetRandomServer() (server VipServer, ok bool) {
	return dl.servers.getRandomServer()
}

func (dl *DirectLocator) GetChecksum() string {
	return dl.servers.getChecksum()
}

func (dl *DirectLocator) RefreshServers() error {
	return nil
}
