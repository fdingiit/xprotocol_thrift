package matcher

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/errors"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/router"
	"mosn.io/mosn/pkg/protocol"
	mosn "mosn.io/mosn/pkg/types"
	"regexp"
	"strings"
)

func init() {
	router.GetMatcherFactory().Register(newQueryStringMatcher)
}

const (
	ParamsSplit = "&"
	ParamSplit  = "="
)

// todo support array

type queryStringMatcher struct {
	queryStringMatch []*v1.ValueMatcher
}

func (pm *queryStringMatcher) Match(headers mosn.HeaderMap) bool {
	params := map[string]string{}

	// build params map
	// /name=jack&sex=man
	if queryStringValue, ok := headers.Get(protocol.MosnHeaderQueryStringKey); ok {
		buildParamMap(queryStringValue, params)
	} else {
		return false
	}

	// params map match
	for _, parameter := range pm.queryStringMatch {
		if param := params[parameter.Name]; param != "" {
			if parameter.Regex {
				// 如果是正则匹配
				matched, err := regexp.MatchString(parameter.Value, param)

				if err != nil || !matched {
					//todo log
					return false
				}
			} else if !strings.EqualFold(param, parameter.Value) {
				return false
			}
		} else {
			return false
		}
	}
	return true
}

func buildParamMap(queryStringValue string, params map[string]string) {
	// name=jack sex=man
	splits := strings.Split(queryStringValue, ParamsSplit)
	for _, split := range splits {
		// name jack
		// sex man
		kv := strings.Split(split, ParamSplit)
		if len(kv) != 2 {
			panic(errors.Errorf("request parameters are wrong [%s]", queryStringValue))
		}
		params[kv[0]] = kv[1]
	}
}

func newQueryStringMatcher(conf *v1.RouterMatch) (router.Matcher, error) {
	if conf.QueryStrings != nil && len(conf.QueryStrings) > 0 {
		return &queryStringMatcher{conf.QueryStrings}, nil
	}

	return nil, nil
}
