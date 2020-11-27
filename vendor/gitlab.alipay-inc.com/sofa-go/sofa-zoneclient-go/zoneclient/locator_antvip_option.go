package zoneclient

import (
	"gitlab.alipay-inc.com/sofa-go/sofa-antvip-client-go/sofaantvip"
	"gitlab.alipay-inc.com/sofa-go/sofa-logger-go/sofalogger"
)

type AntVipLocatortOptionSetter interface {
	set(locator *AntVipLocator)
}

type AntVipLocatortOptionSetterFunc func(*AntVipLocator)

func (f AntVipLocatortOptionSetterFunc) set(c *AntVipLocator) {
	f(c)
}

func WithAntVipLocatorLogger(logger sofalogger.Logger) AntVipLocatortOptionSetterFunc {
	return AntVipLocatortOptionSetterFunc(func(c *AntVipLocator) {
		c.logger = logger
	})
}

func WithAntVipLocatorConfig(config *Config) AntVipLocatortOptionSetterFunc {
	return AntVipLocatortOptionSetterFunc(func(c *AntVipLocator) {
		c.config = config
	})
}

func WithAntVipLocatorClient(client *sofaantvip.AntvipClient) AntVipLocatortOptionSetterFunc {
	return AntVipLocatortOptionSetterFunc(func(c *AntVipLocator) {
		c.client = client
	})
}
