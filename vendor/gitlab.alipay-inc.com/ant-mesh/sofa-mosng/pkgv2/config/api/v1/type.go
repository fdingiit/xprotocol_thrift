package v1

import (
	v2 "mosn.io/mosn/pkg/config/v2"
	"time"

	"gitlab.alipay-inc.com/ant-mesh/sofa-mosng/pkgv2/config/api"
)

func NewObject(t api.Type) api.Object {
	switch t {
	case api.GW_CONFIG:
		return &GwConfig{}
	case api.GATEWAY:
		return &Gateway{}
	case api.ROUTER:
		return &Router{}
	case api.ROUTER_GROUP:
		return &RouterGroup{}
	case api.FILTER_CHAIN:
		return &FilterChain{}
	case api.GLOBAL_FILTER:
		return &GlobalFilter{}
	case api.SERVICE:
		return &GatewayService{}
	case api.METADATA:
		return &Metadata{}
	case api.EXTENSION:
		return &Extension{}
	}

	return &GwConfig{}
}

type GwConfig struct {
	api.Object
	Gateways         []*Gateway        `yaml:"gateways,omitempty" json:"gateways,omitempty"`
	FilterChains     []*FilterChain    `yaml:"filterChains,omitempty" json:"filter_chains,omitempty" json:"filterChains,omitempty"`
	RouterGroups     []*RouterGroup    `yaml:"routerGroups,omitempty" json:"router_groups,omitempty" json:"routerGroups,omitempty"`
	Routers          []*Router         `yaml:"routers,omitempty" json:"routers,omitempty" json:"routers,omitempty"`
	GatewayServices  []*GatewayService `yaml:"gatewayServices,omitempty" json:"gateway_services,omitempty" json:"gatewayServices,omitempty"`
	GatewayMetadatas []*Metadata       `yaml:"gatewayMetadatas,omitempty" json:"gateway_metadatas,omitempty" json:"gatewayMetadatas,omitempty"`
	GlobalFilters    GlobalFilter      `yaml:",inline" json:"global_filter,omitempty" json:"globalFilter,omitempty"`
	Extension        Extension         `yaml:"extension,omitempty" json:"extension,omitempty"`
}

func (GwConfig) Type() api.Type {
	return api.GW_CONFIG
}

type Metadata struct {
	Key   string      `yaml:"key"`
	Value interface{} `yaml:"value,omitempty"`
}

func (Metadata) Type() api.Type {
	return api.METADATA
}

type Gateway struct {
	Name                  string                 `yaml:"name"`
	AddrConfig            string                 `yaml:"addr"`
	ListenerType          v2.ListenerType        `yaml:"ListenerType,omitempty"`
	BindToPort            bool                   `yaml:"bindPort,omitempty"`
	UpstreamProtocol      string                 `yaml:"upstreamProtocol"`
	DownstreamProtocol    string                 `yaml:"downstreamProtocol"`
	Inspector             bool                   `yaml:"inspector,omitempty"`
	ConnectionIdleTimeout time.Duration          `yaml:"connectionIdleTimeout,omitempty"`
	TLSConfig             TLSConfig              `yaml:"tls,omitempty"`
	TLSConfigs            []TLSConfig            `yaml:"tlsSet,omitempty"`
	FilterChain           string                 `yaml:"filterChain,omitempty"`
	Metadata              map[string]interface{} `yaml:"meta,omitempty"`
}

func (Gateway) Type() api.Type {
	return api.GATEWAY
}

type Extension struct {
	LuaPath string            `yaml:"luaPath,omitempty"`
	Filters []FilterExtension `yaml:"filters,omitempty"`
}

func (Extension) Type() api.Type {
	return api.EXTENSION
}

type FilterExtension struct {
	Name string `yaml:"name,omitempty"`
	File string `yaml:"file,omitempty"`
}

type FilterChain struct {
	ChainName string    `yaml:"name" json:"name"`
	Filters   []*Filter `yaml:"filters" json:"filters"`
}

func (FilterChain) Type() api.Type {
	return api.FILTER_CHAIN
}

type GlobalFilter struct {
	Filters []*Filter `yaml:"globalFilters,omitempty" json:"globalFilters,omitempty" json:"global_filters,omitempty"`
}

func (GlobalFilter) Type() api.Type {
	return api.GLOBAL_FILTER
}

type Filter struct {
	Name     string      `yaml:"name" json:"name"`
	Priority int64       `yaml:"priority,omitempty" json:"priority,omitempty"`
	Metadata interface{} `yaml:"meta,omitempty" json:"meta,omitempty"`
}

//func (Filter) ResourceType() ResourceType {
//	return FILTER
//}

type Status string

const (
	OPEN  Status = "OPEN"
	CLOSE Status = "CLOSE"
)

type Router struct {
	api.Object
	Status       Status                 `yaml:"status,omitempty" json:"status,omitempty"`
	GroupBind    string                 `yaml:"groupBind,omitempty" json:"groupBind,omitempty"`
	Name         string                 `yaml:"name" json:"name,omitempty"`
	Matches      *RouterMatch           `yaml:"matches" json:"matches,omitempty"`
	Proxy        *Proxy                 `yaml:"proxy,omitempty" json:"proxy,omitempty"`
	WeightProxy  []WeightProxy          `yaml:"weightProxy,omitempty" json:"weightProxy,omitempty"`
	FilterChains []string               `yaml:"filterChains,omitempty" json:"filterChains,omitempty" json:"filter_chains,omitempty"`
	Filters      []*Filter              `yaml:"filters,omitempty" json:"filters,omitempty"`
	Timeout      uint64                 `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	Metadata     map[string]interface{} `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

func (r *Router) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type rawRouter Router
	raw := rawRouter{Status: OPEN}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	*r = Router(raw)
	return nil
}

type RouterMatch struct {
	Prefix       string                 `yaml:"prefix,omitempty" json:"prefix,omitempty"`                                                               // Match request's Path with Prefix Comparing
	Path         string                 `yaml:"path,omitempty" json:"path,omitempty"`                                                                   // Match request's Path with Exact Comparing
	Regex        bool                   `yaml:"regex,omitempty" json:"regex,omitempty"`                                                                 // Match request's Path with Regex Comparing
	Headers      []*ValueMatcher        `yaml:"headerMatcher,omitempty" json:"headerMatcher,omitempty" json:"header_matcher,omitempty"`                 // Match request's Headers
	QueryStrings []*ValueMatcher        `yaml:"queryStringMatcher,omitempty" json:"queryStringMatcher,omitempty" json:"query_string_matcher,omitempty"` // Match request's Parameters
	Metadata     map[string]interface{} `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

type ValueMatcher struct {
	Name  string `yaml:"name,omitempty" json:"name,omitempty"`
	Value string `yaml:"value,omitempty" json:"value,omitempty"`
	Regex bool   `yaml:"regex,omitempty" json:"regex,omitempty"`
}

func (Router) Type() api.Type {
	return api.ROUTER
}

type RouterGroup struct {
	GroupName string    `yaml:"name" json:"name"`
	Gateways  []string  `yaml:"gateways" json:"gateways"`
	Routers   []*Router `yaml:"routers,omitempty" json:"routers,omitempty"`
}

func (RouterGroup) Type() api.Type {
	return api.ROUTER_GROUP
}

type Proxy struct {
	Protocol  string                 `yaml:"protocol,omitempty" json:"protocol,omitempty"`
	Interface string                 `yaml:"interface,omitempty" json:"interface,omitempty"`
	Method    string                 `yaml:"method,omitempty" json:"method,omitempty"`
	Service   string                 `yaml:"service" json:"service"`
	Timeout   uint64                 `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	Metadata  map[string]interface{} `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

type WeightProxy struct {
	Protocol  string                 `yaml:"protocol,omitempty" json:"protocol,omitempty"`
	Interface string                 `yaml:"interface,omitempty" json:"interface,omitempty"`
	Method    string                 `yaml:"method,omitempty" json:"method,omitempty"`
	Service   string                 `yaml:"service" json:"service"`
	Timeout   uint64                 `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	Weight    int32                  `yaml:"weight" json:"weight"`
	Metadata  map[string]interface{} `yaml:"metadata,omitempty" json:"metadata,omitempty"`
}

type GatewayService struct {
	Name                 string            `yaml:"name" json:"name"`
	Protocol             string            `yaml:"protocol,omitempty"`
	Hosts                []*HostConfig     `yaml:"hosts"`
	Timeout              uint64            `yaml:"timeout,omitempty"`
	LbType               v2.LbType         `yaml:"loadBalance"`
	ClusterType          v2.ClusterType    `yaml:"type,omitempty"`
	SubType              string            `json:"sub_type,omitempty"`
	CircuitBreakers      CircuitBreakers   `yaml:"circuitBreakers,omitempty"`
	Filters              []*Filter         `yaml:"filters,omitempty"`
	FilterChains         []string          `yaml:"filterChains,omitempty"`
	MaxRequestPerConn    uint32            `yaml:"maxRequestPerConn,omitempty"`
	ConnBufferLimitBytes uint32            `yaml:"connBufferLimitBytes,omitempty"`
	Spec                 ClusterSpecInfo   `yaml:"spec,omitempty"`
	LBSubSetConfig       LBSubsetConfig    `yaml:"lbSubsetConfig,omitempty"`
	TLS                  TLSConfig         `yaml:"tls,omitempty"`
	HealthCheck          HealthCheckConfig `yaml:"healthCheck,omitempty"`
}

type HostConfig struct {
	Address  string `json:"address,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	Weight   uint32 `json:"weight,omitempty"`
	//MetaDataConfig *MetadataConfig `json:"meta,omitempty"`
	TLSDisable bool `json:"tls_disable,omitempty"`
}

type HealthCheckConfig struct {
	Protocol             string                 `json:"protocol,omitempty"`
	TimeoutConfig        time.Duration          `json:"timeout,omitempty"`
	IntervalConfig       time.Duration          `json:"interval,omitempty"`
	IntervalJitterConfig time.Duration          `json:"interval_jitter,omitempty"`
	HealthyThreshold     uint32                 `json:"healthy_threshold,omitempty"`
	UnhealthyThreshold   uint32                 `json:"unhealthy_threshold,omitempty"`
	ServiceName          string                 `json:"service_name,omitempty"`
	SessionConfig        map[string]interface{} `json:"check_config,omitempty"`
	CommonCallbacks      []string               `json:"common_callbacks,omitempty"` // HealthCheck support register some common callbacks that are not related to specific cluster
}

type TLSConfig struct {
	Status            bool                   `yaml:"status,omitempty"`
	Type              string                 `yaml:"type,omitempty"`
	ServerName        string                 `yaml:"serverName,omitempty"`
	CACert            string                 `yaml:"caCert,omitempty"`
	CertChain         string                 `yaml:"certChain,omitempty"`
	PrivateKey        string                 `yaml:"privateKey,omitempty"`
	VerifyClient      bool                   `yaml:"verifyClient,omitempty"`
	RequireClientCert bool                   `yaml:"requireClientCert,omitempty"`
	InsecureSkip      bool                   `yaml:"insecureSkip,omitempty"`
	CipherSuites      string                 `yaml:"cipherSuites,omitempty"`
	EcdhCurves        string                 `yaml:"ecdhCurves,omitempty"`
	MinVersion        string                 `yaml:"minVersion,omitempty"`
	MaxVersion        string                 `yaml:"maxVersion,omitempty"`
	ALPN              string                 `yaml:"alpn,omitempty"`
	Ticket            string                 `yaml:"ticket,omitempty"`
	Fallback          bool                   `yaml:"fallback,omitempty"`
	ExtendVerify      map[string]interface{} `yaml:"extendVerify,omitempty"`
}

type LBSubsetConfig struct {
	FallBackPolicy  uint8             `yaml:"fallBackPolicy,omitempty"`
	DefaultSubset   map[string]string `yaml:"defaultSubset,omitempty"`
	SubsetSelectors [][]string        `yaml:"subsetSelectors,omitempty"`
}

type ClusterSpecInfo struct {
	Subscribes []SubscribeSpec `yaml:"subscribe,omitempty"`
}

type SubscribeSpec struct {
	Subscriber  string `yaml:"subscriber,omitempty"`
	ServiceName string `yaml:"serviceName,omitempty"`
}

func (GatewayService) Type() api.Type {
	return api.SERVICE
}

type CircuitBreakers struct {
	Thresholds []Thresholds
}

type Thresholds struct {
	MaxConnections     uint32 `yaml:"maxConnections,omitempty"`
	MaxPendingRequests uint32 `yaml:"maxPendingRequests,omitempty"`
	MaxRequests        uint32 `yaml:"maxRequests,omitempty"`
	MaxRetries         uint32 `yaml:"maxRetries,omitempty"`
}
