package sofaantvip

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"gitlab.alipay-inc.com/sofa-go/sofa-logger-go/sofalogger"
	"gitlab.alipay-inc.com/sofa-go/sofa-registry-client-go/sofaregistry"
)

const (
	DefaultAntvipServerListDataID = "com.alipay.antvip.serverlist.client"
)

var _ Locator = (*RegistryLocator)(nil)

type RegistryLocator struct {
	logger      sofalogger.Logger
	client      *sofaregistry.Client
	config      *Config
	serversLock sync.RWMutex
	servers     vipServers
	subscriber  *sofaregistry.Subscriber
	pushedCh    chan struct{}
}

func NewRegistryLocator(options ...RegistryLocatortOptionSetter) (*RegistryLocator, error) {
	rl := &RegistryLocator{}

	for _, option := range options {
		option.set(rl)
	}

	if err := rl.polyfill(); err != nil {
		return nil, err
	}

	rl.subscriber = rl.client.CreateSubscriber(sofaregistry.
		NewSubscribeContext(DefaultAntvipServerListDataID, sofaregistry.DefaultGroup).
		SetResident(true).
		SetScope(sofaregistry.ScopeDataCenter),
	)
	// refresh server list
	// nolint
	rl.RefreshServers()

	select {
	case <-time.After(rl.config.registryLocator.timeout):
	case <-rl.pushedCh:

	}

	return rl, nil
}

func (rl *RegistryLocator) buildVIPServerFromRegistryData(value []string) []VipServer {
	// nolint
	// map[CZ00B:[11.167.218.68:12200,5 11.166.233.181:12200,5 11.167.194.115:12200,5 11.166.233.182:12200,5 11.167.175.230:12200,5]]
	vs := make([]VipServer, 0, 8)
	for i := range value {
		var (
			addr   string
			weight int64
		)

		x := strings.Split(value[i], ",")
		switch len(x) {
		case 0:
			continue
		case 1:
			addr = x[0]
		default:
			addr = x[0]
			weight, _ = strconv.ParseInt(x[1], 10, 64)
		}

		host, _, err := net.SplitHostPort(addr)
		if err != nil {
			rl.logger.Errorf("failed to split hostport: %q", value[i])
			continue
		}

		vs = append(vs, VipServer{
			Host:     host,
			HostName: host,
			Weight:   int32(weight),
		})
	}

	return vs
}

func (rl *RegistryLocator) OnRegistryPush(dataID string, data map[string][]string, localZone string) {
	rl.logger.Infof("registry push dataID=%s data=%v localZone=%s", dataID, data, localZone)

	select {
	case rl.pushedCh <- struct{}{}:
	default:
	}

	servers := make([]VipServer, 0, 16)
	for zone := range data {
		value := data[zone]
		s := rl.buildVIPServerFromRegistryData(value)
		servers = append(servers, s...)
	}

	rl.setServers(servers)
}

func (rl *RegistryLocator) setServers(servers []VipServer) {
	vs := newVipServers(servers)
	rl.serversLock.Lock()
	rl.servers = vs
	rl.serversLock.Unlock()
}

func (rl *RegistryLocator) GetServers() (servers []VipServer) {
	rl.serversLock.RLock()
	defer rl.serversLock.RUnlock()
	return rl.servers.getServers()
}

func (rl *RegistryLocator) GetChecksum() string {
	rl.serversLock.RLock()
	defer rl.serversLock.RUnlock()
	return rl.servers.getChecksum()
}

func (rl *RegistryLocator) GetRandomServer() (server VipServer, ok bool) {
	rl.serversLock.RLock()
	defer rl.serversLock.RUnlock()
	return rl.servers.getRandomServer()
}

func (rl *RegistryLocator) RefreshServers() error {
	rl.serversLock.Lock()
	defer rl.serversLock.Unlock()
	// resubscribe to the registry
	return rl.subscriber.Sub(rl)
}

func (rl *RegistryLocator) polyfill() error {
	rl.pushedCh = make(chan struct{})

	if rl.config == nil {
		return fmt.Errorf("registrylocator: config is nil")
	}

	if rl.client == nil {
		return fmt.Errorf("registrylocator: client is nil")
	}

	if rl.logger == nil {
		rl.logger = sofalogger.StdoutLogger
	}

	return nil
}
