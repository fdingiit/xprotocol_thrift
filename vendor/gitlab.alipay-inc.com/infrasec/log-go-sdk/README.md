# usage
```
conf := &types.SlsClientConf{
	Endpoint:     "cn-shanghai-ant-share.log.aliyuncs.com",
	ProjectName:  "ant-meshsecurity-dev",
	LogstoreName: "rbac-events-log",
	AccessId:     "xxx",
	AccessSecret: "yyy",
}

client, err := sls.NewLogClient(conf)

if err != nil {
    return
}
defer client.Stop()

testData := map[string]string{
	"ut-key": "ut-value",
}
ret := client.Producer(testData)
```

# test
```
ok      gitlab.alipay-inc.com/infrasec/log-go-sdk  5.393s  coverage: 86.4% of statements
```
# benchmark
```
pkg: gitlab.alipay-inc.com/infrasec/log-go-sdk
BenchmarkClientProducer-8        4600032               223 ns/op
```
