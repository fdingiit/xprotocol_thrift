module github.com/fdingiit/xprotocol_thrift

go 1.15

require (
	go.uber.org/zap v1.16.0
	mosn.io/mosn v0.18.0
	mosn.io/pkg v0.0.0-20200729115159-2bd74f20be0f
)

replace mosn.io/mosn => github.com/rickey17/mosn v0.15.12
