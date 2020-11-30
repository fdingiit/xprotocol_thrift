module gitlab.alipay-inc.com/ant-mesh/mosn

go 1.12

require (
	github.com/c2h5oh/datasize v0.0.0-20200112174442-28bbd4740fee // indirect
	github.com/kr/pretty v0.2.0 // indirect
	go.uber.org/zap v1.14.1
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/tools v0.0.0-20201014231627-1610a49f37af // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	mosn.io/mosn v0.15.0
	mosn.io/pkg v0.0.0-20200729115159-2bd74f20be0f
)

replace (
	github.com/apache/dubbo-go-hessian2 => github.com/apache/dubbo-go-hessian2 v1.4.1-0.20200516085443-fa6429e4481d // perf: https://github.com/apache/dubbo-go-hessian2/pull/188
	github.com/envoyproxy/go-control-plane => gitlab.alipay-inc.com/cloudnative/cloudmesh-go-control-plane v0.0.0-20200602015852-5413b57f5d72
	istio.io/api => gitlab.alipay-inc.com/cloudnative/cloudmesh-api v0.0.0-20191220062600-8ef8a28afc04
	mosn.io/mosn => github.com/rickey17/mosn v0.15.12
)
