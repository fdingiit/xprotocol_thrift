package sofadrm

import (
	"gitlab.alipay-inc.com/sofa-go/sofa-antvip-client-go/sofaantvip"
	"gitlab.alipay-inc.com/sofa-go/sofa-logger-go/sofalogger"
)

type AntvipLocatortOptionSetter interface {
	set(*AntvipLocator)
}

type AntvipLocatortOptionSetterFunc func(*AntvipLocator)

func (f AntvipLocatortOptionSetterFunc) set(c *AntvipLocator) {
	f(c)
}

func WithAntvipLocatorLogger(logger sofalogger.Logger) AntvipLocatortOptionSetterFunc {
	return AntvipLocatortOptionSetterFunc(func(c *AntvipLocator) {
		c.logger = logger
	})
}

func WithAntvipLocatorConfig(config *Config) AntvipLocatortOptionSetterFunc {
	return AntvipLocatortOptionSetterFunc(func(c *AntvipLocator) {
		c.config = config
	})
}

func WithAntvipLocatorClient(client *sofaantvip.AntvipClient) AntvipLocatortOptionSetterFunc {
	return AntvipLocatortOptionSetterFunc(func(c *AntvipLocator) {
		c.client = client
	})
}
