package sofaantvip

type Locator interface {
	GetServers() (servers []VipServer)
	GetRandomServer() (server VipServer, ok bool)
	GetChecksum() string
	RefreshServers() error
}
