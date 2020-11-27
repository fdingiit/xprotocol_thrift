package zoneclient

import (
	"gitlab.alipay-inc.com/sofa-go/sofa-drm-client-go/sofadrm"
	"gitlab.alipay-inc.com/sofa-go/sofa-logger-go/sofalogger"
)

// AlipayRouterOptionSetter configures a AlipayRouter.
type AlipayRouterOptionSetter interface {
	set(*AlipayRouter)
}

type AlipayRouterOptionSetterFunc func(*AlipayRouter)

func (f AlipayRouterOptionSetterFunc) set(c *AlipayRouter) {
	f(c)
}

func WithAlipayRouterConfig(config *Config) AlipayRouterOptionSetterFunc {
	return AlipayRouterOptionSetterFunc(func(c *AlipayRouter) {
		c.config = config
	})
}

func WithAlipayRouterLogger(logger sofalogger.Logger) AlipayRouterOptionSetterFunc {
	return AlipayRouterOptionSetterFunc(func(c *AlipayRouter) {
		c.logger = logger
	})
}

func WithAlipayRouterLocator(locator Locator) AlipayRouterOptionSetterFunc {
	return AlipayRouterOptionSetterFunc(func(c *AlipayRouter) {
		c.locator = locator
	})
}

func WithAlipayRouterDRMClient(client *sofadrm.Client) AlipayRouterOptionSetterFunc {
	return AlipayRouterOptionSetterFunc(func(c *AlipayRouter) {
		c.drm = client
	})
}
