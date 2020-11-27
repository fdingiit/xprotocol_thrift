package sofaantvip

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"

	"gitlab.alipay-inc.com/sofa-go/sofa-logger-go/sofalogger"
)

// type alias for sofabolt.Dialer
type Dialer interface {
	Dial() (net.Conn, error)
}

var _ Dialer = (*AntvipDialer)(nil)

type AntvipDialerOptionSetter interface {
	set(*AntvipDialer)
}

type AntvipDialerSetterFunc func(*AntvipDialer)

func (f AntvipDialerSetterFunc) set(c *AntvipDialer) {
	f(c)
}

func WithAntvipDialerClient(client *AntvipClient) AntvipDialerSetterFunc {
	return AntvipDialerSetterFunc(func(c *AntvipDialer) {
		c.client = client
	})
}

func WithAntvipDialerDomain(domain string) AntvipDialerSetterFunc {
	return AntvipDialerSetterFunc(func(c *AntvipDialer) {
		c.domain = domain
	})
}

func WithAntvipDialerTimeout(timeout time.Duration) AntvipDialerSetterFunc {
	return AntvipDialerSetterFunc(func(c *AntvipDialer) {
		c.timeout = timeout
	})
}

func WithAntvipDialerLogger(logger sofalogger.Logger) AntvipDialerSetterFunc {
	return AntvipDialerSetterFunc(func(c *AntvipDialer) {
		c.logger = logger
	})
}

type AntvipDialer struct {
	serversLock sync.RWMutex
	servers     []RealServer
	domain      string
	timeout     time.Duration
	client      *AntvipClient
	logger      sofalogger.Logger
}

func NewAntvipDialer(options ...AntvipDialerOptionSetter) (*AntvipDialer, error) {
	ad := &AntvipDialer{}

	for _, option := range options {
		option.set(ad)
	}

	if err := ad.polyfill(); err != nil {
		return nil, err
	}

	ad.doLookupAndListen()

	return ad, nil
}

func (ad *AntvipDialer) polyfill() error {
	if ad.logger == nil {
		ad.logger = sofalogger.StdoutLogger
	}

	if ad.client == nil {
		return fmt.Errorf("antvipdialer: client cannot be nil")
	}

	if ad.domain == "" {
		return fmt.Errorf("antvipdialer: domain cannoe be nil")
	}

	if ad.timeout == 0 {
		ad.timeout = 3 * time.Second
	}

	return nil
}

func (ad *AntvipDialer) OnDomainChanged(err error, name string, domain *VipDomain) {
	if err != nil {
		ad.logger.Errorf("antvip push name=%q with error: %+v", name, err)
		return
	}

	ad.logger.Infof("antvip push name=%q servers=%+v version=%s deleted=%t",
		name, domain.GetRealServers(), domain.Version, domain.IsDeleted)

	ad.setServers(ad.filterServers(domain.GetRealServers()))
}

func (ad *AntvipDialer) getServers() (servers []RealServer) {
	ad.serversLock.RLock()
	servers = ad.servers
	ad.serversLock.RUnlock()
	return servers
}

func (ad *AntvipDialer) filterServers(realServers []RealServer) []RealServer {
	servers := make([]RealServer, 0, len(realServers))
	for i := range realServers {
		s := realServers[i]
		if s.IsAvailable() && s.weight > 0 {
			servers = append(servers, s)
		}
	}
	return servers
}

func (ad *AntvipDialer) setServers(servers []RealServer) {
	ad.serversLock.Lock()
	ad.servers = servers
	ad.serversLock.Unlock()
}

func (ad *AntvipDialer) doLookupAndListen() {
	ad.client.AddListener(ad.domain, ad)

	vd, err := ad.client.GetVipDomain(ad.domain, ad.timeout)
	if err != nil {
		return
	}

	ad.setServers(ad.filterServers(vd.GetRealServers()))
}

func (ad *AntvipDialer) Dial() (net.Conn, error) {
	servers := ad.getServers()
	if len(servers) == 0 {
		ad.logger.Errorf("redial choose nil: no available server")
		return nil, errors.New("antvip: no available server")
	}

	s := servers[rand.Intn(len(servers))]
	ip := s.ip
	port := s.healthCheckPort
	if port == 0 {
		port = 9600
	}

	addr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", addr, ad.timeout)
	if err != nil {
		ad.logger.Errorf("redial choose %s %s", addr, err.Error())
	} else {
		ad.logger.Infof("redial choose %s", addr)
	}

	return conn, err
}
