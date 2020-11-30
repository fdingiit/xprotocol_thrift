package sofaregistry

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	DefaultBOLTTransportMaxPendingCommands   = 128
	DefaultBOLTTransportMaxHeartbeatAttempts = 3
	DefaultBOLTTransportHost                 = "confreg-pool"
	DefaultBOLTTransportPort                 = 9603
	DefaultBOLTTransportTimeout              = 10 * time.Second
	DefaultGroup                             = "DEFAULT_GROUP"
	DefaultZone                              = "DEFAULT_ZONE"
)

type BOLTTransportConfig struct {
	host                 string
	port                 int
	timeout              time.Duration
	maxPendingCommands   int
	maxHeartbeatAttempts int
}

func (b *BOLTTransportConfig) GetMaxPendingCommands() int {
	return b.maxPendingCommands
}

func (b *BOLTTransportConfig) SetMaxPendingCommands(m int) *BOLTTransportConfig {
	b.maxPendingCommands = m
	return b
}

func (b *BOLTTransportConfig) SetPort(p int) *BOLTTransportConfig {
	b.port = p
	return b
}

func (b *BOLTTransportConfig) SetHost(h string) *BOLTTransportConfig {
	b.host = h
	return b
}

func (b *BOLTTransportConfig) SetTimeout(t time.Duration) *BOLTTransportConfig {
	b.timeout = t
	return b
}

func (b *BOLTTransportConfig) GetTimeout() time.Duration {
	return b.timeout
}

type Config struct {
	env                 string
	zone                string
	dataCenter          string
	appName             string
	instanceID          string
	accessKey           string
	secretKey           string
	boltTransportConfig BOLTTransportConfig
}

func NewConfig() *Config {
	c := &Config{}
	c.boltTransportConfig = BOLTTransportConfig{
		host:                 DefaultBOLTTransportHost,
		port:                 DefaultBOLTTransportPort,
		timeout:              DefaultBOLTTransportTimeout,
		maxPendingCommands:   DefaultBOLTTransportMaxPendingCommands,
		maxHeartbeatAttempts: DefaultBOLTTransportMaxHeartbeatAttempts,
	}
	return c
}

func (o *Config) GetEnv() string { return o.env }

func (o *Config) SetEnv(n string) *Config {
	o.env = n
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

func (o *Config) GetAccessKey() string { return o.accessKey }

func (o *Config) SetAccesskey(a string) *Config {
	o.accessKey = a
	return o
}

func (o *Config) GetSecretKey() string { return o.secretKey }

func (o *Config) SetSecretKey(s string) *Config {
	o.secretKey = s
	return o
}

func (c *Config) GetSignature() map[string]string {
	if c.accessKey == "" {
		return nil
	}
	return getSignatureMap(c.accessKey, c.secretKey, c.instanceID)
}

func (c *Config) GetBOLTTransportConfig() *BOLTTransportConfig {
	return &c.boltTransportConfig
}

func (c *Config) GetBOLTTransportLocatorURL() string {
	return fmt.Sprintf("http://%s:%d/api/servers/query?env=%s&zone=%s&dataCenter=%s&appName=%s&instanceId=%s",
		c.boltTransportConfig.host,
		c.boltTransportConfig.port,
		c.env,
		c.zone,
		c.dataCenter,
		c.appName,
		c.instanceID)
}

func (c *Config) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"env":        c.env,
		"zone":       c.zone,
		"dataCenter": c.dataCenter,
		"appName":    c.appName,
		"instanceID": c.instanceID,
		"accessKey":  c.accessKey,
		"secretKey":  c.secretKey,
		"bolt_host":  c.boltTransportConfig.host,
		"bolt_port":  fmt.Sprintf("%d", c.boltTransportConfig.port),
	})
}
