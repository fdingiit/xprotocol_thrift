package sofadrm

import (
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"gitlab.alipay-inc.com/sofa-go/sofa-drm-client-go/sofadrm/model"

	"gitlab.alipay-inc.com/sofa-go/sofa-logger-go/sofalogger"
	"gitlab.alipay-inc.com/sofa-go/sofa-registry-client-go/sofaregistry"
)

const (
	DefaultDRMDataID = "com.alipay.zdrmdata.pub.server.url@DRM"
	DefaultDRMGroup  = "DRM"
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
		NewSubscribeContext(DefaultDRMDataID, DefaultDRMGroup).SetResident(true))
	// refresh server list
	// nolint
	rl.RefreshServers()

	select {
	case <-time.After(rl.config.registryLocator.timeout):
	case <-rl.pushedCh:

	}

	return rl, nil
}

func (rl *RegistryLocator) OnRegistryPush(dataID string, data map[string][]string, localZone string) {
	rl.logger.Infof("registry push dataID=%s data=%v localZone=%s", dataID, data, localZone)

	select {
	case rl.pushedCh <- struct{}{}:
	default:
	}

	// 100.88.64.59:9880?_SERIALIZETYPE=hessian2&_CONNECTTIMEOUT=5000&_TIMEOUT=10000
	value := data[localZone]
	if len(value) == 0 {
		return
	}

	servers := make([]model.Server, 0, len(value))
	for i := range value {
		target := value[i]
		if strings.Index(target, "://") == -1 {
			target = "tcp://" + target
		}
		u, err := url.ParseRequestURI(target)
		if err != nil {
			rl.logger.Errorf("failed to parse url with %q: %s", value[i], err.Error())
			continue
		}

		port := defaultDRMPort
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
	servers := rl.GetServers()
	if len(servers) == 0 {
		return model.Server{}, false
	}

	return servers[rand.Intn(len(servers))], true
}

func (rl *RegistryLocator) RefreshServers() error {
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
