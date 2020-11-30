package sofaantvip

import (
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"

	"gitlab.alipay-inc.com/sofa-go/sofa-logger-go/sofalogger"
)

// HTTPLocatorOptionSetter configures a HTTPLocator.
type HTTPLocatorOptionSetter interface {
	set(*HTTPLocator)
}

type HTTPLocatorOptionSetterFunc func(*HTTPLocator)

func (f HTTPLocatorOptionSetterFunc) set(c *HTTPLocator) {
	f(c)
}

func WithHTTPLocatorConfig(config *Config) HTTPLocatorOptionSetterFunc {
	return HTTPLocatorOptionSetterFunc(func(c *HTTPLocator) {
		c.config = config
	})
}

func WithHTTPLocatorLogger(logger sofalogger.Logger) HTTPLocatorOptionSetterFunc {
	return HTTPLocatorOptionSetterFunc(func(c *HTTPLocator) {
		c.logger = logger
	})
}

type HTTPLocator struct {
	sync.RWMutex
	config  *Config
	servers vipServers
	logger  sofalogger.Logger
}

func NewHTTPLocator(options ...HTTPLocatorOptionSetterFunc) (*HTTPLocator, error) {
	hl := &HTTPLocator{}
	for i := range options {
		options[i].set(hl)
	}

	if err := hl.polyfill(); err != nil {
		return nil, err
	}

	// Try to load the servers even if it's error
	_ = hl.LoadAndStoreVipServers()

	if hl.config.syncInterval > 0 {
		go hl.doIntervalSync()
	}

	return hl, nil
}

func (hl *HTTPLocator) polyfill() error {
	if hl.config == nil {
		return errors.New("httplocator: config cannot be nil")
	}

	if hl.logger == nil {
		hl.logger = sofalogger.StdoutLogger
	}

	return nil
}

func (hl *HTTPLocator) doIntervalSync() {
	for {
		time.Sleep(hl.config.httpLocator.interval)
		if err := hl.RefreshServers(); err != nil {
			hl.logger.Errorf("failed to refresh servers: %s", err.Error())
		}
	}
}

func (hl *HTTPLocator) RefreshServers() error {
	return hl.LoadAndStoreVipServers()
}

func (hl *HTTPLocator) GetRandomServer() (VipServer, bool) {
	hl.RLock()
	defer hl.RUnlock()
	return hl.servers.getRandomServer()
}

func (hl *HTTPLocator) GetServers() []VipServer {
	hl.RLock()
	defer hl.RUnlock()
	return hl.servers.getServers()
}

func (hl *HTTPLocator) GetChecksum() string {
	hl.RLock()
	defer hl.RUnlock()
	return hl.servers.getChecksum()
}

func (hl *HTTPLocator) LoadAndStoreVipServers() error {
	return hl.LoadAndStoreVipServersFromURL(hl.config.GetHTTPLocatorURL())
}

func (hl *HTTPLocator) LoadAndStoreVipServersFromURL(url string) error {
	csv, err := hl.LoadVipServersFromURL(url)
	if err != nil {
		return err
	}

	servers, err := newVipServersFromCSV(csv)
	if err != nil {
		return err
	}

	if len(servers.servers) > 0 {
		hl.Lock()
		hl.servers = *servers
		hl.Unlock()
	}

	return nil
}

func (hl *HTTPLocator) LoadVipServers() (string, error) {
	return hl.LoadVipServersFromURL(hl.config.GetHTTPLocatorURL())
}

func (hl *HTTPLocator) LoadVipServersFromURL(url string) (string, error) {
	var (
		res        *http.Response
		err        error
		data       []byte
		statuscode int
	)

	if hl.config.httpLocator.accesslog {
		started := time.Now()
		defer func() {
			hl.logger.Infof("HTTPLocator %s %s %d <%s> %d %d <%s> %s", "GET", url,
				statuscode, string(data), 0, len(data), errstring(err), time.Since(started).String())
		}()
	}

	res, err = hl.getHTTPClient().Get(url)
	if err != nil {
		return "", err
	}
	// nolint
	defer res.Body.Close()
	statuscode = res.StatusCode

	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (hl *HTTPLocator) getHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   hl.config.httpLocator.timeout,
				KeepAlive: 0,
				DualStack: false,
			}).DialContext,
			DisableKeepAlives:   true, // short-lived connection
			MaxIdleConns:        0,
			TLSHandshakeTimeout: hl.config.httpLocator.timeout,
		},
		Timeout: hl.config.httpLocator.timeout,
	}
}
