package filter

import (
	v1 "gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	"strings"

	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/common/constants"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/types"
)

func init() {
	Register("AddRequestHeaderFilter", &AddRequestHeaderFilterFactory{})
	Register("AddResponseHeaderFilter", &AddResponseHeaderFilterFactory{})
}

type AddRequestHeaderFilterFactory struct {
}

func (*AddRequestHeaderFilterFactory) CreateFilter(filter *v1.Filter) types.GatewayFilter {
	return &AddHeaderFilter{
		base:      filter,
		headers:   parseConfigHeaders(filter),
		direction: constants.Request,
	}
}

type AddResponseHeaderFilterFactory struct {
}

func (*AddResponseHeaderFilterFactory) CreateFilter(filter *v1.Filter) types.GatewayFilter {
	return &AddHeaderFilter{
		base:      filter,
		headers:   parseConfigHeaders(filter),
		direction: constants.Response,
	}
}

func parseConfigHeaders(filter *v1.Filter) []Header {
	var headers []Header
	if value, ok := filter.Metadata.(string); ok {
		hs := strings.Split(strings.TrimSpace(value), ",")
		for _, header := range hs {
			kv := strings.Split(strings.TrimSpace(header), ":")
			if len(kv) == 2 {
				headers = append(headers, Header{
					kv[0],
					kv[1],
				})
			}
		}
	}
	return headers
}

type Header struct {
	key   string
	value string
}

type AddHeaderFilter struct {
	base *v1.Filter

	direction constants.Direction
	headers   []Header
}

func (filter *AddHeaderFilter) Name() string {
	return filter.base.Name
}

func (filter *AddHeaderFilter) Priority() int64 {
	return filter.base.Priority
}

func (filter *AddHeaderFilter) InBound(ctx types.Context) (types.FilterStatus, error) {
	if filter.direction == constants.Request {
		for _, header := range filter.headers {
			ctx.Request().SetHeader(header.key, header.value)
		}
	}

	return types.Success, nil
}

func (filter *AddHeaderFilter) OutBound(ctx types.Context) (types.FilterStatus, error) {
	if filter.direction == constants.Response {
		for _, header := range filter.headers {
			ctx.Response().SetHeader(header.key, header.value)
		}
	}
	return types.Success, nil
}
