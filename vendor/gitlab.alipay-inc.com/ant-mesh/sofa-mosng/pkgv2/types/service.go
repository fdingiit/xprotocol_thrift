package types

import (
	v1 "gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
)

type Service interface {
	Name() string
	Upstream() Upstream
	Conf() *v1.GatewayService
}
