package sofaregistry

import (
	"sync"
	"time"

	"github.com/google/uuid"
	sofaregistryproto "gitlab.alipay-inc.com/sofa-go/sofa-registry-proto-go/proto"
)

type PublishContext struct {
	DataID  string `json:"data_id"`  // must
	Group   string `json:"gruop"`    // must
	AppName string `json:"app_name"` // optional
	Ip      string `json:"ip"`       // optional, pod id
}

func NewPublishContext(dataID, group string) *PublishContext {
	return &PublishContext{DataID: dataID, Group: group}
}

func (s *PublishContext) Equal(d *PublishContext) bool {
	// no appname: historical reason
	return s.DataID == d.DataID && s.Group == d.Group
}

func (s *PublishContext) SetAppName(name string) *PublishContext {
	s.AppName = name
	return s
}

func (s *PublishContext) SetIp(ip string) *PublishContext {
	s.Ip = ip
	return s
}

func (s *PublishContext) GetAppName() string { return s.AppName }
func (s *PublishContext) GetDataID() string  { return s.DataID }
func (s *PublishContext) GetGroup() string   { return s.Group }
func (s *PublishContext) GetIp() string      { return s.Ip }

type Publisher struct {
	c   *Client
	ctx *PublishContext

	// used for registry server to locate a publisher.
	id string

	// current version of pub data.
	version int64

	// client will redo pub/sub task after re-establishing connection.
	lastTask *task

	// canceled indicates if publisher is still in use.
	canceled bool

	// synchronizes pub/cancel request and recovery actions.
	sync.Mutex
}

// newPublisher returns a pointer to new initialized publisher.
func NewPublisher(c *Client, ctx *PublishContext) *Publisher {
	return &Publisher{
		c:        c,
		ctx:      ctx,
		id:       uuid.New().String(),
		version:  1,
		canceled: false,
	}
}

func (p *Publisher) GetContext() *PublishContext {
	return p.ctx
}

func (p *Publisher) GetVersion() int64 {
	p.Lock()
	version := p.version
	p.Unlock()
	return version
}

func (p *Publisher) IsDone() bool {
	p.Lock()
	defer p.Unlock()
	if p.lastTask != nil {
		return p.lastTask.IsDone()
	}
	return false
}

// Pub asynchronously sends new data list to registry server.
// NOTE DO NOT use concurrently!
func (p *Publisher) Pub(dl []string) error {
	p.Lock()
	p.version++
	msg := p.protoMsg(dl, REGISTEREventType)
	t := &task{
		version: p.version,
		id:      p.id,
		class:   PUBLISHPbClass,
		req:     msg,
		res:     new(sofaregistryproto.RegisterResponsePb),
	}
	p.lastTask = t
	p.canceled = false
	p.Unlock()

	// sanity store: avoid orphan publisher
	p.c.publishers.Store(p.id, p)

	if err := p.c.enqueueTask(t); err != nil {
		return err
	}

	p.c.metrics.addPublisherRegister()
	return nil
}

// Cancel unregisters the publish operation from registry server.
// DO NOT keep using thus publisher after canceling.
func (p *Publisher) Cancel() error {
	p.Lock()
	defer p.Unlock()

	p.version++
	msg := p.protoMsg(nil, UNREGISTEREventType)
	t := &task{
		version: p.version,
		id:      p.id,
		class:   PUBLISHPbClass,
		req:     msg,
		res:     new(sofaregistryproto.RegisterResponsePb),
	}

	p.c.metrics.addPublisherUnregister()
	err := p.c.enqueueTask(t)
	if err != nil {
		return err
	}

	p.canceled = true
	p.lastTask = nil
	p.c.publishers.Delete(p.id)

	return nil
}

func (p *Publisher) ID() string {
	return p.id
}

func (p *Publisher) Canceled() bool {
	p.Lock()
	defer p.Unlock()
	return p.canceled
}

func (p *Publisher) protoMsg(dl []string, typ string) *sofaregistryproto.PublisherRegisterPb {
	return &sofaregistryproto.PublisherRegisterPb{
		DataList: dataList2DataBoxesPb(dl),
		BaseRegister: &sofaregistryproto.BaseRegisterPb{
			InstanceId: p.c.config.instanceID,
			Zone:       p.c.config.zone,
			DataId:     p.ctx.DataID,
			Group:      p.ctx.Group,
			AppName:    p.ctx.AppName,
			Ip:         p.ctx.Ip,
			RegistId:   p.ID(),
			ClientId:   p.ID(),
			EventType:  typ,
			Version:    p.version,
			Timestamp:  time.Now().UnixNano() / 1000 / 1000,
			Attributes: p.c.config.GetSignature(),
		},
	}
}

func (p *Publisher) repub() error {
	p.Lock()
	defer p.Unlock()

	if p.canceled || p.lastTask == nil {
		return nil
	}

	switch x := p.lastTask.req.(type) {
	case *sofaregistryproto.PublisherRegisterPb:
		p.version++
		msg := p.protoMsg(dataBoxesPb2DataList(x.DataList), REGISTEREventType)
		t := &task{
			version: p.version,
			id:      p.id,
			class:   PUBLISHPbClass,
			req:     msg,
			res:     new(sofaregistryproto.RegisterResponsePb),
		}
		p.lastTask = t
	}

	return p.c.enqueueTask(p.lastTask)
}
