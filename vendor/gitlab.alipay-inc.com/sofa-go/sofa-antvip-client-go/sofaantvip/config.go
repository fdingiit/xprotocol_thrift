package sofaantvip

import (
	"fmt"
	"time"
)

const (
	DefaultSyncerTimeout  = 90 * time.Second
	DefaultSyncerInterval = 1 * time.Second
)

type DomainChecksumMode uint8

var (
	DomainAlipayChecksumMode DomainChecksumMode = 0
	DomainCloudChecksumMode  DomainChecksumMode = 1
)

type Config struct {
	metrics            *Metrics
	appName            string
	datacenter         string
	zone               string
	instanceID         string
	endpoint           string
	accessKey          string
	accessSecret       string
	trFrom             string
	version            string
	domainChecksumMode DomainChecksumMode
	syncInterval       time.Duration
	syncTimeout        time.Duration
	httpLocator        HTTPLocatorConfig
	registryLocator    RegistryLocatorConfig
	alipaySyncer       AlipaySyncerConfig
	cloudSyncer        CloudSyncerConfig
}

func NewConfig() *Config {
	// antvip server use '-client' suffix for TrFrom to judge request source is client, '-client' can not be changed
	c := &Config{
		appName:            "sofaantvip",
		datacenter:         "",
		zone:               "",
		instanceID:         "",
		endpoint:           "",
		accessKey:          "",
		accessSecret:       "",
		trFrom:             "go-client",
		version:            "1.0.4",
		domainChecksumMode: DomainAlipayChecksumMode,
		metrics:            &Metrics{},
	}

	c.syncTimeout = DefaultSyncerTimeout
	c.syncInterval = DefaultSyncerInterval

	c.cloudSyncer.https = false
	c.cloudSyncer.address = DefaultCloudSyncerServerAddress
	c.cloudSyncer.port = DefaultCloudSyncerPort
	c.cloudSyncer.endpoint = DefaultCloudSyncerEndpoint

	c.alipaySyncer.https = false
	c.alipaySyncer.port = DefaultAlipaySyncerPort
	c.alipaySyncer.endpoint = DefaultAlipaySyncerEndpoint

	c.httpLocator.https = false
	c.httpLocator.timeout = DefaultHTTPLocatorTimeout
	c.httpLocator.address = DefaultHTTPLocatorAntVipServerAddress
	c.httpLocator.interval = DefaultHTTPLocatorInterval
	c.httpLocator.port = DefaultHTTPLocatorAntVipHTTPPort
	c.httpLocator.endpoint = DefaultHTTPLocatorAntVipEndpoint
	c.registryLocator.timeout = DefaultRegistryLocatorTimeout

	return c
}

func (o *Config) GetDomainChecksumMode() DomainChecksumMode { return o.domainChecksumMode }

func (o *Config) SetDomainChecksumMode(m DomainChecksumMode) *Config {
	o.domainChecksumMode = m
	return o
}

func (o *Config) SetMetrics(m *Metrics) *Config {
	o.metrics = m
	return o
}

func (o *Config) GetMetrics() *Metrics { return o.metrics }

func (o *Config) GetAppName() string { return o.appName }

func (o *Config) SetAppName(n string) *Config {
	o.appName = n
	return o
}

func (o *Config) GetDataCenter() string { return o.datacenter }

func (o *Config) SetDataCenter(n string) *Config {
	o.datacenter = n
	return o
}

func (o *Config) GetZone() string { return o.zone }

func (o *Config) SetZone(zone string) *Config {
	o.zone = zone
	return o
}

func (o *Config) GetInstanceId() string { return o.instanceID }

func (o *Config) SetInstanceID(id string) *Config {
	o.instanceID = id
	return o
}

func (o *Config) GetEndpoint() string { return o.endpoint }

func (o *Config) SetEndpoint(e string) *Config {
	o.endpoint = e
	return o
}

func (o *Config) GetAccessKey() string { return o.accessKey }

func (o *Config) SetAccesskey(a string) *Config {
	o.accessKey = a
	return o
}

func (o *Config) GetSecretKey() string { return o.accessSecret }

func (o *Config) SetSecretKey(s string) *Config {
	o.accessSecret = s
	return o
}

func (o *Config) GetTRFrom() string { return o.trFrom }

func (o *Config) SetTRFrom(f string) *Config {
	o.trFrom = f
	return o
}

func (o *Config) GetVersion() string { return o.version }

func (o *Config) SetVersion(v string) *Config {
	o.version = v
	return o
}

func (o *Config) GetSyncTimeout() time.Duration { return o.syncTimeout }

func (o *Config) SetSyncTimeout(t time.Duration) *Config {
	o.syncTimeout = t
	return o
}

func (o *Config) GetSyncInterval() time.Duration { return o.syncInterval }

func (o *Config) SetSyncInterval(interval time.Duration) *Config {
	o.syncInterval = interval
	return o
}

func (o *Config) GetHTTPLocatorConfig() *HTTPLocatorConfig {
	return &o.httpLocator
}

func (o *Config) GetHTTPLocatorURL() string {
	scheme := "http"
	if o.httpLocator.https {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s:%d%s?zone=%s", scheme, o.httpLocator.address, o.httpLocator.port,
		o.httpLocator.endpoint, o.zone)
}

func (o *Config) GetCloudSyncerConfig() *CloudSyncerConfig {
	return &o.cloudSyncer
}

func (o *Config) GetAlipaySyncerConfig() *AlipaySyncerConfig {
	return &o.alipaySyncer
}

func (o *Config) GetRegistryLocatorConfig() *RegistryLocatorConfig {
	return &o.registryLocator
}
