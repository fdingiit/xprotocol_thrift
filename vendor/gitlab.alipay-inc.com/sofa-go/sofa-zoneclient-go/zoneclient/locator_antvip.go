package zoneclient

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"gitlab.alipay-inc.com/sofa-go/sofa-zoneclient-go/zoneclient/model"

	"gitlab.alipay-inc.com/sofa-go/sofa-antvip-client-go/sofaantvip"
	"gitlab.alipay-inc.com/sofa-go/sofa-logger-go/sofalogger"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type AntVipLocator struct {
	sync.RWMutex
	logger  sofalogger.Logger
	config  *Config
	client  *sofaantvip.AntvipClient
	servers []model.Server
}

func NewAntVipLocator(options ...AntVipLocatortOptionSetter) (*AntVipLocator, error) {
	al := &AntVipLocator{}

	for _, option := range options {
		option.set(al)
	}

	if err := al.polyfill(); err != nil {
		return nil, err
	}

	al.client.ReserveDomain(al.config.alipayRouter.domain)
	//add antvip listener
	al.client.AddListener(al.config.alipayRouter.domain, al)

	if err := al.RefreshServers(); err != nil {
		return nil, err
	}

	return al, nil
}

func (al *AntVipLocator) RefreshServers() error {
	vd, err := al.client.GetVipDomain(
		al.config.GetAlipayRouterConfig().GetDomain(),
		al.config.GetAlipayRouterConfig().GetTimeout(),
	)
	if err != nil {
		return err
	}

	list := vd.GetAvailableRealAddress(defaultAlipayRouterAntvipPort)
	if len(list) > 0 {
		servers := make([]model.Server, 0, len(list))
		for _, target := range list {
			ip := target
			port := 80

			if strings.Index(strings.TrimSpace(target), ":") != -1 {
				array := strings.Split(strings.TrimSpace(target), ":")

				ip = array[0]
				up, err := strconv.ParseInt(array[1], 10, 16)
				if err == nil {
					port = int(up)
				} else {
					al.logger.Errorf("antviplocator: failed to parse port with %q using %d instead", up, port)
				}
			}

			servers = append(servers, model.Server{Ip: ip, Port: int32(port)})
		}
		al.setServers(servers)
	}

	return nil
}

func (al *AntVipLocator) OnDomainChanged(err error, name string, domain *sofaantvip.VipDomain) {
	al.logger.Infof("zoneclient: name = %v, domain = %v", name, domain)

	if err := al.RefreshServers(); err != nil {
		al.logger.Errorf("antviplocator: failed to refresh servers, %v", err.Error())
	}
}

func (al *AntVipLocator) GetServers() (servers []model.Server) {
	al.RLock()
	defer al.RUnlock()
	return al.servers
}

func (al *AntVipLocator) GetRandomServer() (server model.Server, ok bool) {
	al.logger.Infof("zoneclient: get server from antvip, %v", al.GetServers())

	servers := al.GetServers()
	if len(servers) == 0 {
		return model.Server{}, false
	}

	return servers[rand.Intn(len(servers))], true
}

func (al *AntVipLocator) setServers(servers []model.Server) {
	al.Lock()
	defer al.Unlock()
	al.servers = servers
}

func (al *AntVipLocator) polyfill() error {
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
