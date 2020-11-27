package tbasego

import (
	"sync"
	"time"

	"gitlab.alipay-inc.com/octopus/tbase-go-client/model"
	tbase_log "gitlab.alipay-inc.com/octopus/tbase-go-client/tloggger"
	"mosn.io/pkg/log"
)

type ClientHandler struct {
	pool         *ClientPool
	shardManager *ShardManager
	closed       bool
	wg           *sync.WaitGroup
}

func newClientHandler(pool *ClientPool, shardManager *ShardManager) *ClientHandler {
	clientHandler := &ClientHandler{
		pool:         pool,
		shardManager: shardManager,
		closed:       false,
		wg:           new(sync.WaitGroup),
	}

	clientHandler.wg.Add(1)

	go func() {
		defer clientHandler.wg.Done()
		clientHandler.startUpdateClientPoolLoop()
	}()

	return clientHandler
}

func (c *ClientHandler) takeClient(ta *model.TBaseAction, key []byte,
	connectTimeout time.Duration, maxQueueSize int,
	movedHints map[int]string, socketTimeout time.Duration) (*ClientResource, error) {
	endpoint, err := c.getEndpoint(key, movedHints)
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[CLIENT_HANDLER] get end point error, error: %v", err)
		return nil, err
	}
	ta.EndpointString = endpoint
	return c.pool.takeClient(endpoint, connectTimeout, maxQueueSize, socketTimeout)
}

func (c *ClientHandler) getEndpoint(key []byte, movedHints map[int]string) (string, error) {
	var endpoint string
	var err error
	var shardId int

	if len(movedHints) > 0 {
		shardId, err = c.shardManager.GetShardId(key)
		if err != nil {
			tbase_log.TBaseLogger.Errorf("[CLIENT_HANDLER] get shard id error, error: %v", err)
			return "", err
		}
		endpoint = movedHints[shardId]
	}

	// here means movedHints is empty or movedHints doesn't have hint for endpoint
	if len(endpoint) <= 0 {
		endpoint, err = c.shardManager.GetEndpointByKey(key)
		if err != nil {
			tbase_log.TBaseLogger.Errorf("[CLIENT_HANDLER] get endpoint by key error, error: %v", err)
			return "", err
		}
	}

	return endpoint, nil
}

func (c *ClientHandler) refresh() {
	c.shardManager.Refresh()
}

func (c *ClientHandler) close() {
	if !c.closed {
		c.closed = true
		//close(c.shardManager.RefreshClientPoolChannel)
		c.wg.Wait()
		c.pool.Close()
		if tbase_log.TBaseLogger.GetLogLevel() >= log.INFO {
			tbase_log.TBaseLogger.Infof("[CLIENT_HANDLER] client handler closed")
		}
	}
}

func (c *ClientHandler) startUpdateClientPoolLoop() {
	for {
		select {
		case endpoints, ok := <-c.shardManager.RefreshClientPoolChannel:
			if !ok {
				tbase_log.TBaseLogger.Infof("[CLIENT_HANDLER] RefreshClientPoolChannel is closed")
				return
			}
			oldEndpoints := endpoints.OldEndpoints
			newEndpoints := endpoints.NewEndpoints
			if len(oldEndpoints) <= 0 {
				continue
			}

			tbase_log.TBaseLogger.Infof("detect cluster layout change")

			time.Sleep(time.Duration(c.shardManager.ConnectionInfo.RedisTimeout*2) * time.Millisecond)

			var err error
			for oldEndpoint, _ := range oldEndpoints {
				if _, ok := newEndpoints[oldEndpoint]; !ok {
					tbase_log.TBaseLogger.Infof("[CLIENT_HANDLER] client %v no longer exists, prepare to close. ", oldEndpoint)
					if err = c.pool.clearClient(oldEndpoint); err != nil {
						// we can do nothing if clear client failed, just log
						tbase_log.TBaseLogger.Errorf("[CLIENT_HANDLER] clear client for endpoint %v error, error: %v",
							oldEndpoint, err)
					}
					tbase_log.TBaseLogger.Infof("[CLIENT_HANDLER] client %v closed", oldEndpoint)
				}
			}
		}
	}
}
