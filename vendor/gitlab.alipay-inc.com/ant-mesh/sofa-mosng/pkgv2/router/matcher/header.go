package matcher

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/router"
	mosn "mosn.io/mosn/pkg/types"
	"regexp"
	"strings"
)

func init() {
	router.GetMatcherFactory().Register(newHeaderMatcher)
}

type headerMatcher struct {
	headerMatch []*v1.ValueMatcher
}

func (hmr *headerMatcher) Match(headers mosn.HeaderMap) bool {
	for _, header := range hmr.headerMatch {
		if headerPathValue, ok := headers.Get(header.Name); ok {
			if header.Regex {
				// 如果是正则匹配
				matched, err := regexp.MatchString(header.Value, headerPathValue)

				if err != nil || !matched {
					//todo log
					return false
				}
			} else if !strings.EqualFold(headerPathValue, header.Value) {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func newHeaderMatcher(conf *v1.RouterMatch) (router.Matcher, error) {
	if conf.Headers != nil && len(conf.Headers) > 0 {
		return &headerMatcher{conf.Headers}, nil
	}

	return nil, nil
}
