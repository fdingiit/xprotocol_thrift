package router

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/common/log"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api"
	v1 "gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/event"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/errors"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/filter/pipeline"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/service"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/types"
	mosn "mosn.io/mosn/pkg/types"
	"strings"
	"sync"
)

func init() {
	event.EventListenerManagerInstance().RegisterList([]event.ResourceEventListener{
		event.ResourceEventListenerFuncs{
			Type:       api.ROUTER,
			AddFunc:    GetRouterManager().AddOrUpdateRouter,
			UpdateFunc: GetRouterManager().AddOrUpdateRouter,
		},
		event.ResourceEventListenerFuncs{
			Type:       api.ROUTER_GROUP,
			AddFunc:    GetRouterManager().AddOrUpdateRouterGroup,
			UpdateFunc: GetRouterManager().AddOrUpdateRouterGroup,
		},
	})
}

var (
	routerManagerInstance *Manager
	once                  sync.Once
)

func GetRouterManager() *Manager {
	once.Do(func() {
		routerManagerInstance = newRouterManager()
	})

	return routerManagerInstance
}

type routerMatcher func(string, mosn.HeaderMap) types.Router
type routerFactory func(rConf *v1.Router) (types.Router, error)

type Manager struct {
	rMux            sync.RWMutex
	routerMap       map[string]map[string]types.Router // map[groupName]map[routerName]router
	sMux            sync.RWMutex
	serverRouterMap map[string]map[string]types.Router // map[serverName]map[routerName]router
	routerMatcher   routerMatcher
	routerFactory   routerFactory
}

func (m *Manager) SetRouterMatcher(rm routerMatcher) {
	m.routerMatcher = rm
}

func (m *Manager) SetRouterFactory(rf routerFactory) {
	m.routerFactory = rf
}

func (m *Manager) Match(listenerName string, headers mosn.HeaderMap) types.Router {
	return m.routerMatcher(listenerName, headers)
}

func (m *Manager) AddOrUpdateRouterGroup(o api.Object) (error, bool) {
	rgConf := o.(*v1.RouterGroup)
	for _, rConf := range rgConf.Routers {
		rConf.GroupBind = rgConf.GroupName

		if err, ok := m.addOrUpdateRouter(rConf, rgConf.Gateways); !ok {
			return err, ok
		}
	}

	return nil, true
}

func (m *Manager) GetRouters(routerGroupName string) map[string]types.Router {
	m.rMux.RLock()
	defer m.rMux.RUnlock()
	return m.routerMap[routerGroupName]
}

func (m *Manager) GetServerRouters(serverName string) map[string]types.Router {
	m.sMux.RLock()
	defer m.sMux.RUnlock()
	return m.serverRouterMap[serverName]
}

func (m *Manager) AddOrUpdateRouter(o api.Object) (error, bool) {
	rConf := o.(*v1.Router)
	return m.addOrUpdateRouter(rConf, nil)

}

func (m *Manager) addOrUpdateRouter(rConf *v1.Router, serverNames []string) (error, bool) {
	if serverNames == nil {
		serverNames = getServerNames(rConf.GroupBind)

		if serverNames == nil {
			log.ConfigLogger().Errorf("[gateway][router][AddOrUpdateRouter] no server name found, add router fail, [%v]", rConf)
			return errors.Errorf("no server name found, add router fail"), false
		}
	}

	log.ConfigLogger().Infof("[gateway][router][AddOrUpdateRouter] start add router %s", rConf.Name)

	if router, err := m.routerFactory(rConf); err != nil {
		// todo err
	} else {
		m.rMux.Lock()
		m.sMux.Lock()
		defer m.sMux.Unlock()
		defer m.rMux.Unlock()

		gr := m.routerMap[rConf.GroupBind]
		if gr == nil {
			gr = map[string]types.Router{}
			m.routerMap[rConf.GroupBind] = gr
		}
		gr[router.Conf().Name] = router

		for _, serverName := range serverNames {
			sr := m.serverRouterMap[serverName]
			if sr == nil {
				sr = map[string]types.Router{}
				m.serverRouterMap[serverName] = sr
			}
			sr[router.Conf().Name] = router
		}
	}

	return nil, true
}

func getServerNames(groupName string) []string {

	if rg := config.StoreInstance().Get(api.ROUTER_GROUP, groupName); rg != nil {
		return rg.(*v1.RouterGroup).Gateways
	}

	return nil

}

func defaultRouterFactory(rConf *v1.Router) (types.Router, error) {
	s := service.GetServiceManagerInstance().GetService(rConf.Proxy.Service, rConf.Proxy.Protocol)

	if s == nil {
		return nil, errors.Errorf("no service found for router")
	}

	p := pipeline.GetPipelineManagerInstance().BuildPipeline(rConf, s.Conf())
	m := GetMatcherFactory().CreateMatchers(rConf.Matches)

	return &router{
		pipeline: p.Copy(),
		service:  s,
		matchers: m,
		conf:     *rConf,
	}, nil
}

func defaultRouterMatcher(listenerName string, headers mosn.HeaderMap) types.Router {
	if sConf := config.StoreInstance().Get(api.GATEWAY, listenerName); sConf != nil {
		serverName := sConf.(*v1.Gateway).Name

		routers := GetRouterManager().GetServerRouters(serverName)

		if routers == nil {
			return nil
		}
		return match(routers, headers)
	}
	return nil
}

// todo 重构
func match(routers map[string]types.Router, headers mosn.HeaderMap) types.Router {
	var maxPrefixRouter types.Router

	for _, router := range routers {
		if ok := router.Match(headers); ok {
			if router.Conf().Matches.Prefix != "" {
				// prefix
				maxPrefixRouter = comparePrefixLen(maxPrefixRouter, router)
			} else if router.Conf().Matches.Path != "" {
				// path
				return router
			} else if router.Conf().Matches.Regex {
				// todo regex
			} else {
				// todo only header
				return router
			}
		}
	}

	return maxPrefixRouter
}

func comparePrefixLen(res, target types.Router) types.Router {
	if res == nil {
		return target
	}

	var splitStr = "/"
	r := strings.Split(res.Conf().Matches.Prefix, splitStr)
	t := strings.Split(target.Conf().Matches.Prefix, splitStr)

	if len(r) > len(t) {
		return res
	}

	return target
}

func newRouterManager() *Manager {
	return &Manager{
		rMux:            sync.RWMutex{},
		routerMap:       map[string]map[string]types.Router{},
		sMux:            sync.RWMutex{},
		serverRouterMap: map[string]map[string]types.Router{},
		routerMatcher:   defaultRouterMatcher,
		routerFactory:   defaultRouterFactory,
	}
}
