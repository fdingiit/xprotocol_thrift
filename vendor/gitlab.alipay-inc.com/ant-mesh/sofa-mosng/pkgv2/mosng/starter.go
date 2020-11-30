package mosng

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/common/log"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/mosn"
	"sync"

	_ "gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/valid"
	_ "gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/extension"
	_ "gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/filter"
	_ "gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/mosn"
	_ "gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/router/matcher"
	_ "gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/types"
)

var (
	mosng *Mosng
	once  sync.Once
)

type Mosng struct {
}

func New() *Mosng {
	once.Do(func() {
		mosng = &Mosng{}
	})

	return mosng
}

func (*Mosng) Start() {
	mosn.ConvertInit()
}

func (*Mosng) StartWithOutConvert() {
	// noop now
	log.StartLogger().Infof("[mosng][start] start without convert")
}
