package tbasego

import (
	"fmt"
	"sync/atomic"
	"time"

	tbase_log "gitlab.alipay-inc.com/octopus/tbase-go-client/tloggger"

	time_helper "gitlab.alipay-inc.com/octopus/tbase-go-client/utils/time.helper"

	"gitlab.alipay-inc.com/octopus/tbase-go-client/model"
	"mosn.io/pkg/log"
)

var clientCounter int32 = 0

type TBaseClient struct {
	ConnectionInfo  *model.ConnectionInfo
	commandExecutor *CommandExecutor
	instanceName    string
	ShardManager    *ShardManager
	closed          bool
}

func NewTBaseClient2(connectionInfo *model.ConnectionInfo) (*TBaseClient, error) {
	tbase_log.InitTBaseLogger()

	instanceName := fmt.Sprintf("TBaseClient-(%v.%v#%v)", connectionInfo.Cluster, connectionInfo.Tenant, atomic.AddInt32(&clientCounter, 1))
	if tbase_log.TBaseLogger.GetLogLevel() >= log.INFO {
		tbase_log.TBaseLogger.Infof("[TBASE_CLIENT] init tbase-client instance '%v' with connection string: '%v'", instanceName, connectionInfo.ToString())
	}
	shardManager, err := NewShardManager(connectionInfo)
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[TBASE_CLIENT] init tbase-client instance '%v' error, initialize shard manager error, error: %v", err)
		return nil, err
	}

	tClient := &TBaseClient{
		ConnectionInfo:  connectionInfo,
		closed:          false,
		ShardManager:    shardManager,
		commandExecutor: newCommandExecutor(connectionInfo, newClientHandler(newClientPool(), shardManager)),
		instanceName:    instanceName,
	}

	if connectionInfo.Warmup {
		if connectionInfo.Handshake {
			// todo handle hand shake
		} else {
			// todo handle warmup
		}
	}

	return tClient, nil
}

func NewTBaseClient1(connectionStr string) (*TBaseClient, error) {
	connectionInfo, err := model.Parse(connectionStr)
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[TBASE_CLIENT] init tbase-client with connection string: '%v' error, error: %v", connectionStr, err)
		return nil, err
	}
	return NewTBaseClient2(connectionInfo)
}

// A print like to call redis command and get back result. But it only support string command
// If the receiver value is a primitive or slice/map a pointer must be passed in.
// Example: Do(nil, "SET", "foo", "bar")
//          Do(&output, "GET", "foo")
// if there's something wrong, an 'error' will be returned
func (tClient *TBaseClient) Do(rcv interface{}, cmd string, args ...string) (string, error) {
	return tClient.commandExecutor.do(tClient.buildTBaseAction(rcv, cmd, args...))
}

// A print like to call redis command and get back result. It like Do(xxx) but the arguments can be of almost any type
// But it does not work for commands whose first parameter isn't a key.
// Likely if the receiver value is a primitive or slice/map a pointer must be passed in.
// Example: FlatDo(nil, "SETEX", "fool", ttl, value)
// if there's something wrong, an 'error' will be returned
func (tClient *TBaseClient) FlatDo(rcv interface{}, cmd, key string, args ...interface{}) (string, error) {
	return tClient.commandExecutor.do(tClient.BuildTBaseActionForFlat(rcv, cmd, key, args))
}

func (tClient *TBaseClient) SetSocketTimeout(socketTimeout time.Duration) {
	tClient.commandExecutor.socketTimeout = socketTimeout
}

func (tClient *TBaseClient) Close() {
	if !tClient.closed {
		tClient.closed = true
		if tClient.ShardManager != nil {
			tClient.ShardManager.Close()
		}
		if tClient.commandExecutor != nil {
			tClient.commandExecutor.Close()
		}
		if tbase_log.TBaseLogger.GetLogLevel() >= log.INFO {
			tbase_log.TBaseLogger.Infof("[TBASE_CLIENT] tbase client %v is closed. ", tClient.instanceName)
		}
	}
}

func (tClient *TBaseClient) buildTBaseAction(rcv interface{}, cmd string, args ...string) *model.TBaseAction {
	return &model.TBaseAction{
		Attempts:   tClient.ConnectionInfo.MaxRetries,
		SubmitTime: time.Now().UnixNano(),
		Timeout:    time_helper.MsToNs(tClient.ConnectionInfo.RedisTimeout),
		Rcv:        rcv,
		Cmd:        cmd,
		Args:       args,
	}
}

func (tClient *TBaseClient) BuildTBaseActionForFlat(rcv interface{}, cmd, key string, args ...interface{}) *model.TBaseAction {
	return &model.TBaseAction{
		Attempts:   tClient.ConnectionInfo.MaxRetries,
		SubmitTime: time.Now().UnixNano(),
		Timeout:    time_helper.MsToNs(tClient.ConnectionInfo.RedisTimeout),
		Rcv:        rcv,
		Cmd:        cmd,
		Flat:       true,
		FlatKey:    [1]string{key},
		FlatArgs:   args,
	}
}
