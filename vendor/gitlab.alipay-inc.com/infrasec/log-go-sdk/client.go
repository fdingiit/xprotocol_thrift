package sls

import (
	"bytes"
	"fmt"
	aliyun_sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/gogo/protobuf/jsonpb"
	"gitlab.alipay-inc.com/infrasec/api/types"
	"os"
	"time"
)

const (
	defaultBufferSize           = 100
	defaultPushPeriod           = 60
	defaultChanSize             = 200
	defaultClientRequestTimeout = 5
	defaultClientRetryTimeout   = 15
)

type LogClient struct {
	logStore *aliyun_sls.LogStore

	topic      string
	source     string
	logBuffer  []*aliyun_sls.Log // remember logBuffer is in critical section
	bufferSize int32
	pushPeriod int32

	logChan chan map[string]string
}

func NewLogClientWithBytes(b []byte) (p *LogClient, err error) {
	if b == nil || len(b) <= 0 {
		return nil, fmt.Errorf("sls config is none. ")
	}
	conf := &types.SlsClientConf{}
	var un jsonpb.Unmarshaler
	if err := un.Unmarshal(bytes.NewReader(b), conf); err != nil {
		return nil, err
	}
	return NewLogClient(conf)
}

func NewLogClient(conf *types.SlsClientConf) (p *LogClient, err error) {
	if conf == nil {
		return nil, fmt.Errorf("the conf is nil")
	}

	switch conf.GetMode() {
	case types.Mode_LOCAL_CRYPT:
		if conf.GetKey() == nil {
			return nil, fmt.Errorf("the key is none")
		}
		accessID, err := Decrypt(conf.GetAccessId(), string(conf.GetKey()))
		if err != nil {
			return nil, err
		}
		accessSecret, err := Decrypt(conf.GetAccessSecret(), string(conf.GetKey()))
		if err != nil {
			return nil, err
		}

		return newLogProject(
			conf,
			accessID,
			accessSecret)

	case types.Mode_MIST_CRYPT:
		if conf.GetHandler() == "" {
			return nil, fmt.Errorf("the handler is none")
		}
		//TODO
		return nil, fmt.Errorf("not support")
	default:
		return newLogProject(
			conf,
			conf.AccessId,
			conf.AccessSecret)
	}
}

func newLogProject(conf *types.SlsClientConf, accessKeyID, accessKeySecret string) (*LogClient, error) {
	logProject, err := aliyun_sls.NewLogProject(conf.ProjectName, conf.Endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	if conf.Token != "" {
		logProject.WithToken(conf.Token)
	}

	logProject.UsingHTTP = true
	//logProject.UsingHTTP = conf.UsingHttp

	if conf.UserAgent != "" {
		logProject.UserAgent = conf.UserAgent
	}

	logProject.WithRequestTimeout(time.Duration(defaultInt32(conf.RequestTimeout, defaultClientRequestTimeout)) * time.Second)
	logProject.WithRetryTimeout(time.Duration(defaultInt32(conf.RetryTimeout, defaultClientRetryTimeout)) * time.Second)

	logStore, err := logProject.GetLogStore(conf.GetLogstoreName())
	if err != nil {
		return nil, err
	}
	client := &LogClient{
		source:     hostname,
		topic:      fmt.Sprintf("%s-topic", conf.GetLogstoreName()),
		logStore:   logStore,
		pushPeriod: defaultInt32(conf.GetPushPeriod(), defaultPushPeriod),
		logBuffer:  make([]*aliyun_sls.Log, 0, defaultInt32(conf.GetBufferSize(), defaultBufferSize)),
		bufferSize: defaultInt32(conf.GetBufferSize(), defaultBufferSize),
		logChan:    make(chan map[string]string, defaultInt32(conf.GetChanSize(), defaultChanSize)),
	}

	client.consumer()

	return client, nil
}

func defaultInt32(i, d int32) int32 {
	if i <= 0 {
		return d
	}
	return i
}
