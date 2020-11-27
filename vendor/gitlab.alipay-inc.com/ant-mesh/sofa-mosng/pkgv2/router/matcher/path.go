package matcher

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/router"
	"mosn.io/mosn/pkg/protocol"
	mosn "mosn.io/mosn/pkg/types"
	"strings"
)

func init() {
	router.GetMatcherFactory().Register(newPathMatcher)
}

type pathMatcher struct {
	path string
}

func (pmr *pathMatcher) Match(headers mosn.HeaderMap) bool {
	if headerPathValue, ok := headers.Get(protocol.MosnHeaderPathKey); ok {
		if strings.EqualFold(headerPathValue, pmr.path) {
			return true
		}
	}
	return false
}

func newPathMatcher(conf *v1.RouterMatch) (router.Matcher, error) {
	if conf.Path != "" {
		return &pathMatcher{conf.Path}, nil
	}

	return nil, nil
}
