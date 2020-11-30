package router

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/common/log"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/types"
	mosn "mosn.io/mosn/pkg/types"
)

var (
	matcherFactoryIns = &matcherFactory{
		matcherCreators: []matcherCreator{},
	}
)

type router struct {
	matchers []Matcher
	pipeline types.Pipeline
	service  types.Service
	conf     v1.Router
}

func (br *router) Conf() v1.Router {
	return br.conf
}

func (br *router) Pipeline() types.Pipeline {
	return br.pipeline.Copy()
}

func (br *router) Service() types.Service {
	return br.service
}

func (br *router) Match(headers mosn.HeaderMap) bool {
	for _, matcher := range br.matchers {
		if match := matcher.Match(headers); match {
			continue
		}
		return false
	}
	return true
}

type weightRouter struct {
	*router
	pipeline types.WeightPipeline
}

func (br *weightRouter) Pipeline() types.Pipeline {
	if br.pipeline != nil {
		return br.pipeline.GetPipeline()
	}
	return br.router.pipeline
}

func GetMatcherFactory() *matcherFactory {
	return matcherFactoryIns
}

type matcherFactory struct {
	matcherCreators []matcherCreator
}

type matcherCreator func(conf *v1.RouterMatch) (Matcher, error)

func (mf *matcherFactory) Register(f matcherCreator) {
	mf.matcherCreators = append(mf.matcherCreators, f)
}

func (mf *matcherFactory) GetCreators() []matcherCreator {
	return mf.matcherCreators
}

func (mf *matcherFactory) CreateMatchers(conf *v1.RouterMatch) []Matcher {
	var matchers []Matcher

	for _, f := range mf.matcherCreators {
		if matcher, err := f(conf); err == nil && matcher != nil {
			matchers = append(matchers, matcher)
		} else {
			log.ConfigLogger().Errorf("create router matcher err, %v, routerMatch: %v", err, conf)
		}
	}

	if len(matchers) == 0 {
		log.ConfigLogger().Warnf("create 0 router matcher, routerMatch: %v", conf)
	}

	return matchers
}
