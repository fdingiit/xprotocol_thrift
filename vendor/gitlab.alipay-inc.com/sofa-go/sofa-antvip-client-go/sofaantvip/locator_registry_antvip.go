package sofaantvip

var _ Locator = (*RegistryAntvipLocator)(nil)

type RegistryAntvipLocator struct {
	registry *RegistryLocator
	http     *HTTPLocator
}

func NewRegistryAntvipLocator(r *RegistryLocator, a *HTTPLocator) *RegistryAntvipLocator {
	return &RegistryAntvipLocator{
		registry: r,
		http:     a,
	}
}

func (r *RegistryAntvipLocator) GetServers() []VipServer {
	servers := r.registry.GetServers()
	if len(servers) == 0 {
		return r.http.GetServers()
	}
	return servers
}

func (r *RegistryAntvipLocator) GetRandomServer() (VipServer, bool) {
	server, ok := r.registry.GetRandomServer()
	if ok {
		return server, true
	}

	return r.http.GetRandomServer()
}

func (r *RegistryAntvipLocator) RefreshServers() error {
	rerr := r.registry.RefreshServers()
	aerr := r.http.RefreshServers()
	return merror(rerr, aerr)
}

func (r *RegistryAntvipLocator) GetChecksum() string {
	return r.registry.GetChecksum()
}
