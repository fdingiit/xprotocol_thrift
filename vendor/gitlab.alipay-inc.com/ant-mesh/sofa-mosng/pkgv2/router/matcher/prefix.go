package matcher

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/router"
	"mosn.io/mosn/pkg/protocol"
	mosn "mosn.io/mosn/pkg/types"
	"strings"
)

func init() {
	router.GetMatcherFactory().Register(newPrefixMatcher)
}

type prefixMatcher struct {
	prefix string
}

func (pmr *prefixMatcher) Match(headers mosn.HeaderMap) bool {
	if headerPathValue, ok := headers.Get(protocol.MosnHeaderPathKey); ok {
		if strings.HasPrefix(headerPathValue, pmr.prefix) {
			return true
		}
	}
	return false
}

func newPrefixMatcher(conf *v1.RouterMatch) (router.Matcher, error) {
	if conf.Prefix != "" {
		return &prefixMatcher{conf.Prefix}, nil
	}

	return nil, nil
}
