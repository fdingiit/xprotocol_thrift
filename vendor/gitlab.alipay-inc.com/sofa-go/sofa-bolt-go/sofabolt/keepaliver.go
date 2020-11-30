package sofabolt

import (
	"context"
	"sync"
	"time"

	"gitlab.alipay-inc.com/sofa-go/sofa-logger-go/sofalogger"
)

type KeepAliverOptions struct {
	Context           context.Context
	MaxClientUsed     int
	MinClientInPool   int
	HeartbeatInterval time.Duration
	HeartbeatTimeout  time.Duration
	CleanupInterval   time.Duration
	CleanupMaxChecks  int
}

type KeepAliver struct {
	logger  sofalogger.Logger
	options *KeepAliverOptions
	raw     PoolMap
	tls     PoolMap
	// dying holds the clients will dies
	dying sync.Map
}

func NewKeepAliver(o *KeepAliverOptions, logger sofalogger.Logger) (*KeepAliver, error) {
	ka := &KeepAliver{
		logger:  logger,
		options: o,
	}

	if err := ka.polyfill(); err != nil {
		return nil, err
	}

	go ka.doCleanup(o.Context)
	go ka.doHeartbeat(o.Context)

	return ka, nil
}

// nolint
func (ca *KeepAliver) polyfill() error {
	if ca.options.CleanupInterval == 0 {
		ca.options.CleanupInterval = 10 * time.Second
	}

	if ca.options.CleanupMaxChecks == 0 {
		ca.options.CleanupMaxChecks = 15
	}

	if ca.options.HeartbeatInterval == 0 {
		ca.options.HeartbeatInterval = 30 * time.Second
	}

	if ca.options.HeartbeatTimeout == 0 {
		ca.options.HeartbeatTimeout = 5 * time.Second
	}

	if ca.options.Context == nil {
		ca.options.Context = context.TODO()
	}

	return nil
}

func (ca *KeepAliver) doCleanup(ctx context.Context) {
	cleanupInterval := ca.options.CleanupInterval

	for {
		time.Sleep(cleanupInterval)
		select {
		case <-ctx.Done():
			ca.logger.Infof("shutdown bolt keepaliver cleanup")
			return
		default:
		}

		ca.dying.Range(func(key, value interface{}) bool {
			client, ok := key.(*Client)
			if !ok {
				panic("failed to type casting")
			}

			n, ok := value.(int)
			if !ok {
				panic("failed to type casting")
			}

			if n >= ca.options.CleanupMaxChecks {
				err := client.Close()
				ca.logger.Infof("close dying client (>= max checks) conn=%+v err=%=v",
					client.GetConn(), err)

			} else {
				if ref := client.GetMetrics().GetReferences(); ref > 0 {
					ca.logger.Infof("Skip close client ref=%d conn=%s", ref, client.GetConn())
					ca.dying.Store(client, n+1)

				} else {
					err := client.Close()
					ca.logger.Infof("Skip close client ref=%d conn=%s error=%+v", ref, client.GetConn(), err)
				}
			}

			return true
		})
	}
}

func (ca *KeepAliver) doHeartbeat(ctx context.Context) {
	heartbeatInterval := ca.options.HeartbeatInterval
	heartbeatTimeout := ca.options.HeartbeatTimeout

	req := AcquireRequest()
	res := AcquireResponse()
	req.SetProto(ProtoBOLTV1)
	req.SetCMDCode(CMDCodeBOLTHeartbeat)
	defer func() {
		ReleaseRequest(req)
		ReleaseResponse(res)
	}()

	clientChecker := func(scheme string, address string, p *Pool, t time.Time) bool {
		p.Iterate(func(client *Client) {
			callback := func(err error, ctx *InvokeContext) {
				ca.logger.Debugf("bolt heartbeat response scheme=%s address=%s error=%+v",
					scheme, address, err)

				if err != nil {
					ca.raw.Delete(address)
					ca.logger.Infof("bolt heartbeat failed scheme=%s address=%s error=%+v",
						scheme, address, err)
					return
				}
			}

			if t.Unix()-client.GetMetrics().GetLasted() >= int64(heartbeatInterval.Seconds()) {
				if err := client.DoCallbackTimeout(req,
					ClientCallbackerFunc(callback), heartbeatTimeout); err != nil {
					ca.logger.Errorf("failed to send heartbeat: %+v", err)
				}
			}
		})

		return true
	}

	timer := time.NewTicker(heartbeatInterval)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			ca.logger.Infof("shutdown bolt clientaliver heartbeat")
			return
		case t := <-timer.C:
			ca.logger.Infof("try send bolt heartbeat")
			ca.raw.Range(func(addr string, pool *Pool) bool {
				return clientChecker("raw", addr, pool, t)
			})

			ca.raw.Range(func(addr string, pool *Pool) bool {
				return clientChecker("tls", addr, pool, t)
			})
		}
	}
}

func (ca *KeepAliver) Put(tls bool, force bool, addr string, client *Client) bool {
	if tls {
		return ca.put(&ca.tls, force, addr, client)
	}
	return ca.put(&ca.raw, force, addr, client)
}

func (ca *KeepAliver) put(m *PoolMap, force bool, addr string, client *Client) bool {
	var (
		loaded bool
		actual *Pool
		p      = NewPool()
	)

	p.Push(client)

	actual, loaded = m.LoadOrStore(addr, p)
	if loaded { // One guy win so release the loser
	}

	if actual.Size() >= 2 && ca.options.MaxClientUsed > 0 &&
		client.GetMetrics().GetUsed() >= int64(ca.options.MaxClientUsed) {
		return false
	}

	if force || (ca.options.MinClientInPool > 0 && actual.Size() <= ca.options.MinClientInPool) {
		if loaded {
			actual.Push(client)
		}
		return true
	}

	return false
}

func (t *KeepAliver) Get(tls bool, addr string) (*Client, bool) {
	if tls {
		return t.get(&t.tls, addr)
	}
	return t.get(&t.raw, addr)
}

func (t *KeepAliver) get(m *PoolMap, addr string) (*Client, bool) {
	p, ok := m.Load(addr)
	if !ok {
		return nil, false
	}

	var c *Client

	for {
		c, ok = p.Get()
		if !ok {
			return nil, false
		}

		if p.Size() >= 2 &&
			t.options.MaxClientUsed > 0 &&
			c.GetMetrics().GetUsed() >= int64(t.options.MaxClientUsed) {
			t.del(m, addr, c)
			t.GracefullyClose(c)
			continue
		}
		break
	}

	return c, true
}

func (ca *KeepAliver) Del(tls bool, address string, client *Client) bool {
	if tls {
		return ca.del(&ca.tls, address, client)
	}
	return ca.del(&ca.raw, address, client)
}

func (k *KeepAliver) del(m *PoolMap, addr string, client *Client) bool {
	p, ok := m.Load(addr)
	if !ok {
		return false
	}

	p.Delete(client)

	return true
}

func (k *KeepAliver) GracefullyClose(client *Client) {
	k.logger.Infof("try to gracefully close client used=%d lasted=%d ref=%d conn=%+v",
		client.GetMetrics().GetUsed(),
		client.GetMetrics().GetLasted(),
		client.GetMetrics().GetReferences(),
		client.GetConn())

	if ref := client.GetMetrics().GetReferences(); ref > 0 {
		k.dying.Store(client, 1)
	} else {
		err := client.Close()
		if err == nil {
			k.logger.Infof("direct close refless client")
		} else {
			k.logger.Infof("direct close refless client: %s", err.Error())
		}
	}
}
