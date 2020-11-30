package sofaregistry

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	sofaregistryproto "gitlab.alipay-inc.com/sofa-go/sofa-registry-proto-go/proto"
)

var ErrSubscriberPeekTimeout = errors.New("subscriber: peek timeout")

// Scope is the realm definition of LDC architecture. Registry server
// will send different sub results according to different scopes.
type Scope string

// Enumeration of supported scopes.
const (
	ScopeZone       Scope = "zone"
	ScopeDataCenter Scope = "dataCenter"
	ScopeGlobal     Scope = "global"
)

type SubscribeContext struct {
	DataID  string `json:"data_id"`  // must
	Group   string `json:"group"`    // must
	AppName string `json:"app_name"` // optional
	Scope   Scope  `json:"scope"`    // optional, default to zone
	Ip      string `json:"ip"`       // optional, pod id

	// Resident: 是否常驻.
	// 使用场景：保证某个 dataid 永远不会被误 cancel, 比如说 DRM 使用 Registry 寻址场景。
	Resident bool `json:"resident"`
}

func NewSubscribeContext(dataID, group string) *SubscribeContext {
	return &SubscribeContext{
		DataID: dataID,
		Group:  group,
		Scope:  ScopeZone,
	}
}

func (s *SubscribeContext) SetResident(b bool) *SubscribeContext {
	s.Resident = b
	return s
}

func (s *SubscribeContext) IsResident() bool {
	return s.Resident
}

func (s *SubscribeContext) Equal(d *SubscribeContext) bool {
	// no appname: historical reason
	return s.DataID == d.DataID && s.Group == d.Group
}

func (s *SubscribeContext) GetAppName() string { return s.AppName }

func (s *SubscribeContext) SetAppName(name string) *SubscribeContext {
	s.AppName = name
	return s
}

func (s *SubscribeContext) GetIp() string { return s.Ip }

func (s *SubscribeContext) SetIp(ip string) *SubscribeContext {
	s.Ip = ip
	return s
}

func (s SubscribeContext) GetScope() Scope { return s.Scope }

func (s *SubscribeContext) SetScope(scope Scope) *SubscribeContext {
	s.Scope = scope
	return s
}

func (s *SubscribeContext) GetDataID() string { return s.DataID }
func (s *SubscribeContext) GetGroup() string  { return s.Group }

type segmentData struct {
	id      string
	data    map[string][]string // zone -> DataList
	version int64
}

type Subscriber struct {
	c       *Client
	ctx     *SubscribeContext
	handler Handler

	notifies   uint64
	notifiedCh chan struct{}
	waitLock   sync.Mutex

	// used for registry server to locate a publisher.
	id string

	// current version of sub data.
	version int64

	localZone string

	// canceled indicates if publisher is still in use.
	canceled bool

	// client will redo pub/sub task after re-establishing connection.
	lastTask *task

	// segments contains all received segments data from server.
	segments map[string]*segmentData

	// synchronizes pub/cancel request and recovery actions.
	sync.RWMutex
}

func NewSubscriber(c *Client, i *SubscribeContext) *Subscriber {
	return &Subscriber{
		c:        c,
		ctx:      i,
		id:       uuid.New().String(),
		version:  0,
		canceled: false,
		lastTask: nil,
		// use bufferd-channel to make sure do not lost channel event
		notifiedCh: make(chan struct{}, 1),
		segments:   make(map[string]*segmentData, 16),
	}
}

func (s *Subscriber) GetContext() *SubscribeContext {
	return s.ctx
}

func (s *Subscriber) GetVersion() int64 {
	s.RLock()
	defer s.RUnlock()
	return s.version
}

func (s *Subscriber) GetData() map[string][]string {
	s.RLock()
	defer s.RUnlock()
	return s.buildFullData()
}

func (s *Subscriber) GetLocalZone() string {
	s.RLock()
	defer s.RUnlock()
	return s.localZone
}

func (s *Subscriber) Peek(timeout time.Duration) (map[string][]string, string, error) {
	err := s.wait(timeout)
	if err != nil {
		return nil, "", err
	}

	s.RLock()
	defer s.RUnlock()

	return s.buildFullData(), s.localZone, nil
}

// wait waits the server push at least one times.
func (s *Subscriber) wait(timeout time.Duration) error {
	if atomic.LoadUint64(&s.notifies) > 0 { // already notify
		return nil
	}

	t := time.NewTimer(timeout)
	defer t.Stop()

	// at most one goroutine can hold the lock
	s.waitLock.Lock()
	defer s.waitLock.Unlock()

	if atomic.LoadUint64(&s.notifies) > 0 { // double check
		return nil
	}

	select {
	case <-s.notifiedCh:
		return nil
	case <-t.C:
		return ErrSubscriberPeekTimeout
	}
}

// notify notifies server push.
func (s *Subscriber) notify() {
	atomic.AddUint64(&s.notifies, 1)
	select {
	case s.notifiedCh <- struct{}{}:
	default:
	}
}

func (s *Subscriber) handleSegmentData(d *segmentData, localZone string) {
	s.Lock()

	oldseg, ok := s.segments[d.id]
	if ok {
		if d.version <= oldseg.version { // stale segment
			s.Unlock()
			return
		}
	}

	s.localZone = localZone

	// Override segment data.
	s.segments[d.id] = d

	// Build full user data and invoke callback.
	if s.handler != nil {
		s.handler.OnRegistryPush(s.ctx.DataID, s.buildFullData(), localZone)
	}

	s.Unlock()

	// do server push notify
	s.notify()
}

// buildFullData merges all segments and returns full user data.
func (s *Subscriber) buildFullData() map[string][]string {
	full := make(map[string][]string)
	for _, seg := range s.segments {
		for z, dl := range seg.data {
			full[z] = append(full[z], dl...)
		}
	}
	return full
}

func (s *Subscriber) ID() string {
	return s.id
}

func (s *Subscriber) Canceled() bool {
	s.Lock()
	defer s.Unlock()
	return s.canceled
}

func (s *Subscriber) Sub(h Handler) error {
	s.Lock()
	s.canceled = false
	s.handler = h
	s.version++
	msg := s.protoMsg(REGISTEREventType)
	t := &task{
		id:      s.id,
		version: s.version,
		class:   SUBSCRIBEPbClass,
		req:     msg,
		res:     new(sofaregistryproto.RegisterResponsePb),
	}
	s.lastTask = t
	s.Unlock()

	// sanity store: avoid orphan subscriber
	s.c.subscribers.Store(s.id, s)

	err := s.c.enqueueTask(t)
	if err != nil {
		return err
	}

	s.c.metrics.addSubscriberRegister()
	return nil
}

// Cancel unregisters the publish operation from registry server.
// DO NOT keep using thus publisher after canceling.
func (s *Subscriber) Cancel() error {
	// cannot be canceled forever
	if s.GetContext().IsResident() {
		return nil
	}

	s.Lock()
	defer s.Unlock()

	s.version++
	msg := s.protoMsg(UNREGISTEREventType)
	t := &task{
		version: s.version,
		class:   SUBSCRIBEPbClass,
		req:     msg,
		res:     new(sofaregistryproto.RegisterResponsePb),
	}

	s.c.metrics.addSubscriberUnregister()
	err := s.c.enqueueTask(t)
	if err != nil {
		return err
	}

	s.canceled = true
	s.lastTask = nil
	s.c.subscribers.Delete(s.id)

	return nil
}

func (s *Subscriber) protoMsg(typ string) *sofaregistryproto.SubscriberRegisterPb {
	return &sofaregistryproto.SubscriberRegisterPb{
		Scope: string(s.ctx.Scope),
		BaseRegister: &sofaregistryproto.BaseRegisterPb{
			InstanceId: s.c.config.instanceID,
			Zone:       s.c.config.zone,
			DataId:     s.ctx.DataID,
			Group:      s.ctx.Group,
			AppName:    s.ctx.AppName,
			Ip:         s.ctx.Ip,
			RegistId:   s.ID(),
			ClientId:   s.ID(),
			EventType:  typ,
			Version:    s.version,
			Timestamp:  time.Now().UnixNano() / 1000 / 1000,
			Attributes: s.c.config.GetSignature(),
		},
	}
}

func (s *Subscriber) resub() error {
	s.Lock()
	defer s.Unlock()

	if s.canceled || s.lastTask == nil {
		return nil
	}

	s.version++
	msg := s.protoMsg(REGISTEREventType)
	t := &task{
		id:      s.id,
		version: s.version,
		class:   SUBSCRIBEPbClass,
		req:     msg,
		res:     new(sofaregistryproto.RegisterResponsePb),
	}
	s.lastTask = t

	return s.c.enqueueTask(s.lastTask)
}
