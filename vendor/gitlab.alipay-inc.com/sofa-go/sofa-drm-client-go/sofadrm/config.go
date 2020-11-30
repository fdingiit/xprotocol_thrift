package sofadrm

import (
	"time"
)

const (
	DefaultDRMDomain                         = "zdrmdata-pool"
	DefaultBOLTTransportTimeout              = 10 * time.Second
	DefaultBOLTTransportMaxHeartbeatAttempts = 3
	DefaultBOLTTransportMaxPendingCommands   = 128
	DefaultHeartbeatInterval                 = 30 * time.Second
	DefaultRegisterInterval                  = 30 * time.Second
	DefaultRegistryTimeout                   = 3 * time.Second
)

type RegistryLocatorConfig struct {
	timeout time.Duration
}

func (rl *RegistryLocatorConfig) SetTimeout(t time.Duration) {
	rl.timeout = t
}

type BOLTTransportConfig struct {
	timeout              time.Duration
	maxPendingCommands   int
	maxHeartbeatAttempts int
}

func (b *BOLTTransportConfig) SetMaxPendingCommands(maxPendingCommands int) *BOLTTransportConfig {
	b.maxPendingCommands = maxPendingCommands
	return b
}

func (b *BOLTTransportConfig) GetMaxPendingCommands() int {
	return b.maxPendingCommands
}

func (b *BOLTTransportConfig) SetTimeout(t time.Duration) *BOLTTransportConfig {
	b.timeout = t
	return b
}

func (b *BOLTTransportConfig) GetTimeout() time.Duration {
	return b.timeout
}

type AntvipLocatorConfig struct {
	domain  string
	timeout time.Duration
}

func (a *AntvipLocatorConfig) SetTimeout(t time.Duration) *AntvipLocatorConfig {
	a.timeout = t
	return a
}

func (a *AntvipLocatorConfig) GetTimeout() time.Duration { return a.timeout }

func (a *AntvipLocatorConfig) GetDomain() string {
	return a.domain
}

func (a *AntvipLocatorConfig) SetDomain(host string) *AntvipLocatorConfig {
	a.domain = host
	return a
}

type Config struct {
	instanceID        string
	zone              string
	dataCenter        string
	appName           string
	profile           string
	advertiseIFace    string
	accessKey         string
	secretKey         string
	heartbeatInterval time.Duration
	registryLocator   RegistryLocatorConfig
	boltTransport     BOLTTransportConfig
	antvipLocator     AntvipLocatorConfig
}

func NewConfig() *Config {
	c := &Config{}
	c.boltTransport.timeout = DefaultBOLTTransportTimeout
	c.boltTransport.maxPendingCommands = DefaultBOLTTransportMaxPendingCommands
	c.boltTransport.maxHeartbeatAttempts = DefaultBOLTTransportMaxHeartbeatAttempts
	c.antvipLocator.timeout = DefaultBOLTTransportTimeout
	c.antvipLocator.domain = DefaultDRMDomain
	c.registryLocator.timeout = DefaultRegistryTimeout
	c.heartbeatInterval = DefaultHeartbeatInterval
	return c
}

func (o *Config) GetAccessKey() string {
	return o.accessKey
}

func (o *Config) SetAccessKey(ak string) *Config {
	o.accessKey = ak
	return o
}

func (o *Config) GetSecretKey() string {
	return o.secretKey
}

func (o *Config) SetSecretKey(sk string) *Config {
	o.secretKey = sk
	return o
}

func (o *Config) SetAdvertiseIFace(iface string) *Config {
	o.advertiseIFace = iface
	return o
}

func (o *Config) GetAppName() string { return o.appName }

func (o *Config) SetAppName(n string) *Config {
	o.appName = n
	return o
}

func (o *Config) GetDataCenter() string { return o.dataCenter }

func (o *Config) SetDataCenter(n string) *Config {
	o.dataCenter = n
	return o
}

func (o *Config) GetInstanceID() string { return o.instanceID }

func (o *Config) SetInstanceID(id string) *Config {
	o.instanceID = id
	return o
}

func (o *Config) GetZone() string { return o.zone }

func (o *Config) SetZone(zone string) *Config {
	o.zone = zone
	return o
}

func (o *Config) GetProfile() string { return o.profile }

func (o *Config) SetProfile(profile string) *Config {
	o.profile = profile
	return o
}

func (c *Config) GetAntvipLocatorConfig() *AntvipLocatorConfig {
	return &c.antvipLocator
}

func (c *Config) GetRegistryLocatorConfig() *RegistryLocatorConfig {
	return &c.registryLocator
}

func (c *Config) GetBOLTTransportConfig() *BOLTTransportConfig {
	return &c.boltTransport
}

func (c *Config) SetHeartbeatInterval(interval time.Duration) *Config {
	c.heartbeatInterval = interval
	return c
}

func (c *Config) GetHeartbeatInterval() time.Duration { return c.heartbeatInterval }
