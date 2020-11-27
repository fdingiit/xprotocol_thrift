package filter

import (
	v1 "gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	"math"

	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/errors"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/types"
	"golang.org/x/time/rate"
	"mosn.io/mosn/pkg/protocol/http"
)

func init() {
	Register("RateLimiterFilter", &RequestRateLimiterFilterFactory{})
}

type RequestRateLimiterFilterFactory struct {
}

func (*RequestRateLimiterFilterFactory) CreateFilter(filter *v1.Filter) types.GatewayFilter {
	rateSpeed := filter.Metadata.(float64)
	return &RateLimiterFilter{
		base:    filter,
		limiter: rate.NewLimiter(rate.Limit(rateSpeed), int(math.Ceil(rateSpeed))),
	}
}

type RateLimiterFilter struct {
	base *v1.Filter

	limiter *rate.Limiter
}

func (f *RateLimiterFilter) Name() string {
	return f.base.Name
}

func (f *RateLimiterFilter) Priority() int64 {
	return f.base.Priority
}

func (f *RateLimiterFilter) InBound(ctx types.Context) (types.FilterStatus, error) {
	if !f.limiter.Allow() {
		// todo
		return types.Error, errors.New(http.TooManyRequests)
	}
	return types.Success, nil
}

func (f *RateLimiterFilter) OutBound(ctx types.Context) (types.FilterStatus, error) {
	panic("implement me")
}
