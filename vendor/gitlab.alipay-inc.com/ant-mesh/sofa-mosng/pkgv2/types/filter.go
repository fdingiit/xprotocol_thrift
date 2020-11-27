package types

import (
	v1 "gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
)

type GatewayFilter interface {
	// filter name
	Name() string

	// the priority of filter 0 - 1000
	Priority() int64

	// invoke before proxy
	InBound(ctx Context) (FilterStatus, error)

	// invoke after proxy
	OutBound(ctx Context) (FilterStatus, error)
}

type FilterFactory interface {
	CreateFilter(conf *v1.Filter) GatewayFilter
}

type FilterStatus string

var (
	Success FilterStatus = "SUCCESS"
	Error   FilterStatus = "ERROR"
)
