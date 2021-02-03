module github.com/fdingiit/xprotocol_thrift

go 1.12

require (
	github.com/kr/pretty v0.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/zap v1.14.1
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/tools v0.0.0-20201014231627-1610a49f37af // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	mosn.io/api v0.0.0-20210113033009-f24f4e59b2bc
	mosn.io/pkg v0.0.0-20201228070559-80e9ae937bd5
)

replace (
	mosn.io/api => github.com/fdingiit/api v0.0.0-20210119063843-f5ac263a02c5
	mosn.io/pkg => github.com/fdingiit/pkg v0.0.0-20210119065649-3fcc1522bbbd
)
