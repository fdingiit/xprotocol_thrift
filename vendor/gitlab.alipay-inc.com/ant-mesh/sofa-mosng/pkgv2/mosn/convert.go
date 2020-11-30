package mosn

import (
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/common/constants"
	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/common/log"
	v1 "gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api/v1"
	"mosn.io/api"
	v2 "mosn.io/mosn/pkg/config/v2"
	"net"
	"os/user"
	"time"

	"encoding/json"
)

func convertRouterGroupConfig(router *v1.RouterGroup) *v2.RouterConfiguration {
	return buildRouterGroup(router)
}

func convertListenerConfig(server *v1.Gateway) *v2.Listener {
	var logFolder string
	if usr, err := user.Current(); err != nil {
		logFolder = "/home/admin/logs/sofa-mosng"
	} else {
		logFolder = usr.HomeDir + "/logs/sofa-mosng"
	}
	// todo convert
	listenerConfig := &v2.Listener{
		ListenerConfig: v2.ListenerConfig{
			//ConnectionIdleTimeout: &v2.DurationConfig{Duration: server.ConnectionIdleTimeout},
			Type:       server.ListenerType,
			Name:       server.Name,
			BindToPort: server.BindToPort,
			Inspector:  server.Inspector,
			AccessLogs: []v2.AccessLog{{Path: logFolder + "/access.log", Format: "%start_time% %request_received_duration% %response_received_duration% %protocol%"}},
		},
		Addr:                    convertAddress(server.AddrConfig),
		PerConnBufferLimitBytes: 65536,
	}

	// virtual listener need none filters
	if listenerConfig.Name == "virtual" {
		return nil
	}

	// network filter: proxy & connection_manager
	listenerConfig.FilterChains = convertNetWorkFilter(server)

	listenerConfig.StreamFilters = convertStreamFilter()

	return listenerConfig
}

func convertClustersConfig(service *v1.GatewayService) *v2.Cluster {
	cluster := &v2.Cluster{
		Name:                 service.Name,
		ClusterType:          service.ClusterType,
		LbType:               service.LbType,
		SubType:              service.SubType,
		LBSubSetConfig:       convertLbSubType(service.LBSubSetConfig),
		MaxRequestPerConn:    service.MaxRequestPerConn,
		ConnBufferLimitBytes: service.ConnBufferLimitBytes,
		HealthCheck:          convertHealthCheckCfg(service.HealthCheck),
		CirBreThresholds:     convertCir(service.CircuitBreakers),
		Hosts:                buildMosnHosts(service.Hosts),
		Spec:                 convertSpec(service.Spec),
		TLS:                  convertTlCfg(service.TLS),
	}
	return cluster
}

func convertLbSubType(config v1.LBSubsetConfig) v2.LBSubsetConfig {
	return v2.LBSubsetConfig{
		FallBackPolicy:  config.FallBackPolicy,
		DefaultSubset:   config.DefaultSubset,
		SubsetSelectors: config.SubsetSelectors,
	}
}

func convertSpec(info v1.ClusterSpecInfo) v2.ClusterSpecInfo {
	var subscribes []v2.SubscribeSpec
	for _, sub := range info.Subscribes {
		subscribes = append(subscribes, v2.SubscribeSpec{
			Subscriber:  sub.Subscriber,
			ServiceName: sub.ServiceName,
		})
	}
	return v2.ClusterSpecInfo{
		Subscribes: subscribes,
	}
}

func convertCir(breakers v1.CircuitBreakers) v2.CircuitBreakers {
	var thresholds []v2.Thresholds
	for _, b := range breakers.Thresholds {
		thresholds = append(thresholds, v2.Thresholds{
			MaxConnections:     b.MaxConnections,
			MaxPendingRequests: b.MaxPendingRequests,
			MaxRequests:        b.MaxRequests,
			MaxRetries:         b.MaxRetries,
		})
	}
	return v2.CircuitBreakers{
		Thresholds: thresholds,
	}
}

func buildMosnHosts(hostConfigs []*v1.HostConfig) []v2.Host {
	hosts := make([]v2.Host, 0, len(hostConfigs))

	for _, hostConfig := range hostConfigs {
		hosts = append(hosts, v2.Host{
			HostConfig: convertHost(hostConfig),
		})
	}
	return hosts
}

func convertHost(config *v1.HostConfig) v2.HostConfig {
	return v2.HostConfig{
		Address:    config.Address,
		Hostname:   config.Hostname,
		Weight:     config.Weight,
		TLSDisable: config.TLSDisable,
	}
}

func convertAddress(addr string) net.Addr {

	if addr == "" {
		return nil
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.ConfigLogger().Errorf("Invalid address: %v", err)
		return nil
	}
	return tcpAddr
}

func convertNetWorkFilter(server *v1.Gateway) []v2.FilterChain {
	filters := make([]v2.Filter, 0, 2)
	filterChains := make([]v2.FilterChain, 0, 1)

	// create proxy filter: proxy tcpproxy WebSocketProxy
	if err, proxyFilter := buildProxyFilter(server); err != nil {
		// todo log
	} else {
		filters = append(filters, *proxyFilter)
	}

	// create connection_manager filter
	if err, proxyFilter := buildConnectionManagerFilter(server); err != nil {
		// todo log
	} else {
		filters = append(filters, proxyFilter)
	}

	// todo tlsConfig := convertTLS(xdsFilterChain.GetTlsContext())
	cfg := convertTlCfg(server.TLSConfig)

	var cfgs []v2.TLSConfig
	for _, tls := range server.TLSConfigs {
		cfgs = append(cfgs, convertTlCfg(tls))
	}

	filterChain := v2.FilterChain{
		FilterChainConfig: v2.FilterChainConfig{
			Filters:    filters,
			TLSConfig:  &cfg,
			TLSConfigs: cfgs,
		},
		TLSContexts: convertTls(cfg, cfgs),
	}

	filterChains = append(filterChains, filterChain)
	return filterChains
}

func convertTls(cfg v2.TLSConfig, cfgs []v2.TLSConfig) []v2.TLSConfig {
	if &cfg != nil && len(cfgs) > 0 {
		// todo err
	}

	var TLSContexts []v2.TLSConfig
	if len(cfgs) > 0 {
		copy(TLSContexts, cfgs)
	} else { // no tls_context_set, use tls_context
		if &cfg == nil { // no tls_context, generate a default one
			TLSContexts = append(TLSContexts, v2.TLSConfig{})
		} else { // use tls_context
			TLSContexts = append(TLSContexts, cfg)
		}
	}

	return TLSContexts
}

func convertStreamFilter() []v2.Filter {
	filter := v2.Filter{
		Type: "mosng",
	}

	return []v2.Filter{filter}
}

func buildProxyFilter(server *v1.Gateway) (error, *v2.Filter) {
	// todo rpc protocol
	proxyConfig := v2.Proxy{
		Name:               server.Name,
		DownstreamProtocol: server.DownstreamProtocol,
		UpstreamProtocol:   server.UpstreamProtocol,
		RouterConfigName:   server.Name,
		ValidateClusters:   false,
	}

	filter := &v2.Filter{
		Type:   "proxy",
		Config: toMap(proxyConfig),
	}

	return nil, filter
}

func buildConnectionManagerFilter(server *v1.Gateway) (error, v2.Filter) {

	configuration := buildRouterForServer(server)

	filter := v2.Filter{
		Type:   "connection_manager",
		Config: toMap(configuration),
	}

	return nil, filter
}

func buildRouterGroup(routerGroup *v1.RouterGroup) *v2.RouterConfiguration {
	configuration := &v2.RouterConfiguration{
		RouterConfigurationConfig: v2.RouterConfigurationConfig{
			// todo support more
			RouterConfigName: routerGroup.Gateways[0],
			// todo header setter
		},

		VirtualHosts: convertVirtualHosts(routerGroup.Routers),
	}

	return configuration
}

func buildRouterForServer(server *v1.Gateway) *v2.RouterConfiguration {
	//routerGroup := resource.StoreInstance().Get(api.ROUTER_GROUP, server.GroupBind)
	//if routerGroup == nil {
	//	// todo err
	//}
	configuration := &v2.RouterConfiguration{
		RouterConfigurationConfig: v2.RouterConfigurationConfig{
			RouterConfigName: server.Name,
			// todo header setter
		},

		//VirtualHosts: convertVirtualHosts(routerGroup.(v1.RouterGroup).Routers),
	}
	return configuration
}

func convertVirtualHosts(routers []*v1.Router) []*v2.VirtualHost {
	hosts := make([]*v2.VirtualHost, 0, 1) // todo more vh
	mosnRouters := make([]v2.Router, 0, len(routers))

	host := v2.VirtualHost{
		Domains: []string{"*"}, // todo
	}

	for _, router := range routers {

		mosnRouter := v2.Router{
			RouterConfig: v2.RouterConfig{
				Match: buildMosnMatch(router.Proxy),
				Route: buildMosnRoute(router),
			},
		}
		mosnRouters = append(mosnRouters, mosnRouter)
	}

	host.Routers = mosnRouters
	hosts = append(hosts, &host)

	return hosts
}

func buildMosnMatch(proxy *v1.Proxy) v2.RouterMatch {
	return v2.RouterMatch{
		Prefix: "/",
		Headers: []v2.HeaderMatcher{
			{
				Name:  constants.MosngServiceHeader,
				Value: proxy.Service,
			},
		},
	}
}

func buildMosnRoute(router *v1.Router) v2.RouteAction {
	action := v2.RouteAction{
		RouterActionConfig: v2.RouterActionConfig{
			// todo
			ClusterName: router.Proxy.Service,
		},
		Timeout: time.Duration(router.Timeout),
	}
	return action
}

func toMap(in interface{}) map[string]interface{} {
	var out map[string]interface{}
	data, _ := json.Marshal(in)
	json.Unmarshal(data, &out)
	return out
}

func convertHealthCheckCfg(config v1.HealthCheckConfig) v2.HealthCheck {
	return v2.HealthCheck{
		HealthCheckConfig: v2.HealthCheckConfig{
			Protocol:             config.Protocol,
			TimeoutConfig:        api.DurationConfig{Duration: config.TimeoutConfig},
			IntervalConfig:       api.DurationConfig{Duration: config.IntervalConfig},
			IntervalJitterConfig: api.DurationConfig{Duration: config.IntervalJitterConfig},
			HealthyThreshold:     config.HealthyThreshold,
			UnhealthyThreshold:   config.UnhealthyThreshold,
			ServiceName:          config.ServiceName,
			SessionConfig:        config.SessionConfig,
			CommonCallbacks:      config.CommonCallbacks,
		},
	}
}

func convertTlCfg(config v1.TLSConfig) v2.TLSConfig {
	return v2.TLSConfig{
		Status:       config.Status,
		Type:         config.Type,
		ServerName:   config.ServerName,
		CACert:       config.CACert,
		CertChain:    config.CertChain,
		PrivateKey:   config.PrivateKey,
		VerifyClient: config.VerifyClient,
		//RequireClientCert: config.RequireClientCert,
		InsecureSkip: config.InsecureSkip,
		CipherSuites: config.CipherSuites,
		EcdhCurves:   config.EcdhCurves,
		MinVersion:   config.MinVersion,
		MaxVersion:   config.MaxVersion,
		ALPN:         config.ALPN,
		Ticket:       config.Ticket,
		Fallback:     config.Fallback,
		ExtendVerify: config.ExtendVerify,
		//SdsConfig:         config.SdsConfig,
	}
}
