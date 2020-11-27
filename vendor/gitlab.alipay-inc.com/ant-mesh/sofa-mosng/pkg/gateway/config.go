package gateway

import (
	"github.com/json-iterator/go"
	"mosn.io/api"
	"mosn.io/mosn/pkg/protocol"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

const (
	GatewayHandlerRuleName  = "handler_rule"
	GatewayKeystoreRuleName = "keystore_rule"

	AppSignConfigName     = "sign_config"
	ClusterSignConfigName = "sign_config"
)

func init() {
	RegisterGatewayRuleType(GatewayHandlerRuleName, &GatewayHandlerRule{})
	RegisterGatewayRuleType(GatewayKeystoreRuleName, &GatewayKeystoreRule{})

	RegisterAppSpecType(AppSignConfigName, &SignConfig{})
	RegisterClusterSpecType(ClusterSignConfigName, &SignConfig{})
}

type GatewayFilterConfig struct {
	DownstreamProtocol string   `json:"downstream_protocol"`
	UpstreamProtocol   string   `json:"upstream_protocol"`
	Handlers           []string `json:"handlers"`
	ConfigPath         string   `json:"config_path"`
	LogFolder          string   `json:"log_folder"`
}

type GatewayConfig struct {
	ServiceConfigs []ServiceConfig `json:"service_list"`
	GatewayRules   []GatewayRule   `json:"gateway_rule"`
	AppConfigs     []AppConfig     `json:"app_list"`
	ClusterConfigs []ClusterConfig `json:"cluster_list"`
}

type ServiceConfig struct {
	Id       string         `json:"id"`
	Status   ApiStatus      `json:"status"`
	Timeout  uint64         `json:"timeout,omitempty"`
	Upstream UpstreamConfig `json:"upstream"`
	Rules    []ServiceRule  `json:"rules"`
}

type UpstreamConfig struct {
	Protocol    api.Protocol           `json:"protocol"`
	ClusterName string                 `json:"cluster_name"`
	TimeOut     uint64                 `json:"timeout"`
	Config      map[string]interface{} `json:"config"`
}

type ServiceRule struct {
	Name   string
	Config interface{}
}

type GatewayRule struct {
	Name   string      `json:"name"`
	Config interface{} `json:"config"`
}

type GatewayHandlerRule struct {
	Handlers []string `json:"handlers"`
}

type GatewayKeystoreRule struct {
	Secrets map[string]string `json:"secrets"`
}

type AppConfig struct {
	Id   string                 `json:"id"`
	Name string                 `json:"name"`
	Spec map[string]interface{} `json:"spec"`
}

type ClusterConfig struct {
	Id   string                 `json:"id"`
	Name string                 `json:"name"`
	Spec map[string]interface{} `json:"spec"`
}

type SignConfig struct {
	Algorithm string   `json:"algorithm"`
	KeyId     string   `json:"key_id"`
	Headers   []string `json:"headers"`
}

var ProtocolsSupported = map[string]bool{
	string(protocol.Auto):      true,
	string(protocol.SofaRPC):   true,
	string(protocol.HTTP2):     true,
	string(protocol.HTTP1):     true,
	string(protocol.Xprotocol): true,
}

func ParseGatewayStreamFilter(cfg map[string]interface{}) (*GatewayFilterConfig, error) {
	filterConfig := &GatewayFilterConfig{}
	data, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, filterConfig); err != nil {
		return nil, err
	}
	return filterConfig, nil
}
