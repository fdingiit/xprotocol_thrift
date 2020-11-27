module gitlab.alipay-inc.com/ant-mesh/mosn

go 1.12

require (
	github.com/Knetic/govaluate v3.0.1-0.20171022003610-9aa49832a739+incompatible
	github.com/SkyAPM/go2sky v0.5.0
	github.com/alibaba/sentinel-golang v0.2.1-0.20200509115140-6d505e23ef30
	github.com/aliyun/alibaba-cloud-sdk-go v1.61.577
	github.com/aliyun/aliyun-log-go-sdk v0.1.6
	github.com/apache/dubbo-go-hessian2 v1.5.0
	github.com/apache/thrift v0.13.0
	github.com/deckarep/golang-set v1.7.1
	github.com/eapache/queue v0.0.0-20180227141424-093482f3f8ce
	github.com/envoyproxy/go-control-plane v0.9.4
	github.com/gin-gonic/gin v1.6.2
	github.com/go-playground/validator/v10 v10.2.0
	github.com/go-zookeeper/zk v1.0.2
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.5
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/json-iterator/go v1.1.9
	github.com/julienschmidt/httprouter v1.3.0
	github.com/magiconair/properties v1.8.1
	github.com/mediocregopher/mediocre-go-lib v0.0.0-20190730033908-c20f884d6844 // indirect
	github.com/nacos-group/nacos-sdk-go v0.0.0-20190723125407-0242d42e3dbb
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d
	github.com/openzipkin/zipkin-go v0.1.1
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.2.1
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.9.1 // indirect
	github.com/prometheus/procfs v0.0.8 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20200313005456-10cdbea86bc0
	github.com/satori/go.uuid v1.2.0
	github.com/shirou/gopsutil v2.20.6+incompatible
	github.com/stretchr/testify v1.5.1
	github.com/urfave/cli v1.20.0
	github.com/valyala/fasthttp v1.14.1-0.20200605121233-ac51d598dc54
	github.com/valyala/fastjson v1.4.5 // indirect
	github.com/yuin/gopher-lua v0.0.0-20191220021717-ab39c6098bdb
	gitlab.alipay-inc.com/AntTechRisk/awatch-go v0.0.0-20200818032039-9fb24a779546
	gitlab.alipay-inc.com/ant-mesh/sofa-mosng v0.0.0-20200806014231-5a4f2cfbde89
	gitlab.alipay-inc.com/ant_agent/client_golang v0.0.0-20191106131121-98759bd0281a
	gitlab.alipay-inc.com/infrasec/api v1.3.0
	gitlab.alipay-inc.com/infrasec/log-go-sdk v1.0.1
	gitlab.alipay-inc.com/infrasec/mist-go-sdk v1.0.3
	gitlab.alipay-inc.com/infrasec/opa v1.3.0
	gitlab.alipay-inc.com/octopus/radix v0.0.0-20200205031526-430a99e67ade // indirect
	gitlab.alipay-inc.com/octopus/tbase-go-client v1.1.0
	gitlab.alipay-inc.com/sofa-go/sofa-antvip-client-go v0.3.6
	gitlab.alipay-inc.com/sofa-go/sofa-drm-client-go v1.1.3
	gitlab.alipay-inc.com/sofa-go/sofa-hessian-go v0.1.5
	gitlab.alipay-inc.com/sofa-go/sofa-logger-go v0.2.4
	gitlab.alipay-inc.com/sofa-go/sofa-registry-client-go v0.5.3
	gitlab.alipay-inc.com/sofa-go/sofa-zoneclient-go v0.3.8
	gitlab.alipay-inc.com/sofa-open/eureka-go v0.0.0-20201105112636-32e44f2a9a48
	go.uber.org/atomic v1.6.0
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/net v0.0.0-20200822124328-c89045814202
	golang.org/x/text v0.3.2
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	golang.org/x/tools v0.0.0-20201014231627-1610a49f37af // indirect
	google.golang.org/grpc v1.28.0
	istio.io/api v0.0.0-20200227213531-891bf31f3c32
	mosn.io/api v0.0.0-20200729124336-c71e8f2074cb
	mosn.io/mosn v0.15.0
	mosn.io/pkg v0.0.0-20200729115159-2bd74f20be0f
)

replace (
	github.com/apache/dubbo-go-hessian2 => github.com/apache/dubbo-go-hessian2 v1.4.1-0.20200516085443-fa6429e4481d // perf: https://github.com/apache/dubbo-go-hessian2/pull/188
	github.com/envoyproxy/go-control-plane => gitlab.alipay-inc.com/cloudnative/cloudmesh-go-control-plane v0.0.0-20200602015852-5413b57f5d72
	istio.io/api => gitlab.alipay-inc.com/cloudnative/cloudmesh-api v0.0.0-20191220062600-8ef8a28afc04
	mosn.io/mosn => github.com/rickey17/mosn v0.15.12
)
