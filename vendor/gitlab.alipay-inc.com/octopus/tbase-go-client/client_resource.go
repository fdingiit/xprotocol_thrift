package tbasego

import (
	"gitlab.alipay-inc.com/octopus/tbase-go-client/async"
	tbase_log "gitlab.alipay-inc.com/octopus/tbase-go-client/tloggger"
)

/**

 */
type ClientResource struct {
	asyncClient *async.AsyncClient
	endpoint    string
	closed      bool
}

func NewClientResource(asyncClient *async.AsyncClient, endpoint string) *ClientResource {
	return &ClientResource{asyncClient: asyncClient, endpoint: endpoint, closed: false}
}

func (clientResource *ClientResource) GetEndpoint() string {
	return clientResource.endpoint
}

func (clientResource *ClientResource) GetAsyncClient() *async.AsyncClient {
	return clientResource.asyncClient
}

func (clientResource *ClientResource) Close() {
	if !clientResource.closed {
		clientResource.closed = true
		clientResource.asyncClient.Close()
		tbase_log.TBaseLogger.Infof("[CLIENT_RESOURCE] client resource is closed")
	}
}
