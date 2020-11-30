package types

import (
	"mosn.io/mosn/pkg/protocol"
)

type Pipeline interface {
	//
	DoInBound(Context) error

	//
	DoOutBound(Context) error

	AddOrUpdateFilter(GatewayFilter) error

	AddNewFilter([]GatewayFilter) error

	DelFilter(GatewayFilter) error

	GetFilters() []GatewayFilter

	Sort()

	Copy() Pipeline

	SetErrHandler(ErrorHandler)

	HandleErr(Context, GatewayError) (httpCode int, headers protocol.CommonHeader, res []byte)
}

type WeightPipeline interface {
	GetPipeline() Pipeline
}
