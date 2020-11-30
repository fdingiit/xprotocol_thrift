package types

import (
	v1 "gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	mosn "mosn.io/mosn/pkg/types"
)

type RouterManager interface {
	MatchRouter() Router
	AddRouter(router Router)
	DeleteRouter(router Router)
}

type Router interface {
	Match(headers mosn.HeaderMap) bool
	Pipeline() Pipeline
	Service() Service
	Conf() v1.Router
}
