package tbasego

import (
	"net"
	"sync"
	"sync/atomic"
	"time"

	tbase_log "gitlab.alipay-inc.com/octopus/tbase-go-client/tloggger"

	"gitlab.alipay-inc.com/octopus/tbase-go-client/async"

	. "gitlab.alipay-inc.com/octopus/tbase-go-client/model/error"
	"mosn.io/pkg/log"
)

type ClientPool struct {
	endpoint2ClientResourceMap sync.Map
	syncMutex                  *sync.RWMutex
	closed                     bool
}

func newClientPool() *ClientPool {
	return &ClientPool{
		syncMutex: new(sync.RWMutex),
		closed:    false,
	}
}

func (clientPool *ClientPool) takeClient(endpoint string, connectTimeout time.Duration,
	maxQueueSize int, socketTimeout time.Duration) (*ClientResource, error) {
	if clientPool.closed {
		tbase_log.TBaseLogger.Errorf("[CLIENT_POOL] client pool closed")
		return nil, NewTBaseClientInternalError("client pool closed")
	}

	var value interface{}
	var has bool
	startTime := time.Now().UnixNano()
	if value, has = clientPool.endpoint2ClientResourceMap.Load(endpoint); !has {
		asyncClient, err := clientPool.newAsyncClient(endpoint, connectTimeout, maxQueueSize, socketTimeout)
		if err != nil {
			tbase_log.TBaseLogger.Errorf("[CLIENT_POOL] new async client error. "+
				"endpoint: %v, error: %v", endpoint, err)
			return nil, err
		}
		newClientResource := NewClientResource(asyncClient, endpoint)

		existing, loaded := clientPool.endpoint2ClientResourceMap.LoadOrStore(endpoint, newClientResource)
		if loaded {
			value = existing
			newClientResource.Close()
		} else {
			value = newClientResource
		}
	}

	if clientResource, ok := value.(*ClientResource); ok {
		if atomic.LoadUint32(&clientResource.GetAsyncClient().ConnUsable) == 0 {
			tbase_log.TBaseLogger.Infof("[CLIENT_POOL] detect the client %v "+
				"has an net err, create a new one", endpoint)
			remainTime := int64(connectTimeout) - (time.Now().UnixNano() - startTime)
			if remainTime < 0 {
				return nil, NewTBaseClientTimeoutError("not enough time to re-create a client")
			}
			asyncClient, err := clientPool.newAsyncClient(endpoint, connectTimeout, maxQueueSize, socketTimeout)
			if err != nil {
				tbase_log.TBaseLogger.Errorf("[CLIENT_POOL] new async client error. endpoint: %v, error: %v", endpoint, err)
				return nil, err
			}
			newClientResource := NewClientResource(asyncClient, endpoint)

			var oldValue interface{}
			var needCloseClient *ClientResource = nil

			clientPool.syncMutex.Lock()
			oldValue, _ = clientPool.endpoint2ClientResourceMap.Load(endpoint)
			if oldClientResource, ok := oldValue.(*ClientResource); ok {
				if oldClientResource == nil || atomic.LoadUint32(&oldClientResource.GetAsyncClient().ConnUsable) == 0 {
					clientPool.endpoint2ClientResourceMap.Store(endpoint, newClientResource)
					needCloseClient = oldClientResource
					clientResource = newClientResource
				} else {
					needCloseClient = newClientResource
					clientResource = oldClientResource
				}
				clientPool.syncMutex.Unlock()
			} else {
				clientPool.syncMutex.Unlock()
				return nil, NewTBaseClientInternalError("close old client error, " +
					"can't convert to \"*ClientResource\"")
			}

			if needCloseClient != nil {
				needCloseClient.Close()
			}
		}

		return clientResource, nil

	} else {
		return nil, NewTBaseClientInternalError("can't convert to \"*ClientResource\"")
	}

}

func (clientPool *ClientPool) newAsyncClient(endpoint string, connectTimeout time.Duration,
	maxQueueSize int, socketTimeout time.Duration) (*async.AsyncClient, error) {

	asyncClient, err := async.NewAsyncClient(endpoint, connectTimeout, maxQueueSize, socketTimeout)
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[CLIENT_POOL] new async executor error. endpoint: %v, error: %v", endpoint, err)
		if netError, ok := err.(net.Error); ok {
			if netError.Timeout() {
				return nil, NewTBaseClientTimeoutError("initialize pool error. " + err.Error())
			} else {
				return nil, NewTBaseClientConnectionError("initialize pool error. " + err.Error())
			}
		} else {
			return nil, NewTBaseClientInternalError("initialize pool error. " + err.Error())
		}
	}

	return asyncClient, nil
}

func (clientPool *ClientPool) Close() {
	if !clientPool.closed {
		clientPool.closed = true
		clientPool.endpoint2ClientResourceMap.Range(func(key, value interface{}) bool {
			clientResource, ok := value.(*ClientResource)
			if !ok {
				tbase_log.TBaseLogger.Errorf("[CLIENT_POOL] can't convert to \"*ClientResource\"")
				// please see sync.Map#Range function details if confused by `return true`, continue iteration whatever happens
				return true
			}
			clientResource.Close()
			clientPool.endpoint2ClientResourceMap.Delete(key)
			return true
		})
		if tbase_log.TBaseLogger.GetLogLevel() > log.INFO {
			tbase_log.TBaseLogger.Infof("[CLIENT_POOL] client pool is closed")
		}
	}
}

func (clientPool *ClientPool) clearClient(endpoint string) error {
	if clientPool.closed {
		return NewTBaseClientInternalError("client pool closed")
	}

	if value, has := clientPool.endpoint2ClientResourceMap.Load(endpoint); has {
		clientPool.endpoint2ClientResourceMap.Delete(endpoint)
		if clientResource, ok := value.(*ClientResource); ok {
			clientResource.Close()
			return nil
		} else {
			return NewTBaseClientInternalError("can't convert map value to ClientResource")
		}
	} else {
		if tbase_log.TBaseLogger.GetLogLevel() >= log.INFO {
			tbase_log.TBaseLogger.Infof("[CLIENT_POOL] instance %v is not "+
				"currently part of this pool.", endpoint)
		}
		return nil
	}
}
