package api

import (
	"reflect"
)

func init() {
	register(GW_CONFIG, GATEWAY, ROUTER_GROUP, ROUTER, FILTER_CHAIN, GLOBAL_FILTER, SERVICE, METADATA, EXTENSION, INNER)
}

type Type string

const (
	V1Version         = "v1"
	K8S_BETA1_VERSION = "apiextensions.k8s.io/v1beta1"
)

const (
	GW_CONFIG     Type = "GatewayConfig"
	GATEWAY       Type = "gateway"
	ROUTER_GROUP  Type = "routerGroup"
	ROUTER        Type = "router"
	FILTER_CHAIN  Type = "filterChain"
	GLOBAL_FILTER Type = "globalFilter"
	SERVICE       Type = "service"
	METADATA      Type = "Metadata"
	EXTENSION     Type = "extension"
	INNER         Type = "inner"
)

type Object interface {
	Type() Type
}

func GetType(name string) Type {
	return typeMapping[name]
}

var typeMapping = make(map[string]Type)

func register(types ...Type) {
	for _, t := range types {
		valueOf := reflect.ValueOf(t)
		typeMapping[valueOf.String()] = t
	}
}
