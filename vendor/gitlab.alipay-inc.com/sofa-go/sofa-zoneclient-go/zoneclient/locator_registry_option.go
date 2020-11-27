package zoneclient

import (
	"gitlab.alipay-inc.com/sofa-go/sofa-logger-go/sofalogger"
	"gitlab.alipay-inc.com/sofa-go/sofa-registry-client-go/sofaregistry"
)

type RegistryLocatortOptionSetter interface {
	set(*RegistryLocator)
}

type RegistryLocatortOptionSetterFunc func(*RegistryLocator)

func (f RegistryLocatortOptionSetterFunc) set(c *RegistryLocator) {
	f(c)
}

func WithRegistryLocatorLogger(logger sofalogger.Logger) RegistryLocatortOptionSetterFunc {
	return RegistryLocatortOptionSetterFunc(func(c *RegistryLocator) {
		c.logger = logger
	})
}

func WithRegistryLocatorConfig(config *Config) RegistryLocatortOptionSetterFunc {
	return RegistryLocatortOptionSetterFunc(func(c *RegistryLocator) {
		c.config = config
	})
}

func WithRegistryLocatorClient(client *sofaregistry.Client) RegistryLocatortOptionSetterFunc {
	return RegistryLocatortOptionSetterFunc(func(c *RegistryLocator) {
		c.client = client
	})
}
