package zoneclient

import (
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"gitlab.alipay-inc.com/sofa-go/sofa-zoneclient-go/zoneclient/model"

	"gitlab.alipay-inc.com/sofa-go/sofa-logger-go/sofalogger"
	"gitlab.alipay-inc.com/sofa-go/sofa-registry-client-go/sofaregistry"
)

const (
	DefaultDataID                = "com.alipay.zoneclient.pub.server.url"
	DefaultGroup                 = "DEFAULT_GROUP"
	DefaultRegistyLocatorTimeout = 3 * time.Second
)

var _ Locator = (*RegistryLocator)(nil)

type RegistryLocator struct {
	logger      sofalogger.Logger
	client      *sofaregistry.Client
	config      *Config
	serversLock sync.RWMutex
	servers     []model.Server
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
		NewSubscribeContext(DefaultDataID, DefaultGroup).SetScope(sofaregistry.ScopeDataCenter).SetResident(true))
	if err := rl.RefreshServers(); err != nil {
		return nil, err
	}

	if _, _, err := rl.subscriber.Peek(DefaultRegistyLocatorTimeout); err != nil {
		rl.logger.Errorf("registry peek failed,%v", err.Error())
	}

	return rl, nil
}

func (rl *RegistryLocator) OnRegistryPush(dataID string, data map[string][]string, localZone string) {
	rl.logger.Infof("registry push dataID=%s data=%v localZone=%s", dataID, data, localZone)

	servers := make([]model.Server, 0)
	if data != nil {
		for _, rsList := range data {
			if len(rsList) > 0 {
				for i := range rsList {
					target := rsList[i]
					if strings.Index(target, "://") == -1 {
						target = "http://" + target
					}
					u, err := url.ParseRequestURI(target)
					if err != nil {
						rl.logger.Errorf("failed to parse url with %q: %s", rsList[i], err.Error())
						continue
					}

					port := 80
					up, err := strconv.ParseUint(u.Port(), 10, 16)
					if err == nil {
						port = int(up)
					} else {
						rl.logger.Errorf("failed to parse port with %q using %d instead", u.Port(), port)
					}

					servers = append(servers, model.Server{
						Ip:   u.Hostname(),
						Port: int32(port),
					})
				}
			}
		}
	}

	rl.setServers(servers)
}

func (rl *RegistryLocator) setServers(servers []model.Server) {
	rl.serversLock.Lock()
	rl.servers = servers
	rl.serversLock.Unlock()
}

func (rl *RegistryLocator) GetServers() (servers []model.Server) {
	rl.serversLock.RLock()
	servers = rl.servers
	rl.serversLock.RUnlock()
	return servers
}

func (rl *RegistryLocator) GetRandomServer() (server model.Server, ok bool) {
	rl.logger.Infof("zoneclient: get server from registry, %v", rl.GetServers())

	servers := rl.GetServers()
	if len(servers) == 0 {
		return model.Server{}, false
	}

	return servers[rand.Intn(len(servers))], true
}

func (rl *RegistryLocator) RefreshServers() error {
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
