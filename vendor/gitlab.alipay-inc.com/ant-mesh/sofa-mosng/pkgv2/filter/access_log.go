package filter

import (
	v1 "gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/types"
)

func init() {
	Register("AccessLogFilter", &AccessLogFilterFactory{})
}

type AccessLogFilter struct {
	priority int64
	name     string
}

type AccessLogFilterFactory struct{}

func (alff *AccessLogFilterFactory) CreateFilter(filter *v1.Filter) types.GatewayFilter {
	return &AccessLogFilter{name: "AccessLogFilter", priority: filter.Priority}
}

func (alf *AccessLogFilter) Name() string {
	return alf.name
}

func (alf *AccessLogFilter) Priority() int64 {
	return alf.priority
}

func (alf *AccessLogFilter) InBound(ctx types.Context) (types.FilterStatus, error) {
	return types.Success, nil
}

func (alf *AccessLogFilter) OutBound(ctx types.Context) (types.FilterStatus, error) {
	return types.Success, nil
}
