package sofaantvip

import (
	"context"

	"gitlab.alipay-inc.com/sofa-go/sofa-logger-go/sofalogger"
)

// AntvipClientOptionSetter configures a client.
type AntvipClientOptionSetter interface {
	set(*AntvipClient)
}

type AntvipClientOptionSetterFunc func(*AntvipClient)

func (f AntvipClientOptionSetterFunc) set(c *AntvipClient) {
	f(c)
}

func WithAntvipClientSyncer(syncer Syncer) AntvipClientOptionSetterFunc {
	return AntvipClientOptionSetterFunc(func(c *AntvipClient) {
		c.syncer = syncer
	})
}

func WithAntvipClientLogger(logger sofalogger.Logger) AntvipClientOptionSetterFunc {
	return AntvipClientOptionSetterFunc(func(c *AntvipClient) {
		c.logger = logger
	})
}

func WithAntvipClientConfig(config *Config) AntvipClientOptionSetterFunc {
	return AntvipClientOptionSetterFunc(func(c *AntvipClient) {
		c.config = config
	})
}

func WithAntvipClientContext(ctx context.Context) AntvipClientOptionSetterFunc {
	return AntvipClientOptionSetterFunc(func(c *AntvipClient) {
		c.context = ctx
	})
}
