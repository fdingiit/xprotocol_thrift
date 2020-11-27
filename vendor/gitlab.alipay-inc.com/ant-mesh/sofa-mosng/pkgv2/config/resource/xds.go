package resource

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/common/log"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/dispatcher"
)

var xds = &xdsResource{
	baseResource{
		dispatcher: dispatcher.DispatcherInstance(),
	},
}

func XDSResource() *xdsResource {
	return xds
}

const (
	MosngConfig       = "type.googleapis.com/networking.istio.io.v1alpha3.GatewayConfig"
	MosngGateway      = "type.alipayapis.com/mosn.api.v1.gateway.MosngGateway"
	MosngRouterGroup  = "type.alipayapis.com/mosn.api.v1.gateway.MosngRouterGroup"
	MosngRouter       = "type.alipayapis.com/mosn.api.v1.gateway.MosngRouter"
	MosngService      = "type.alipayapis.com/mosn.api.v1.gateway.MosngService"
	MosngFilterChain  = "type.alipayapis.com/mosn.api.v1.gateway.MosngFilterChain"
	MosngGlobalFilter = "type.alipayapis.com/mosn.api.v1.gateway.MosngGlobalFilter"
	MosngMetadata     = "type.alipayapis.com/mosn.api.v1.gateway.metadata"
)

type xdsResource struct {
	baseResource
}

func (b *xdsResource) Start() {
	log.StartLogger().Infof("[mosng][resource][start] xdsResource start")
	//xdsv2.RegisterTypeURLHandleFunc(MosngConfig, HandleGwConfig)
	//xdsv2.RegisterTypeURLHandleFunc(MosngGateway, HandleGateway)
	//xdsv2.RegisterTypeURLHandleFunc(MosngRouterGroup, HandleRouterGroup)
	//xdsv2.RegisterTypeURLHandleFunc(MosngRouter, HandleRouter)
	//xdsv2.RegisterTypeURLHandleFunc(MosngService, HandleService)
	//xdsv2.RegisterTypeURLHandleFunc(MosngFilterChain, HandleFilterChain)
	//xdsv2.RegisterTypeURLHandleFunc(MosngGlobalFilter, HandleGlobalFilter)
	//xdsv2.RegisterTypeURLHandleFunc(MosngMetadata, HandleMetadata)
}

//
//func HandleGwConfig(client *xdsv2.ADSClient, response *envoy_api_v2.DiscoveryResponse) {
//
//	defer func() {
//		if err := recover(); err != nil {
//			log.ConfigLogger().Errorf("[xds][HandleGwConfig][recover], err msg %+v", err)
//		}
//	}()
//
//	log.ConfigLogger().Infof("[xds][HandleGwConfig] receive config:  %v", response.Resources)
//
//	gwcfg := &v1alpha3.GatewayConfig{}
//	var err error
//	for _, res := range response.Resources {
//		if err = gwcfg.Unmarshal(res.GetValue()); err != nil {
//			log.ConfigLogger().Errorf("[xds][HandleGwConfig] gwcfg.Unmarshal err %v", err)
//		}
//
//		cfg := &v1.GwConfig{}
//
//		cfg.Gateways = convertGateways(gwcfg.Gateways)
//
//		cfg.RouterGroups = convertRouterGroups(gwcfg.RouterGroups)
//
//		cfg.Routers = convertRouters(gwcfg.Routers)
//
//		cfg.FilterChains = convertFilterChains(gwcfg.FilterChains)
//
//		cfg.GlobalFilters = convertGlobalFilter(gwcfg.GlobalFilter)
//
//		cfg.GatewayServices = convertServices(gwcfg.GatewayServices)
//
//		cfg.GatewayMetadatas = convertMetaDatas(gwcfg.GatewayMetadatas)
//
//		bytes, _ := json.Marshal(cfg)
//		log.ConfigLogger().Infof("[xds][HandleGwConfig] convert config %s", bytes)
//
//		if err, ok := XDSResource().dispatcher.Dispatch(event.Add, cfg); !ok {
//			log.ConfigLogger().Errorf("[xds][HandleGwConfig] Dispatch err %v", err)
//		}
//	}
//
//}
//
//func HandleGateway(client *xdsv2.ADSClient, response *envoy_api_v2.DiscoveryResponse) {
//	log.ConfigLogger().Infof("[xds][HandleGateway] receive config:  %v", response.Resources)
//
//	gwpb := &v1alpha3.Gateway{}
//	var err error
//	for _, res := range response.Resources {
//		if err = gwpb.Unmarshal(res.GetValue()); err != nil {
//			log.ConfigLogger().Errorf("[xds][HandleGateway] gwcfg.Unmarshal err %v", err)
//		}
//
//		gwv1 := convertGateway(gwpb)
//
//		for _, s := range gwv1 {
//			if err, ok := XDSResource().dispatcher.Dispatch(event.Add, s); !ok {
//				log.ConfigLogger().Errorf("[xds][HandleGateway] Dispatch err %v", err)
//			}
//		}
//
//	}
//
//}
//
//func HandleRouterGroup(client *xdsv2.ADSClient, response *envoy_api_v2.DiscoveryResponse) {
//	log.ConfigLogger().Infof("[xds][HandleRouterGroup] receive config:  %v", response.Resources)
//
//	rgpb := &v1alpha3.RouterGroup{}
//	var err error
//	for _, res := range response.Resources {
//		if err = rgpb.Unmarshal(res.GetValue()); err != nil {
//			log.ConfigLogger().Errorf("[xds][HandleRouterGroup] gwcfg.Unmarshal err %v", err)
//			return
//		}
//
//		rgv1 := convertRouterGroup(rgpb)
//
//		if err, ok := XDSResource().dispatcher.Dispatch(event.Add, rgv1); !ok {
//			log.ConfigLogger().Errorf("[xds][HandleRouterGroup] Dispatch err %v", err)
//		}
//	}
//
//}
//
//func HandleRouter(client *xdsv2.ADSClient, response *envoy_api_v2.DiscoveryResponse) {
//	log.ConfigLogger().Infof("[xds][HandleRouter] receive config:  %v", response.Resources)
//
//	rpb := &v1alpha3.GatewayRouter{}
//	var err error
//	for _, res := range response.Resources {
//		if err = rpb.Unmarshal(res.GetValue()); err != nil {
//			log.ConfigLogger().Errorf("[xds][HandleRouter] gwcfg.Unmarshal err %v", err)
//			return
//		}
//
//		rv1 := convertRouter(rpb)
//
//		if err, ok := XDSResource().dispatcher.Dispatch(event.Add, rv1); !ok {
//			log.ConfigLogger().Errorf("[xds][HandleRouter] Dispatch err %v", err)
//		}
//	}
//
//}
//
//func HandleService(client *xdsv2.ADSClient, response *envoy_api_v2.DiscoveryResponse) {
//	log.ConfigLogger().Infof("[xds][HandleService] receive config:  %v", response.Resources)
//
//	spb := &v1alpha3.GatewayService{}
//	var err error
//	for _, res := range response.Resources {
//		if err = spb.Unmarshal(res.GetValue()); err != nil {
//			log.ConfigLogger().Errorf("[xds][HandleService] gwcfg.Unmarshal err %v", err)
//			return
//		}
//
//		sv1 := convertService(spb)
//
//		if err, ok := XDSResource().dispatcher.Dispatch(event.Add, sv1); !ok {
//			log.ConfigLogger().Errorf("[xds][HandleService] Dispatch err %v", err)
//		}
//	}
//
//}
//
//func HandleGlobalFilter(client *xdsv2.ADSClient, response *envoy_api_v2.DiscoveryResponse) {
//	log.ConfigLogger().Infof("[xds][HandleGlobalFilter] receive config:  %v", response.Resources)
//
//	gfpb := &v1alpha3.GlobalFilter{}
//	var err error
//	for _, res := range response.Resources {
//		if err = gfpb.Unmarshal(res.GetValue()); err != nil {
//			log.ConfigLogger().Errorf("[xds][HandleGlobalFilter] gwcfg.Unmarshal err %v", err)
//			return
//		}
//
//		gfv1 := convertGlobalFilter(gfpb)
//
//		if err, ok := XDSResource().dispatcher.Dispatch(event.Add, gfv1); !ok {
//			log.ConfigLogger().Errorf("[xds][HandleGlobalFilter] Dispatch err %v", err)
//		}
//	}
//
//}
//
//func HandleFilterChain(client *xdsv2.ADSClient, response *envoy_api_v2.DiscoveryResponse) {
//	log.ConfigLogger().Infof("[xds][HandleFilterChain] receive config:  %v", response.Resources)
//
//	fcpb := &v1alpha3.FilterChain{}
//	var err error
//	for _, res := range response.Resources {
//		if err = fcpb.Unmarshal(res.GetValue()); err != nil {
//			log.ConfigLogger().Errorf("[xds][HandleFilterChain] gwcfg.Unmarshal err %v", err)
//			return
//		}
//
//		fcv1 := convertFilterChain(fcpb)
//
//		if err, ok := XDSResource().dispatcher.Dispatch(event.Add, fcv1); !ok {
//			log.ConfigLogger().Errorf("[xds][HandleFilterChain] Dispatch err %v", err)
//		}
//	}
//}
//
//func HandleMetadata(client *xdsv2.ADSClient, response *envoy_api_v2.DiscoveryResponse) {
//	log.ConfigLogger().Infof("[xds][HandleMetadata] receive config:  %v", response.Resources)
//
//	metapb := &v1alpha3.GatewayMetadata{}
//	var err error
//	for _, res := range response.Resources {
//		if err = metapb.Unmarshal(res.GetValue()); err != nil {
//			log.ConfigLogger().Errorf("[xds][HandleMetadata] gwcfg.Unmarshal err %v", err)
//			return
//		}
//
//		metav1 := convertMetaData(metapb)
//
//		if err, ok := XDSResource().dispatcher.Dispatch(event.Add, metav1); !ok {
//			log.ConfigLogger().Errorf("[xds][HandleMetadata] Dispatch err %v", err)
//		}
//	}
//}
//
//func convertMetaDatas(metadata []*v1alpha3.GatewayMetadata) (vms []*v1.Metadata) {
//	if metadata == nil {
//		return
//	}
//	for _, m := range metadata {
//		vms = append(vms, convertMetaData(m))
//	}
//
//	return
//}
//
//func convertMetaData(m *v1alpha3.GatewayMetadata) *v1.Metadata {
//	if m == nil {
//		return nil
//	}
//	vm := &v1.Metadata{
//		Key:   m.Key,
//		Value: m.Value,
//	}
//
//	return vm
//}
//
//func convertServices(services []*v1alpha3.GatewayService) (vss []*v1.GatewayService) {
//	if services == nil {
//		return
//	}
//	for _, s := range services {
//		vss = append(vss, convertService(s))
//	}
//	return
//}
//
//func convertService(s *v1alpha3.GatewayService) *v1.GatewayService {
//	if s == nil {
//		return nil
//	}
//	vs := &v1.GatewayService{
//		Name:         s.Name,
//		Protocol:     s.Protocol,
//		Hosts:        convertHost(s.Hosts),
//		LbType:       convertLb(s.LbType),
//		Filters:      convertFilters(s.Filters),
//		FilterChains: s.FilterChains,
//	}
//	return vs
//}
//
//func convertLb(lbType v1alpha3.LbType) v2.LbType {
//	if lbType == 1 {
//		return v2.LB_RANDOM
//	} else {
//		return v2.LB_ROUNDROBIN
//	}
//}
//
//func convertHost(configs []*v1alpha3.HostConfig) (vhs []*v1.HostConfig) {
//	if configs == nil {
//		return
//	}
//	for _, c := range configs {
//		vh := &v1.HostConfig{
//			Address:    c.Address,
//			Hostname:   c.Hostname,
//			Weight:     c.Weight,
//			TLSDisable: c.TLSDisable,
//		}
//
//		vhs = append(vhs, vh)
//	}
//
//	return
//}
//
//func convertGlobalFilter(filter *v1alpha3.GlobalFilter) v1.GlobalFilter {
//	gf := v1.GlobalFilter{}
//
//	if filter != nil {
//		gf.Filters = convertFilters(filter.Filters)
//
//	}
//	return gf
//}
//
//func convertFilterChains(chains []*v1alpha3.FilterChain) (vfs []*v1.FilterChain) {
//	if chains == nil {
//		return
//	}
//	for _, c := range chains {
//		vfs = append(vfs, convertFilterChain(c))
//	}
//
//	return
//}
//
//func convertFilterChain(c *v1alpha3.FilterChain) *v1.FilterChain {
//	if c == nil {
//		return nil
//	}
//	vf := &v1.FilterChain{
//		ChainName: c.Name,
//		Filters:   convertFilters(c.Filters),
//	}
//
//	return vf
//}
//
//func convertRouterGroups(groups []*v1alpha3.RouterGroup) (vrs []*v1.RouterGroup) {
//	if groups == nil {
//		return
//	}
//	for _, g := range groups {
//		vrs = append(vrs, convertRouterGroup(g))
//	}
//
//	return
//}
//
//func convertRouterGroup(g *v1alpha3.RouterGroup) *v1.RouterGroup {
//
//	v1g := &v1.RouterGroup{}
//	if g == nil {
//		return v1g
//	}
//	v1g.Gateways = g.Gateways
//	v1g.GroupName = g.Name
//	v1g.Routers = convertRouters(g.Routers)
//
//	return v1g
//}
//
//func convertRouters(routers []*v1alpha3.GatewayRouter) (rs []*v1.Router) {
//	if routers == nil {
//		return
//	}
//	for _, r := range routers {
//		rs = append(rs, convertRouter(r))
//	}
//	return
//}
//
//func convertRouter(r *v1alpha3.GatewayRouter) *v1.Router {
//	v1r := &v1.Router{}
//	if r == nil {
//		return v1r
//	}
//	v1r.Name = r.Name
//	v1r.Metadata = convertMeta(r.Meta)
//	v1r.Proxy = convertProxy(r.GetProxy())
//	v1r.FilterChains = r.FilterChains
//	v1r.Filters = convertFilters(r.Filters)
//	v1r.GroupBind = r.GroupBind
//	v1r.Matches = convertMatches(r.Matches)
//	v1r.Status = convertStatus(r.Status)
//	v1r.Timeout = uint64(r.Timeout)
//
//	return v1r
//}
//
//func convertStatus(s string) v1.Status {
//	if s == string(v1.OPEN) {
//		return v1.OPEN
//	} else {
//		return v1.CLOSE
//	}
//}
//
//func convertMatches(match *v1alpha3.RouterMatch) *v1.RouterMatch {
//	v1m := &v1.RouterMatch{}
//
//	if match == nil {
//		return v1m
//	}
//	v1m.Metadata = convertMeta(match.Meta)
//	v1m.Path = match.GetPath()
//	v1m.Headers = convertHeader(match.HeaderMatcher)
//	v1m.QueryStrings = convertHeader(match.QueryStringMatcher)
//	v1m.Prefix = match.GetPrefix()
//	v1m.Regex = match.GetRegex()
//	return v1m
//
//}
//
//func convertHeader(matchers []*v1alpha3.StringMatcher) (vms []*v1.ValueMatcher) {
//	if matchers == nil {
//		return
//	}
//	for _, m := range matchers {
//		vm := &v1.ValueMatcher{
//			Regex: m.Regex,
//			Name:  m.Name,
//			Value: m.Value,
//		}
//		vms = append(vms, vm)
//	}
//	return
//}
//
//func convertFilters(filters []*v1alpha3.Filter) (vfs []*v1.Filter) {
//	if filters == nil {
//		return
//	}
//	for _, f := range filters {
//		if p, err := strconv.ParseInt(f.Priority, 10, 64); err == nil {
//			vf := &v1.Filter{
//				Name:     f.Name,
//				Priority: p,
//				Metadata: convertMeta(f.Meta),
//			}
//			vfs = append(vfs, vf)
//		}
//	}
//	return
//}
//
//func convertProxy(proxy *v1alpha3.Proxy) (v1p *v1.Proxy) {
//	v1p = &v1.Proxy{}
//	if proxy == nil {
//		return
//	}
//
//	v1p.Metadata = convertMeta(proxy.Meta)
//	v1p.Service = proxy.Service
//	v1p.Method = proxy.Method
//	v1p.Interface = proxy.Interface
//	v1p.Protocol = proxy.Protocol
//	v1p.Timeout = uint64(proxy.Timeout)
//	return
//}
//
//func convertMeta(meta map[string]string) (m map[string]interface{}) {
//	m = map[string]interface{}{}
//
//	if meta == nil {
//		return
//	}
//
//	for k, v := range meta {
//		m[k] = v
//	}
//
//	return
//}
//
//func convertGateways(gateways []*v1alpha3.Gateway) (vgs []*v1.Gateway) {
//	if gateways == nil {
//		return
//	}
//	for _, g := range gateways {
//		for _, s := range g.Servers {
//			gateway := &v1.Gateway{}
//			gateway.Name = fmt.Sprintf("%s/%s/%s", g.GatewayExt.Namespace, g.GatewayExt.Name, s.Port.Name)
//			gateway.AddrConfig = fmt.Sprintf("%s:%d", s.Port.PortExt.Host, s.Port.Number)
//			gateway.BindToPort = s.Port.PortExt.BindPort
//			gateway.DownstreamProtocol = s.Port.PortExt.DownstreamProtocol
//			gateway.UpstreamProtocol = s.Port.PortExt.UpstreamProtocol
//			// gateway.ListenerType = s.Port.PortExt.m
//			vgs = append(vgs, gateway)
//		}
//	}
//	return
//}
//
//func convertGateway(g *v1alpha3.Gateway) (vgs []*v1.Gateway) {
//	if g == nil {
//		return
//	}
//	for _, s := range g.Servers {
//		gateway := &v1.Gateway{}
//		gateway.Name = fmt.Sprintf("%s/%s/%s", g.GatewayExt.Namespace, g.GatewayExt.Name, s.Port.Name)
//		gateway.AddrConfig = fmt.Sprintf("%s:%s", s.Port.PortExt.Host, s.Port)
//		gateway.BindToPort = s.Port.PortExt.BindPort
//		gateway.DownstreamProtocol = s.Port.PortExt.DownstreamProtocol
//		gateway.UpstreamProtocol = s.Port.PortExt.UpstreamProtocol
//		vgs = append(vgs, gateway)
//	}
//
//	return
//}
