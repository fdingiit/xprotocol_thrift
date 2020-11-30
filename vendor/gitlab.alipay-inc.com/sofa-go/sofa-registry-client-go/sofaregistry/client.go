package sofaregistry

import (
	"errors"
	"net"
	"sync"
	"time"

	"gitlab.alipay-inc.com/sofa-go/sofa-registry-client-go/sofaregistry/queue"
	sofaregistryproto "gitlab.alipay-inc.com/sofa-go/sofa-registry-proto-go/proto"

	"github.com/gogo/protobuf/proto"
	"gitlab.alipay-inc.com/sofa-go/sofa-logger-go/sofalogger"
)

//go:generate syncmap -pkg sofaregistry -o publisher_generated.go -name PublisherMap map[string]*Publisher
//go:generate syncmap -pkg sofaregistry -o subscriber_generated.go -name SubscriberMap map[string]*Subscriber

const (
	defaultResetTimeout = 5 * time.Second
)

type Client struct {
	sync.Mutex
	publishers  PublisherMap
	subscribers SubscriberMap
	logger      sofalogger.Logger
	transport   Transport
	config      *Config
	queue       queue.Queue
	metrics     *Metrics
	redialCh    chan struct{}
}

func New(options ...ClientOptionSetter) (*Client, error) {
	c := &Client{}
	for i := range options {
		options[i].set(c)
	}

	if err := c.polyfill(); err != nil {
		return nil, err
	}

	c.run()

	return c, nil
}

func (c *Client) Publish(ctx *PublishContext, dl []string) error {
	c.Lock()
	defer c.Unlock()

	var (
		p     *Publisher
		found = false
	)

	// try find existent publisher
	c.publishers.Range(func(key string, publisher *Publisher) bool {
		pctx := publisher.GetContext()
		if pctx.Equal(ctx) {
			found = true
			p = publisher
			return false
		}
		return true
	})

	if p == nil {
		p = NewPublisher(c, ctx)
	}

	if err := p.Pub(dl); err != nil {
		return err
	}

	if !found {
		c.publishers.Store(p.id, p)
	}

	return nil
}

func (c *Client) UnPublish(ctx *PublishContext) error {
	c.Lock()
	defer c.Unlock()

	// unpublish all publishers with same dataid
	c.publishers.Range(func(key string, p *Publisher) bool {
		pctx := p.GetContext()
		if !pctx.Equal(ctx) {
			return true
		}

		if err := p.Cancel(); err != nil {
			c.logger.Errorf("failed to unpublish dataID=%s group=%s appname=%s: %s",
				ctx.DataID, ctx.Group, ctx.AppName, err)
			return true
		}
		// remove the publisher
		c.publishers.Delete(key)
		return true
	})

	return nil
}

// HACK: A temporary solution to make sure only one publisher with the same dataid
func (c *Client) CreatePublisher(ctx *PublishContext) *Publisher {
	c.Lock()
	defer c.Unlock()

	var p *Publisher
	c.publishers.Range(func(key string, publisher *Publisher) bool {
		if publisher.GetContext().Equal(ctx) {
			p = publisher
			return false
		}
		return true
	})

	if p == nil {
		p = NewPublisher(c, ctx)
	}

	c.publishers.Store(p.id, p)
	return p
}

// HACK: A temporary solution to make sure only one subscriber with the same dataid
func (c *Client) CreateSubscriber(ctx *SubscribeContext) *Subscriber {
	c.Lock()
	defer c.Unlock()

	var s *Subscriber
	c.subscribers.Range(func(key string, subscriber *Subscriber) bool {
		if subscriber.GetContext().Equal(ctx) {
			s = subscriber
			return false
		}
		return true
	})

	if s == nil {
		s = NewSubscriber(c, ctx)
	}

	c.subscribers.Store(s.id, s)
	return s
}

func (c *Client) GetPublishers(fn func(id string, publisher *Publisher) bool) {
	c.publishers.Range(fn)
}

func (c *Client) GetSubscribers(fn func(id string, subscriber *Subscriber) bool) {
	c.subscribers.Range(fn)
}

func (c *Client) GetMetrics() *Metrics {
	return c.metrics
}

func (c *Client) polyfill() error {
	if c.logger == nil {
		c.logger = sofalogger.StdoutLogger
	}

	if c.redialCh == nil {
		c.redialCh = make(chan struct{})
	}

	if c.metrics == nil {
		c.metrics = NewMetrics()
	}

	if c.transport == nil {
		return errors.New("sofaregistry: transport cannot be nil")
	}

	if c.config == nil {
		return errors.New("sofaregistry: config cannot be nil")
	}

	if c.queue == nil {
		c.queue = queue.NewFIFOWithSize(1024)
	}

	return nil
}

func (c *Client) send() {
	for {
		item, err := c.queue.Pop()
		if err != nil {
			c.logger.Errorf("failed to get item from queue but continue: %+v", err.Error())
			continue
		}

		t, ok := item.(*task)
		if !ok {
			panic("failed to type casting")
		}

		if err = c.handleTask(t); err != nil {
			c.logger.Errorf("failed to handle task: %+v", err.Error())
		}
	}
}

func (c *Client) redial() {
	c.transport.OnRedial(func(conn net.Conn) {
		c.metrics.addRedial()

		c.logger.Infof("client redial success conn=%s->%s to republish and resubscribe",
			conn.LocalAddr(), conn.RemoteAddr())

		// notify redial event
		select {
		case c.redialCh <- struct{}{}:
		default:
		}

		// sleep 500ms to wait transport to ready
		time.Sleep(500 * time.Millisecond)

		c.publishers.Range(func(key string, publisher *Publisher) bool {
			if err := publisher.repub(); err != nil {
				c.logger.Errorf("failed to republish: %v", err)
			}
			return true
		})

		// wait 1s sofaregistry server to process publisher
		time.Sleep(1 * time.Second)

		c.subscribers.Range(func(key string, subscriber *Subscriber) bool {
			if err := subscriber.resub(); err != nil {
				c.logger.Errorf("failed to resubscribe: %v", err)
			}
			return true
		})
	})
}

func (c *Client) recv() {
	err := c.transport.OnRecv(
		func(err error, req proto.Message) {
			c.metrics.addServerPush()

			if err != nil {
				c.logger.Errorf("failed to receive request from transport: %v", err)
				return
			}

			received, ok := req.(*sofaregistryproto.ReceivedDataPb)
			if !ok {
				c.logger.Errorf("received an unexpected request")
				return
			}

			c.logger.Infof("server push request=%s", prettyReceivedDataPb(received))

			data := make(map[string][]string)
			for z, db := range received.Data {
				data[z] = dataBoxesPb2DataList(db.Data)
			}

			for _, id := range received.SubscriberRegistIds {
				subscriber, ok := c.subscribers.Load(id)
				if !ok {
					c.logger.Infof("stale subscriber: %s", id)
					continue
				}

				subscriber.handleSegmentData(
					&segmentData{
						id:      received.Segment,
						data:    data,
						version: received.Version,
					},
					received.LocalZone,
				)
			}
		})
	if err != nil { // should never happen
		c.logger.Errorf("received unexpected error: %+v", err)
	}
}

func (c *Client) run() {
	go c.recv()
	go c.send()
	go c.redial()
}

func (c *Client) enqueueTask(t *task) error {
	return c.queue.Push(t)
}

func (c *Client) handleTask(t *task) error {
	if c.handle(t) {
		t.Done()
		c.metrics.addSuccessTask()
		return nil
	}

	c.metrics.addFailureTask()
	t.addFailures()
	return c.queue.Push(t)
}

func (c *Client) handle(t *task) bool {
	started := time.Now()

	// try delay when the task always be failed
	t.tryDelay()

	err := c.invoke(t.class, t.req, t.res)
	if err != nil {
		c.logger.Errorf("failed to invoke transport: %v", err)
		return false
	}

	switch req := t.req.(type) {
	case *sofaregistryproto.PublisherRegisterPb:
		return c.loggingPublisherTask(started, t, req)

	case *sofaregistryproto.SubscriberRegisterPb:
		return c.loggingSubscriberTask(started, t, req)

	default:
		c.logger.Errorf("unknown task: %+v", t.req)
	}

	return false
}

func (c *Client) invoke(class string, req proto.Message, res proto.Message) error {
	return c.transport.Send(class, req, res)
}

// nolint
func (c *Client) loggingSubscriberTask(started time.Time, t *task, req *sofaregistryproto.SubscriberRegisterPb) bool {
	res, ok := t.res.(*sofaregistryproto.RegisterResponsePb)
	if !ok {
		return false
	}

	if !res.Success {
		c.logger.Errorf("subscriber register failure queuue=%d response=%+v elapsed=%s",
			c.queue.Len(), res.String(), time.Since(started).String())
		return false
	}

	c.logger.Infof("subscriber register success queue=%d base=%+v elapsed=%s", c.queue.Len(), req.String(), time.Since(started).String())

	return true
}

// nolint
func (c *Client) loggingPublisherTask(started time.Time, t *task, req *sofaregistryproto.PublisherRegisterPb) bool {
	res, ok := t.res.(*sofaregistryproto.RegisterResponsePb)
	if !ok {
		return false
	}

	if !res.Success {
		c.logger.Errorf("publisher register failure queue=%d response=%+v elapsed=%s", c.queue.Len(), res.String(),
			time.Since(started).String())
		return false
	}

	c.logger.Infof("publisher register success queue=%d base=%+v elapsed=%s", c.queue.Len(),
		req.String(), time.Since(started).String())

	return true
}

// Reset resets all publishers and subscribers then close connection until redial connection or timeout
// for registry server GC the connection
func (c *Client) Reset() {
	c.logger.Infof("reset the client")
	c.subscribers.Range(func(id string, subscriber *Subscriber) bool {
		// skip resident subscriber
		if subscriber.GetContext().IsResident() {
			return true
		}

		if err := subscriber.Cancel(); err != nil {
			c.logger.Errorf("subscriber cancel failed: %+v", err)
			return true
		}

		return true
	})

	c.publishers.Range(func(id string, publisher *Publisher) bool {
		if err := publisher.Cancel(); err != nil {
			c.logger.Errorf("publisher cancel failed: %+v", err)
		}
		return true
	})

	if err := c.transport.Close(); err != nil {
		c.logger.Errorf("failed to close transport: %+v", err)
	}

	// wait redial event
	select {
	case <-c.redialCh:
	case <-time.After(defaultResetTimeout):
	}
}
