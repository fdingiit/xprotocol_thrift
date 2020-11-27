package sofaantvip

import (
	"fmt"
)

const (
	DefaultCloudSyncerServerAddress = "antvip-pool"
	DefaultCloudSyncerPort          = 9003
	DefaultCloudSyncerEndpoint      = "/antcloud/antvip/instances/get"
)

type CloudSyncerConfig struct {
	https     bool
	address   string
	port      int16
	endpoint  string
	accesslog bool
}

func (o *CloudSyncerConfig) EnableAccessLog() *CloudSyncerConfig {
	o.accesslog = true
	return o
}

func (o *CloudSyncerConfig) EnableHTTPS() *CloudSyncerConfig {
	o.https = true
	return o
}

func (o *CloudSyncerConfig) SetAddress(address string) *CloudSyncerConfig {
	o.address = address
	return o
}

func (o *CloudSyncerConfig) SetPort(port int16) *CloudSyncerConfig {
	o.port = port
	return o
}

func (o *CloudSyncerConfig) SetEndpoint(endpoint string) *CloudSyncerConfig {
	o.endpoint = endpoint
	return o
}

func (o *CloudSyncerConfig) GetURL() string {
	scheme := "http"
	if o.https {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s:%d%s", scheme, o.address, o.port,
		o.endpoint)
}
