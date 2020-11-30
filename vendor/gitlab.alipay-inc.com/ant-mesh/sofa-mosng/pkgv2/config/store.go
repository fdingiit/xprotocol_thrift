package config

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api"
	v1 "gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/event"
	"reflect"
	"sync"
)

var (
	storeInstance Store = &store{
		gMux:         sync.RWMutex{},
		g:            &v1.GwConfig{},
		sMux:         sync.RWMutex{},
		serverMap:    make(map[string]*v1.Gateway),
		svcMux:       sync.RWMutex{},
		svcMap:       make(map[string]*v1.GatewayService),
		fcMux:        sync.RWMutex{},
		filterChains: make(map[string]*v1.FilterChain),
		gfMux:        sync.RWMutex{},
		rgMux:        sync.RWMutex{},
		routerGroup:  make(map[string]*v1.RouterGroup),
		cfgMux:       sync.RWMutex{},
		msMap:        make(map[string]*v1.Metadata),
	}
)

func StoreInstance() Store {
	return storeInstance
}

// todo 重构
type store struct {
	gMux         sync.RWMutex
	g            *v1.GwConfig
	sMux         sync.RWMutex
	serverMap    map[string]*v1.Gateway
	svcMux       sync.RWMutex
	svcMap       map[string]*v1.GatewayService
	fcMux        sync.RWMutex
	filterChains map[string]*v1.FilterChain
	gfMux        sync.RWMutex
	globalFilter *v1.GlobalFilter
	rgMux        sync.RWMutex
	routerGroup  map[string]*v1.RouterGroup
	rMux         sync.RWMutex
	routers      map[string]*v1.Router // todo groupName|routerName = router
	cfgMux       sync.RWMutex
	msMap        map[string]*v1.Metadata
}

func (s *store) Diff(o api.Object) []event.DifferEvent {
	events := []event.DifferEvent{{
		Object:    o,
		EventType: event.Update,
	}}

	switch o.Type() {
	case api.GW_CONFIG:
		return s.g.Diff(o.(*v1.GwConfig))
	case api.GATEWAY:
		ns := o.(*v1.Gateway)
		if os := s.serverMap[ns.Name]; os != nil {
			return os.Diff(ns)
		}
	case api.ROUTER_GROUP:
		if e := s.diffRouterGroup(o); e != nil {
			return e
		}
	case api.ROUTER:
		nr := o.(*v1.Router)
		if or := s.routers[nr.GroupBind+"|"+nr.Name]; or != nil {
			return or.Diff(nr)
		}
	case api.FILTER_CHAIN:
		nc := o.(*v1.FilterChain)
		if oc := s.filterChains[nc.ChainName]; oc != nil {
			return oc.Diff(nc)
		}
	case api.GLOBAL_FILTER:
		ng := o.(v1.GlobalFilter)
		if og := s.globalFilter; og != nil {
			return og.Diff(ng)
		}
	case api.SERVICE:
		ns := o.(*v1.GatewayService)
		if os := s.svcMap[ns.Name]; os != nil {
			return os.Diff(ns)
		}
	case api.METADATA:
		nc := o.(*v1.Metadata)
		if oc := s.msMap[nc.Key]; oc != nil {
			return oc.Diff(nc)
		}
	}
	return events
}

func (s *store) diffRouterGroup(o api.Object) []event.DifferEvent {
	nrg := o.(*v1.RouterGroup)
	var des []event.DifferEvent
	if org := s.routerGroup[nrg.GroupName]; org != nil {

		for _, nr := range nrg.Routers {
			// 原来就有
			if or := s.routers[nrg.GroupName+"|"+nr.Name]; or != nil {
				differEvents := or.Diff(nr)
				des = append(des, differEvents...)
				continue
			}
			// 原来没有
			des = append(des, event.DifferEvent{Object: nr, EventType: event.Add})
		}

		// 原来有，现在没有
		ina := routerOnlyInA(org.Routers, nrg.Routers)
		for _, dr := range ina {
			des = append(des, event.DifferEvent{Object: dr, EventType: event.Delete})
		}
	}
	return des
}

func (s *store) Store(o api.Object) {
	switch o.Type() {
	case api.GW_CONFIG:
		s.SaveGateway(o.(v1.GwConfig))
	case api.GATEWAY:
		s.SaveServer(o.(*v1.Gateway))
	case api.ROUTER_GROUP:
		s.SaveRouterGroup(o.(*v1.RouterGroup))
	case api.FILTER_CHAIN:
		s.SaveFilterChain(o.(*v1.FilterChain))
	case api.GLOBAL_FILTER:
		s.SaveGlobalFilter(o.(v1.GlobalFilter))
	case api.SERVICE:
		s.SaveService(o.(*v1.GatewayService))
	case api.METADATA:
		s.SaveConfig(o.(*v1.Metadata))
	}
}

func (s *store) Get(t api.Type, name string) api.Object {
	switch t {
	case api.GW_CONFIG:
		return objOrNil(s.g)
	case api.GATEWAY:
		return objOrNil(s.serverMap[name])
	case api.ROUTER_GROUP:
		return objOrNil(s.routerGroup[name])
	case api.FILTER_CHAIN:
		return objOrNil(s.filterChains[name])
	case api.GLOBAL_FILTER:
		return objOrNil(s.globalFilter)
	case api.SERVICE:
		return objOrNil(s.serverMap[name])
	case api.METADATA:
		return objOrNil(s.msMap[name])
	}

	return nil
}

func objOrNil(i interface{}) api.Object {
	if IsNil(i) {
		return nil
	}
	return i.(api.Object)
}

func IsNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}

func (s *store) GetAll(t api.Type) interface{} {
	switch t {
	case api.GW_CONFIG:
		return s.g
	case api.GATEWAY:
		return s.serverMap
	case api.ROUTER_GROUP:
		return s.routerGroup
	case api.FILTER_CHAIN:
		return s.filterChains
	case api.GLOBAL_FILTER:
		return s.globalFilter
	case api.SERVICE:
		return s.serverMap
	case api.METADATA:
		return s.msMap
	}

	return nil
}

func (s *store) SaveGateway(g v1.GwConfig) {
	for _, rg := range g.RouterGroups {
		s.SaveRouterGroup(rg)
	}

	for _, cfg := range g.GatewayMetadatas {
		s.SaveConfig(cfg)
	}

	for _, fc := range g.FilterChains {
		s.SaveFilterChain(fc)
	}

	s.SaveGlobalFilter(g.GlobalFilters)

	for _, svc := range g.GatewayServices {
		s.SaveService(svc)
	}

	for _, ser := range g.Gateways {
		s.SaveServer(ser)
	}
}

func (s *store) SaveServer(ser *v1.Gateway) {
	s.sMux.Lock()
	s.gMux.Lock()

	s.serverMap[ser.Name] = ser
	s.g.Gateways = append(s.g.Gateways, ser)

	s.gMux.Unlock()
	s.sMux.Unlock()
}

func (s *store) SaveRouterGroup(rg *v1.RouterGroup) {
	s.rgMux.Lock()
	s.gMux.Lock()

	s.routerGroup[rg.GroupName] = rg
	s.g.RouterGroups = append(s.g.RouterGroups, rg)

	s.gMux.Unlock()
	s.rgMux.Unlock()
}

func (s *store) SaveFilterChain(fc *v1.FilterChain) {
	s.fcMux.Lock()
	s.gMux.Lock()

	s.filterChains[fc.ChainName] = fc
	s.g.FilterChains = append(s.g.FilterChains, fc)

	s.gMux.Unlock()
	s.fcMux.Unlock()
}

func (s *store) SaveGlobalFilter(gf v1.GlobalFilter) {
	s.gfMux.Lock()
	s.gMux.Lock()

	s.globalFilter = &gf
	s.g.GlobalFilters = gf

	s.gMux.Unlock()
	s.gfMux.Unlock()
}

func (s *store) SaveService(svc *v1.GatewayService) {
	s.svcMux.Lock()
	s.gMux.Lock()

	s.svcMap[svc.Name] = svc
	s.g.GatewayServices = append(s.g.GatewayServices, svc)

	s.gMux.Unlock()
	s.svcMux.Unlock()
}

func (s *store) SaveConfig(cfg *v1.Metadata) {
	s.cfgMux.Lock()
	s.gMux.Lock()

	s.msMap[cfg.Key] = cfg
	s.g.GatewayMetadatas = append(s.g.GatewayMetadatas, cfg)

	s.gMux.Unlock()
	s.cfgMux.Unlock()
}

func routerOnlyInA(a, b []*v1.Router) (ina []*v1.Router) {
	for _, v := range a {
		if !InSliceIface(v, b) {
			ina = append(ina, v)
		}
	}
	return
}

func InSliceIface(v *v1.Router, sl []*v1.Router) bool {
	for _, vv := range sl {
		if compare(vv, v) {
			return true
		}
	}
	return false
}

func compare(a, b *v1.Router) bool {
	return a.Name == b.Name
}
