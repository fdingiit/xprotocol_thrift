package zoneclient

import "time"

const (
	DefaultAlipayRouterAntvipDomain  = "zonemng-pool.stable.alipay.net"
	DefaultAlipayRouterAntvipTimeout = 5 * time.Second
)

type AlipayRouterConfig struct {
	domain          string
	timeout         time.Duration
	zoneRoutePath   string
	elasticRulePath string
}

func (c *AlipayRouterConfig) SetDomain(domain string) *AlipayRouterConfig {
	c.domain = domain
	return c
}

func (al *AlipayRouterConfig) SetZoneRoutePath(s string) *AlipayRouterConfig {
	al.zoneRoutePath = s
	return al
}

func (al *AlipayRouterConfig) GetZoneRoutePath() string { return al.zoneRoutePath }

func (al *AlipayRouterConfig) SetElasticRulePath(s string) *AlipayRouterConfig {
	al.elasticRulePath = s
	return al
}

func (al *AlipayRouterConfig) GetElasticRulePath() string { return al.elasticRulePath }

func (al *AlipayRouterConfig) GetDomain() string { return al.domain }

func (al *AlipayRouterConfig) GetTimeout() time.Duration { return al.timeout }

type Config struct {
	appName      string
	zone         string
	alipayRouter AlipayRouterConfig
}

func NewConfig() *Config {
	cf := &Config{}
	cf.alipayRouter.domain = DefaultAlipayRouterAntvipDomain
	cf.alipayRouter.timeout = DefaultAlipayRouterAntvipTimeout

	return cf
}

func (o *Config) GetAlipayRouterConfig() *AlipayRouterConfig {
	return &o.alipayRouter
}

func (o *Config) GetZone() string { return o.zone }

func (o *Config) SetZone(zone string) *Config {
	o.zone = zone
	return o
}

func (o *Config) GetAppName() string { return o.appName }

func (o *Config) SetAppName(n string) *Config {
	o.appName = n
	return o
}
