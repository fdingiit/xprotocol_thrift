package sofadrm

import (
	"fmt"
	"math/rand"
	"sync"

	"gitlab.alipay-inc.com/sofa-go/sofa-drm-client-go/sofadrm/model"

	"gitlab.alipay-inc.com/sofa-go/sofa-antvip-client-go/sofaantvip"
	"gitlab.alipay-inc.com/sofa-go/sofa-logger-go/sofalogger"
)

type AntvipLocator struct {
	logger      sofalogger.Logger
	client      *sofaantvip.AntvipClient
	config      *Config
	serversLock sync.RWMutex
	servers     []model.Server
}

func NewAntvipLocator(options ...AntvipLocatortOptionSetter) (*AntvipLocator, error) {
	al := &AntvipLocator{}

	for _, option := range options {
		option.set(al)
	}

	if err := al.polyfill(); err != nil {
		return nil, err
	}

	// refresh server list
	// nolint
	al.RefreshServers()

	return al, nil
}

func (al *AntvipLocator) polyfill() error {
	if al.config == nil {
		return fmt.Errorf("antviplocator: config is nil")
	}

	if al.client == nil {
		return fmt.Errorf("antviplocator: client is nil")
	}

	if al.logger == nil {
		al.logger = sofalogger.StdoutLogger
	}

	return nil
}

func (al *AntvipLocator) OnDomainChanged(err error, name string, domain *sofaantvip.VipDomain) {
	if err != nil {
		al.logger.Errorf("failed to get servers from antvip: %v", err)
		return
	}

	if domain := al.config.GetAntvipLocatorConfig().GetDomain(); domain != name {
		al.logger.Errorf("expect domain %s but got %s", domain, name)
		return
	}

	al.logger.Infof("get servers from antvip: %v", domain.GetRealServers())

	finalServers := model.NewAvailableServers(domain.GetRealServers(), domain.ProtectThreshold, defaultDRMPort)
	if len(finalServers) > 0 {
		al.setServers(finalServers)
	}
}

func (al *AntvipLocator) setServers(servers []model.Server) {
	al.serversLock.Lock()
	al.servers = servers
	al.serversLock.Unlock()
}

func (al *AntvipLocator) GetServers() (servers []model.Server) {
	al.serversLock.RLock()
	servers = al.servers
	al.serversLock.RUnlock()
	return servers
}

func (al *AntvipLocator) GetRandomServer() (server model.Server, ok bool) {
	servers := al.GetServers()
	if len(servers) == 0 {
		return model.Server{}, false
	}

	return servers[rand.Intn(len(servers))], true
}

func (al *AntvipLocator) RefreshServers() error {
	domain := al.config.GetAntvipLocatorConfig().GetDomain()
	al.client.ReserveDomain(domain)
	al.client.AddListener(domain, al)
	vd, err := al.client.GetVipDomain(domain, al.config.GetAntvipLocatorConfig().GetTimeout())
	al.OnDomainChanged(err, domain, vd)
	return err
}
