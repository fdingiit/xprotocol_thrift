package async

import (
	"time"

	"gitlab.alipay-inc.com/octopus/radix"
	tbase_log "gitlab.alipay-inc.com/octopus/tbase-go-client/tloggger"
)

// this async pool maintains one connection only
type AsyncPool struct {
	addr              string
	connectionTimeout time.Duration
	socketTimeout     time.Duration
	conn              *radix.IoErrConn
}

func NewAsyncPool(addr string, connectionTimeout time.Duration, socketTimeout time.Duration) (*AsyncPool, error) {
	conn, err := radix.Dial("tcp", addr, radix.DialConnectTimeout(connectionTimeout),
		radix.DialWriteTimeout(socketTimeout), radix.DialReadTimeout(socketTimeout))
	if err != nil {
		tbase_log.TBaseLogger.Errorf("[ASYNC_POOL] dial to %v error, error %v", addr, err)
		return nil, err
	}

	return &AsyncPool{addr: addr, connectionTimeout: connectionTimeout, socketTimeout: socketTimeout, conn: &radix.IoErrConn{Conn: conn}}, nil
}

func (asyncPool *AsyncPool) Get() *radix.IoErrConn {
	return asyncPool.conn
}

func (asyncPool *AsyncPool) Close() {
	if err := asyncPool.conn.Close(); err != nil {
		// we can nothing if close failed, just log
		tbase_log.TBaseLogger.Errorf("[ASYNC_POOL] close connection for %v error, error %v", asyncPool.addr, err)
	}
}
