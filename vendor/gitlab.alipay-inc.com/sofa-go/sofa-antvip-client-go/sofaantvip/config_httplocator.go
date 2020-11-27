package sofaantvip

import (
	"time"
)

const (
	DefaultHTTPLocatorAntVipServerAddress = "antvip-pool"
	DefaultHTTPLocatorTimeout             = 5 * time.Second
	DefaultHTTPLocatorInterval            = 30 * time.Second
	DefaultHTTPLocatorAntVipHTTPPort      = 9500
	DefaultHTTPLocatorAntVipEndpoint      = "/antvip/serversByZone.do"
)

type HTTPLocatorConfig struct {
	https     bool
	address   string
	port      int16
	endpoint  string
	timeout   time.Duration
	interval  time.Duration
	accesslog bool
}

func (o *HTTPLocatorConfig) EnableHTTPS() *HTTPLocatorConfig {
	o.https = true
	return o
}

func (o *HTTPLocatorConfig) SetAddress(address string) *HTTPLocatorConfig {
	o.address = address
	return o
}

func (o *HTTPLocatorConfig) SetPort(port int16) *HTTPLocatorConfig {
	o.port = port
	return o
}

func (o *HTTPLocatorConfig) SetEndpoint(endpoint string) *HTTPLocatorConfig {
	o.endpoint = endpoint
	return o
}

func (o *HTTPLocatorConfig) SetInterval(interval time.Duration) *HTTPLocatorConfig {
	o.interval = interval
	return o
}

func (o *HTTPLocatorConfig) EnableAccessLog() *HTTPLocatorConfig {
	o.accesslog = true
	return o
}
