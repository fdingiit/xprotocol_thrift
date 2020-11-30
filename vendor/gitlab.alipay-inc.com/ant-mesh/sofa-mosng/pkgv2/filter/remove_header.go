package filter

import (
	v1 "gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	"strings"

	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/common/constants"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/types"
)

func init() {
	Register("RemoveRequestHeaderFilter", &RemoveRequestHeaderFilterFactory{})
	Register("RemoveResponseHeaderFilter", &RemoveResponseHeaderFilterFactory{})
}

type RemoveRequestHeaderFilterFactory struct {
}

func (*RemoveRequestHeaderFilterFactory) CreateFilter(filter *v1.Filter) types.GatewayFilter {
	return &RemoveHeaderFilter{
		base:       filter,
		direction:  constants.Request,
		headerKeys: parseConfigHeaderKeys(filter),
	}
}

type RemoveResponseHeaderFilterFactory struct {
}

func (*RemoveResponseHeaderFilterFactory) CreateFilter(filter *v1.Filter) types.GatewayFilter {
	return &RemoveHeaderFilter{
		base:       filter,
		direction:  constants.Response,
		headerKeys: parseConfigHeaderKeys(filter),
	}
}

func parseConfigHeaderKeys(filter *v1.Filter) []string {
	headerKeys := []string{}
	if value, ok := filter.Metadata.(string); ok {
		hs := strings.Split(strings.TrimSpace(value), ",")
		for _, header := range hs {
			headerKeys = append(headerKeys, strings.TrimSpace(header))
		}
	}
	return headerKeys
}

type RemoveHeaderFilter struct {
	base *v1.Filter

	direction  constants.Direction
	headerKeys []string
}

func (filter *RemoveHeaderFilter) Name() string {
	return filter.base.Name
}

func (filter *RemoveHeaderFilter) Priority() int64 {
	return filter.base.Priority
}

func (filter *RemoveHeaderFilter) InBound(ctx types.Context) (types.FilterStatus, error) {
	if filter.direction == constants.Request {
		for _, header := range filter.headerKeys {
			ctx.Request().DelHeader(header)
		}
	}
	return types.Success, nil
}

func (filter *RemoveHeaderFilter) OutBound(ctx types.Context) (types.FilterStatus, error) {
	if filter.direction == constants.Response {
		for _, header := range filter.headerKeys {
			ctx.Response().DelHeader(header)
		}
	}
	return types.Success, nil
}
