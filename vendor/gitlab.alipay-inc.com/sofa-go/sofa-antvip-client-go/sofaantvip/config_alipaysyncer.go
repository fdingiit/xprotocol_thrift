package sofaantvip

import (
	"fmt"
)

const (
	DefaultAlipaySyncerPort     = 9500
	DefaultAlipaySyncerEndpoint = "/antvipDomains"
)

type AlipaySyncerConfig struct {
	https     bool
	port      int16
	endpoint  string
	accesslog bool
	zi        *ZoneInfoList
}

func (o *AlipaySyncerConfig) GetZoneInfoList() *ZoneInfoList { return o.zi }

func (o *AlipaySyncerConfig) SetZoneInfoList(zi *ZoneInfoList) *AlipaySyncerConfig {
	o.zi = zi
	return o
}

func (o *AlipaySyncerConfig) SetEndpoint(endpoint string) *AlipaySyncerConfig {
	o.endpoint = endpoint
	return o
}

func (o *AlipaySyncerConfig) SetPort(p int16) *AlipaySyncerConfig {
	o.port = p
	return o
}

func (o *AlipaySyncerConfig) EnableAccessLog() *AlipaySyncerConfig {
	o.accesslog = true
	return o
}

func (o *AlipaySyncerConfig) GetURL(host string) string {
	scheme := "http"
	if o.https {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s:%d%s", scheme, host, o.port,
		o.endpoint)
}
