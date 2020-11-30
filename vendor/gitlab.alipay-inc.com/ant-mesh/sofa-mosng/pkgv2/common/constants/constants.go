package constants

type Direction int

const (
	MosngServiceHeader = "x-mosng-upstream-name"

	Request Direction = iota
	Response
)

const (
	AppLogConfigFile = "log_config_file"
)

const (
	AttrSofaPenAttrs       = "SofaPenAttrs"
	AttrSofaRpcRouterLocal = "sofa_head_target_local"
	AttrTracerCtx          = "rpc_trace_context"
)

const (
	ContextKeyTraceId = "sofaTraceId"
)
