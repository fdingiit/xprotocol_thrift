package model

import (
	fmt "fmt"
	"time"

	"gitlab.alipay-inc.com/sofa-go/sofa-antvip-client-go/sofaantvip"
)

type AttributeGetResponse struct {
	data []byte
}

func (r AttributeGetResponse) Marshal() ([]byte, error) {
	return r.data, nil
}

func (r *AttributeGetResponse) SetData(s []byte) {
	r.data = append(r.data[:0], s...)
}

func (r *AttributeGetResponse) SetDataString(s string) {
	r.data = append(r.data[:0], s...)
}

func (r *AttributeGetResponse) Reset() {
	r.data = r.data[:0]
}

func (r *AttributeGetResponse) String() string {
	return string(r.data)
}

func (r *AttributeGetResponse) ProtoMessage() {
}

type HeartbeatRequest struct {
	Zone          string
	ClientIp      string
	InstanceId    string
	Profile       string
	VersionMap    map[string]int32
	AckVersionMap map[string]int32
}

type HeartbeatResponse struct {
	WaitTime int64
	DiffMap  map[string]int32
}

type SubscriberRegReq struct {
	Zone       string
	DataId     string
	Uuid       string
	InstanceId string
	Profile    string
	Attributes map[string]string
}

type SubscriberRegResult struct {
	Zone       string
	DataId     string
	Uuid       string
	InstanceId string
	Profile    string
	Attributes map[string]string
	Result     bool
	Message    string
}

type AttributeGetRequest struct {
	DataId string
	Value  string
}

type AttributeSetRequest struct {
	DataId string
	Value  string
}

type ClientConfig struct {
	AppName               string
	DataCenter            string
	Zone                  string
	InstanceId            string
	Profile               string
	EndPoint              string
	DomainPrefix          string
	DomainSuffix          string
	AccessKey             string
	AccessSecret          string
	ConnectRetryDuration  time.Duration
	HeartbeatDuration     time.Duration
	RegisterCheckDuration time.Duration
}

type Server struct {
	Url     string
	Ip      string
	Port    int32
	Timeout int32
}

type DrmValue struct {
	Value   string
	Version int
}

type DrmAck struct {
	ActTime time.Time
	Version int
}

//go:generate syncmap -pkg model -o cachemap_generated.go -name LocalValueCacheMap map[string]LocalValueCache

type LocalValueCache struct {
	DataId   string
	DrmValue DrmValue
	DrmAck   DrmAck
}

func NewAvailableServers(servers []sofaantvip.RealServer, protectThreshold int32, port int32) []Server {
	result := []Server{}
	noAvailable := []Server{}
	for _, server := range servers {
		drmServer := Server{
			Ip:   server.GetIp(),
			Port: port,
		}
		if server.IsAvailable() {
			result = append(result, drmServer)
		} else {
			noAvailable = append(noAvailable, drmServer)
		}
	}

	length := float32(len(servers))
	needCount := float32(protectThreshold) / float32(100) * length
	for _, no := range noAvailable {
		if float32(len(result)) >= needCount {
			return result
		}
		result = append(result, no)
	}
	return result
}

func (srv *Server) GetIPPort() string {
	return fmt.Sprintf("%s:%d", srv.Ip, srv.Port)
}
