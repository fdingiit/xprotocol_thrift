package codec

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/types"
	"sync"
)

type codecFactory struct {
	cMux     sync.RWMutex
	codecMap map[string]types.Codec
}

func (c *codecFactory) Register(codec types.Codec) {
	c.cMux.Lock()
	defer c.cMux.Unlock()

	c.codecMap[codec.Protocol()] = codec
}

func (c *codecFactory) GetCodec(protocol string) types.Codec {
	return c.codecMap[protocol]
}

var (
	codecFactoryInstance types.CodecFactory
	once                 sync.Once
)

func newCodecFactory() types.CodecFactory {
	return &codecFactory{
		cMux:     sync.RWMutex{},
		codecMap: make(map[string]types.Codec),
	}
}

func GetCodeFactoryInstance() types.CodecFactory {
	once.Do(func() {
		codecFactoryInstance = newCodecFactory()
	})
	return codecFactoryInstance
}
