package sofadrm

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"gitlab.alipay-inc.com/sofa-go/sofa-drm-client-go/sofadrm/model"

	"github.com/gogo/protobuf/proto"
	"github.com/google/uuid"
	"gitlab.alipay-inc.com/sofa-go/sofa-logger-go/sofalogger"
)

const (
	defaultDRMPort           = 9880
	SubscriberRegReqPbClass  = "com.alipay.drm.client.api.model.pb.SubscriberRegReqPb"
	SubscriberRegResultClass = "com.alipay.drm.client.api.model.pb.SubscriberRegResultPb"
	AttributeSetRequestClass = "com.alipay.drm.client.api.model.pb.AttributeSetRequestPb"
	AttributeGetRequestClass = "com.alipay.drm.client.api.model.pb.AttributeGetRequestPb"
	HeartbeatRequestPbClass  = "com.alipay.drm.client.api.model.pb.HeartbeatRequestPb"
	DefaultDRMValueVersion   = -1
)

type Client struct {
	transport Transport
	logger    sofalogger.Logger
	cache     model.LocalValueCacheMap
	listeners MutListenerMap
	config    *Config
	metrics   *Metrics
	lock      sync.Mutex
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

func (c *Client) polyfill() error {
	if c.logger == nil {
		c.logger = sofalogger.StdoutLogger
	}

	if c.transport == nil {
		return errors.New("sofadrm: transport cannot be nil")
	}

	if c.metrics == nil {
		c.metrics = &Metrics{}
	}

	return nil
}

func (c *Client) GetMetrics() *Metrics {
	return c.metrics
}

func (c *Client) GetValue(dataID string) (string, int, error) {
	// Load or store the dummy listener
	_, ok := c.listeners.Load(dataID)
	if !ok {
		c.listeners.Store(dataID, buildMutListener(dummyListenerFunc))
	}

	local, ok := c.cache.Load(dataID)
	if ok && local.DrmValue.Version != DefaultDRMValueVersion {
		return local.DrmValue.Value, local.DrmValue.Version, nil
	} else {
		// Store the initialized version even if it's failed
		// client will send heartbeat at interval
		c.cache.Store(dataID, model.LocalValueCache{
			DataId: dataID,
			DrmValue: model.DrmValue{
				Version: DefaultDRMValueVersion,
			},
		})
	}

	return c.fetchAndUpdateLocalValue(dataID, -1)
}

func (c *Client) AddListener(dataID string, ln Listener) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	m, ok := c.listeners.Load(dataID)
	if ok { // if already register then *overwrite* it
		m.OverWrite(ln)
		return nil
	}

	// Register the listener in order to reregister at interval
	c.listeners.Store(dataID, buildMutListener(ln))

	// Try send subscribe request even if it's failed
	// client will resubscribe until success
	err := c.doSubscribe(dataID)
	if err != nil {
		err = fmt.Errorf("failed to send subscribe request: %s", err.Error())
		c.logger.Errorf("%s", err.Error())
		return err
	}

	return nil
}

func (c *Client) doSubscribe(dataID string) error {
	subscriberRegReq := &model.SubscriberRegReq{
		Zone:       c.config.GetZone(),
		DataId:     dataID,
		InstanceId: c.config.GetInstanceID(),
		Uuid:       uuid.New().String(),
		Attributes: make(map[string]string),
		Profile:    c.config.GetProfile(),
	}

	if c.config.accessKey != "" &&
		c.config.secretKey != "" &&
		c.config.instanceID != "" {
		sm := getSignatureMap(c.config.accessKey, c.config.secretKey, c.config.instanceID)
		for k, v := range sm {
			subscriberRegReq.Attributes[k] = v
		}
	}

	req := model.ConvertToSubscriberRegReqPb(subscriberRegReq)

	return c.transport.Send(SubscriberRegReqPbClass, req, nil)
}

func (c *Client) run() {
	go c.recv()
	go c.redial()
	go c.doHeartbeatAndRegisterCheck()
}

func (c *Client) doHeartbeatAndRegisterCheck() {
	for {
		time.Sleep(c.config.heartbeatInterval)
		c.heartbeat()
		c.register()
	}
}

func (c *Client) register() {
	c.listeners.Range(func(dataID string, ln *mutlistener) bool {
		if ln.IsRegistered() {
			return true
		}

		if err := c.doSubscribe(dataID); err != nil {
			c.logger.Errorf("failed to subscribe dataID=%s: %v", dataID, err)
		}

		return true
	})
}

func (c *Client) heartbeat() {
	c.GetMetrics().addHeartbeat()

	hb := &model.HeartbeatRequest{
		Zone: c.config.zone,
		// ClientIp:   rc.clientIp,
		InstanceId: c.config.instanceID,
		Profile:    c.config.GetProfile(),
	}

	versions := make(map[string]int32, 64)
	ackVersions := make(map[string]int32, 64)
	c.cache.Range(func(dataID string, value model.LocalValueCache) bool {
		versions[dataID] = int32(value.DrmValue.Version)
		ackVersions[dataID] = int32(value.DrmAck.Version)
		return true
	})

	hb.VersionMap = versions
	hb.AckVersionMap = ackVersions

	var (
		preq = model.ConvertToHeartbeatRequestPb(hb)
		res  = new(model.HeartbeatResponsePb)
	)
	res.DiffMap = make(map[string]int32, 16)

	if err := c.transport.Send(HeartbeatRequestPbClass, preq, res); err != nil {
		c.logger.Errorf("failed to send heartbeat: %v", err)
		return
	}

	for dataID, version := range res.DiffMap {
		if _, _, err := c.fetchAndUpdateLocalValue(dataID, int(version)); err != nil {
			c.logger.Errorf("failed to update local value: %v", err)
		}
	}

	c.logger.Infof("send heartbeat success with diffmap=%+v", res.DiffMap)
}

func (c *Client) redial() {
	c.transport.OnRedial(func(conn net.Conn) {
		c.metrics.addRedial()

		c.logger.Infof("client redial success conn=%s->%s resubscribe",
			conn.LocalAddr(), conn.RemoteAddr())

		time.Sleep(1 * time.Second)

		c.listeners.Range(func(dataID string, ln *mutlistener) bool {
			ln.MarkUnRegistered()

			err := c.doSubscribe(dataID)
			if err != nil {
				c.logger.Errorf("failed to resubscribe dataID=%s", dataID)
			}

			return true
		})
	})
}

func (c *Client) recv() {
	err := c.transport.OnRecv(
		func(err error, class string, req, res proto.Message) {
			c.metrics.addServerPush()

			if err != nil {
				c.logger.Errorf("failed to receive request from transport: %v", err)
				return
			}

			switch class {
			case SubscriberRegResultClass:
				mreq, ok := req.(*model.SubscriberRegResultPb)
				if !ok {
					return
				}
				c.handleSubscriberRegResult(mreq)

			case AttributeGetRequestClass:
				mreq, ok := req.(*model.AttributeGetRequestPb)
				if !ok {
					return
				}
				mres, ok := req.(*model.AttributeGetResponse)
				if !ok {
					return
				}

				c.handleAttributeGetRequest(mreq, mres)

			case AttributeSetRequestClass:
				mreq, ok := req.(*model.AttributeSetRequestPb)
				if !ok {
					return
				}
				c.handleAttributeSetRequest(mreq)

			default:
				c.logger.Errorf("received an unexpected request: %s", class)
			}
		})
	if err != nil { // should never happen
		c.logger.Errorf("received unexpected error: %+v", err)
	}
}

func (c *Client) handleAttributeSetRequest(req *model.AttributeSetRequestPb) {
	c.logger.Infof("drm push setrequest req=%+v", req.String())
	c.OnDRMPush(req.DataId, req.Value)
}

func (c *Client) handleAttributeGetRequest(req *model.AttributeGetRequestPb, res *model.AttributeGetResponse) {
	c.logger.Infof("drm push getrequest req=%+v", req.String())
	v, ok := c.cache.Load(req.DataId)
	if !ok {
		return
	}

	res.SetDataString(v.DrmValue.Value)
}

func (c *Client) handleSubscriberRegResult(res *model.SubscriberRegResultPb) {
	c.logger.Infof("drm push subscribe res=%+v", res)
	response := model.ConvertToSubscriberRegResult(res)
	if !response.Result {
		c.logger.Errorf("failed to register subscriber: %v+", response.Message)
		return
	}

	ln, ok := c.listeners.Load(response.DataId)
	if !ok { // stale subscriber
		return
	}

	ln.MarkRegistered()
}

func (c *Client) OnDRMPush(dataID string, command string) {
	c.logger.Infof("drm push dataID=%s command=%q", dataID, command)

	var val string

	if strings.HasSuffix(command, model.LocalSuffix) {
		val = command[:(len(command) - len(model.LocalSuffix))]
	} else if strings.HasSuffix(command, model.RemoteSuffix) {
		if strings.HasPrefix(command, model.CmdPrefix) {
			commander := model.NewCommander(command[:(len(command) - len(model.RemoteSuffix))])

			localVersion := -1
			localCache, ok := c.cache.Load(dataID)
			if ok {
				localVersion = localCache.DrmValue.Version
			}

			commandVersion, err := strconv.Atoi(commander.GetProps("version", "-1"))
			if err != nil {
				commandVersion = -1
			}

			zone := commander.GetProps("zone", c.config.zone)
			if localVersion >= commandVersion {
				c.logger.Infof("skip setrequest dataID=%s because local version(%d) is larger than command version(%d)",
					dataID, localVersion, commandVersion)
				return
			}

			remoteValue, remoteVersion, err := c.transport.Fetch(dataID, zone, commandVersion)
			if err != nil {
				c.logger.Errorf("failed to fetch remote version: %v", err)
				return
			}

			if remoteVersion > 0 {
				c.updateLocalValue(dataID, remoteValue, remoteVersion)
				return
			}
		}
	}

	c.tryNotifyListener(dataID, val)
}

func (c *Client) fetchAndUpdateLocalValue(dataID string, version int) (string, int, error) {
	val, version, err := c.transport.Fetch(dataID, c.config.zone, version)
	if err != nil {
		return "", -1, err
	}
	if version > 0 {
		c.updateLocalValue(dataID, val, version)
	} else if version == 0 { // when version == 0, overwrite the local verion but do not notify
		c.cache.Store(dataID, model.LocalValueCache{
			DataId: dataID,
			DrmValue: model.DrmValue{
				Version: version,
			},
			DrmAck: model.DrmAck{
				ActTime: time.Now(),
				Version: version,
			},
		})
	}

	return val, version, nil
}

func (c *Client) updateLocalValue(dataID string, value string, version int) {
	if len(value) > 1024 {
		c.logger.Infof("update local value dataID=%s value=%s...more(%d) version=%d",
			dataID, value[:1024], len(value), version)
	} else {
		c.logger.Infof("update local value dataID=%s value=%s version=%d",
			dataID, value, version)
	}

	cache := model.LocalValueCache{
		DataId: dataID,
		DrmValue: model.DrmValue{
			Value:   value,
			Version: version,
		},
	}

	cache.DrmAck = model.DrmAck{
		ActTime: time.Now(),
		Version: version,
	}

	c.cache.Store(dataID, cache)
	c.tryNotifyListener(dataID, value)
}

func (c *Client) tryNotifyListener(dataID string, value string) {
	ln, ok := c.listeners.Load(dataID)
	if ok {
		ln.OnDRMPush(dataID, value)
	}
}

func (c *Client) LoadLocalValues() []model.LocalValueCache {
	values := make([]model.LocalValueCache, 0, 16)
	c.cache.Range(func(key string, value model.LocalValueCache) bool {
		values = append(values, model.LocalValueCache{
			DataId:   value.DataId,
			DrmAck:   value.DrmAck,
			DrmValue: value.DrmValue,
		})
		return true
	})
	return values
}
