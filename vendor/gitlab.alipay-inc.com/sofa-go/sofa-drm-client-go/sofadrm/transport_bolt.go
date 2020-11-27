package sofadrm

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"gitlab.alipay-inc.com/sofa-go/sofa-drm-client-go/sofadrm/model"
	"go.uber.org/multierr"

	"github.com/gogo/protobuf/proto"
	"gitlab.alipay-inc.com/sofa-go/sofa-bolt-go/sofabolt"
	"gitlab.alipay-inc.com/sofa-go/sofa-logger-go/sofalogger"
)

const (
	defaultRedialNotifyTimeout = 1 * time.Second
	defaultHeartbeatInterval   = 15 * time.Second
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// BOLTTransportOptionSetter configures a BOLTTransport.
type BOLTTransportOptionSetter interface {
	set(*BOLTTransport)
}

type BOLTTransportOptionSetterFunc func(*BOLTTransport)

func (f BOLTTransportOptionSetterFunc) set(c *BOLTTransport) {
	f(c)
}

func WithBOLTTransportLogger(logger sofalogger.Logger) BOLTTransportOptionSetterFunc {
	return BOLTTransportOptionSetterFunc(func(b *BOLTTransport) {
		b.logger = logger
	})
}

func WithBOLTTransportConfig(config *Config) BOLTTransportOptionSetterFunc {
	return BOLTTransportOptionSetterFunc(func(b *BOLTTransport) {
		b.config = config
	})
}

func WithBOLTTransportDialer(dialer sofabolt.Dialer) BOLTTransportOptionSetter {
	return BOLTTransportOptionSetterFunc(func(b *BOLTTransport) {
		b.dialer = dialer
	})
}

func WithBOLTTransportLocator(locator Locator) BOLTTransportOptionSetter {
	return BOLTTransportOptionSetterFunc(func(b *BOLTTransport) {
		b.locator = locator
	})
}

type BOLTTransport struct {
	handlerLock       sync.RWMutex
	handler           sofabolt.Handler
	client            *sofabolt.Client
	config            *Config
	redialCh          chan net.Conn
	logger            sofalogger.Logger
	locator           Locator
	dialer            sofabolt.Dialer
	heartbeatFailures uint64
}

func NewBOLTTransport(options ...BOLTTransportOptionSetter) (*BOLTTransport, error) {
	nt := &BOLTTransport{
		redialCh: make(chan net.Conn),
	}

	for i := range options {
		options[i].set(nt)
	}

	if err := nt.polyfill(); err != nil {
		return nil, err
	}

	return nt, nil
}

func (nt *BOLTTransport) polyfill() error {
	if nt.logger == nil {
		nt.logger = sofalogger.StdoutLogger
	}

	if nt.config == nil {
		return errors.New("bolttransport: config cannot be nil")
	}

	if nt.locator == nil {
		return errors.New("bolttransport: locator cannot be nil")
	}

	if nt.dialer == nil {
		nt.dialer = sofabolt.DialerFunc(nt.dial)
	}

	conn, _ := nt.dialer.Dial()

	client, err := sofabolt.NewClient(
		sofabolt.WithClientConn(conn),
		sofabolt.WithClientRedial(nt.dialer), // always redial
		sofabolt.WithClientTimeout(
			0,                               // no read timeout
			nt.config.boltTransport.timeout, // write timeout
			0,                               // no idle timeout
			0,                               // no flush timeout
		),
		sofabolt.WithClientHandler(nt),
		sofabolt.WithClientMaxPendingCommands( // set max pending commands
			nt.config.GetBOLTTransportConfig().GetMaxPendingCommands()),
		sofabolt.WithClientHeartbeat(
			defaultHeartbeatInterval,
			nt.config.boltTransport.timeout,
			0,
			nt.onHeartbeat,
		),
	)
	if err != nil {
		return err
	}

	nt.client = client

	return nil
}

func (nt *BOLTTransport) onHeartbeat(success bool) {
	if success {
		atomic.StoreUint64(&nt.heartbeatFailures, 0)
		return
	}

	hb := atomic.AddUint64(&nt.heartbeatFailures, 1)
	if int(hb) >= nt.config.boltTransport.maxHeartbeatAttempts { // close the connection and wait redial
		atomic.StoreUint64(&nt.heartbeatFailures, 0)
		nt.logger.Errorf("heartbeat ping failure overflow and try close connection to redial: %d >= %d", hb,
			nt.config.boltTransport.maxHeartbeatAttempts)
		_ = nt.Close()

	} else {
		nt.logger.Errorf("heartbeat ping failure: %d", hb)
	}
}

func (nt *BOLTTransport) Close() error {
	return nt.client.GetConn().Close()
}

func (nt *BOLTTransport) Fetch(dataID string, zone string, localVersion int) (value string, version int, err error) {
	srv, ok := nt.locator.GetRandomServer()
	if !ok {
		return "", 0, nt.locator.RefreshServers()
	}

	var (
		res        *http.Response
		data       []byte
		statuscode int
		url        = fmt.Sprintf("http://%s/queryDrmData.htm?dataId=%s&zone=%s&version=%d&profile=%s&instanceId=%s",
			srv.Ip, dataID, zone, localVersion, nt.config.profile, nt.config.instanceID)
	)

	started := time.Now()
	defer func() {
		if len(data) > 64 {
			nt.logger.Infof("fetch %s %s %d <%s...> [%d] %d <%s> %s", "GET", url,
				statuscode, string(data[:64]), version, len(data), errstr(err), time.Since(started).String())
		} else {
			nt.logger.Infof("fetch %s %s %d <%s> [%d] %d <%s> %s", "GET", url,
				statuscode, string(data), version, len(data), errstr(err), time.Since(started).String())
		}
	}()

	res, err = nt.getHTTPClient().Get(url)
	if err != nil {
		return "", 0, multierr.Append(err, nt.locator.RefreshServers())
	}
	// nolint
	defer res.Body.Close()
	statuscode = res.StatusCode

	headers := res.Header["Drm_version"]
	if len(headers) > 0 {
		version, err = strconv.Atoi(headers[0])
		if err != nil {
			version = -1
		}
	}

	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return "", 0, multierr.Append(err, nt.locator.RefreshServers())
	}

	value = string(data)
	if strings.Contains(value, "@DRM") {
		value = strings.Split(value, "@DRM")[0]
	}

	return value, version, nil
}

func (nt *BOLTTransport) OnRedial(fn func(net.Conn)) {
	for conn := range nt.redialCh {
		go fn(conn)
	}
}

func (nt *BOLTTransport) dial() (net.Conn, error) {
	srv, ok := nt.locator.GetRandomServer()
	if !ok {
		nt.logger.Errorf("redial cannot get server from locator")
		return nil, errors.New("bolttransport: no available server")
	}

	address := srv.GetIPPort()
	conn, err := net.DialTimeout("tcp4", address, nt.config.GetBOLTTransportConfig().GetTimeout())
	if err != nil {
		nt.logger.Errorf("redial dial %s %s", address, err.Error())
	} else {
		nt.logger.Infof("redial choose %s", address)

		// notify we're redialing until timeout
		select {
		case nt.redialCh <- conn:
		case <-time.After(defaultRedialNotifyTimeout):
		}
	}

	return conn, err
}

func (nt *BOLTTransport) ServeSofaBOLT(rw sofabolt.ResponseWriter, req *sofabolt.Request) {
	handler := nt.getHandler()
	if handler == nil {
		nt.logger.Debugf("skip the request")
		return
	}

	handler.ServeSofaBOLT(rw, req)
}

func (nt *BOLTTransport) Send(class string, req, res proto.Message) (err error) {
	started := time.Now()
	defer func() {
		nt.logger.Infof("send class=%s req=%s res=%s err=%s elapsed=%s", class, protostr(req),
			protostr(res), errstr(err), time.Since(started).String())
	}()

	breq := sofabolt.AcquireRequest()
	bres := sofabolt.AcquireResponse()
	defer func() {
		sofabolt.ReleaseRequest(breq)
		sofabolt.ReleaseResponse(bres)
	}()
	breq.SetCodec(sofabolt.CodecProtobuf)
	if res == nil { // oneway
		breq.SetType(sofabolt.TypeBOLTRequestOneWay)
	}

	var data []byte
	data, err = proto.Marshal(req)
	if err != nil {
		return err
	}
	breq.SetClassString(class).SetContent(data)

	if err = nt.client.DoTimeout(breq, bres, nt.config.GetBOLTTransportConfig().GetTimeout()); err != nil {
		return err
	}

	if res != nil {
		err = proto.Unmarshal(bres.GetContent(), res)
	}

	return err
}

func (nt *BOLTTransport) accesslogHandler(fn sofabolt.HandlerFunc) sofabolt.HandlerFunc {
	return sofabolt.HandlerFunc(func(rw sofabolt.ResponseWriter, req *sofabolt.Request) {
		start := time.Now()
		fn(rw, req)
		nt.logger.Infof("accesslog %s %s %d %d %s",
			string(req.GetClass()),
			string(rw.GetResponse().GetClass()),
			len(req.GetContent()),
			len(rw.GetResponse().GetContent()),
			time.Since(start).String(),
		)
	})
}

func (nt *BOLTTransport) recvHandler(fn func(err error, class string, req, res proto.Message)) sofabolt.HandlerFunc {
	return sofabolt.HandlerFunc(func(rw sofabolt.ResponseWriter, breq *sofabolt.Request) {
		defer func() {
			if r := recover(); r != nil {
				nt.logger.Errorf("recover from panic: %+v", r)
			}
		}()

		if breq.GetCMDCode() == sofabolt.CMDCodeBOLTHeartbeat {
			rw.GetResponse().SetCMDCode(sofabolt.CMDCodeBOLTHeartbeat)
			_, err := rw.Write()
			if err != nil {
				nt.logger.Errorf("heartbeat pong failure: %s", err.Error())
			} else {
				nt.logger.Infof("heartbeat pong success")
			}
			return
		}

		var (
			req   proto.Message
			res   proto.Message
			err   error
			class = string(breq.GetClass())
		)

		switch class {
		case SubscriberRegResultClass:
			req = &model.SubscriberRegResultPb{}
			err = proto.Unmarshal(breq.GetContent(), req)
			if err != nil {
				err = fmt.Errorf("failed to unmarshal subscriber result: %s", err.Error())
				break
			}

		case AttributeSetRequestClass:
			req = &model.AttributeSetRequestPb{}
			err = proto.Unmarshal(breq.GetContent(), req)
			if err != nil {
				err = fmt.Errorf("failed to unmarshal setrequest result: %s", err.Error())
				break
			}

		case AttributeGetRequestClass:
			req = &model.AttributeGetRequestPb{}
			err = proto.Unmarshal(breq.GetContent(), req)
			if err != nil {
				err = fmt.Errorf("failed to unmarshal getrequest result: %s", err.Error())
				break
			}
			res = new(model.AttributeGetResponse)

		default:
			err = fmt.Errorf("unknown request: %s", class)
			break
		}

		fn(err, class, req, res)

		if res == nil { // one way
			return
		}

		data, err := proto.Marshal(res)
		if err != nil {
			nt.logger.Errorf("failed to marshal response: %s", err.Error())
		}

		rw.GetResponse().SetCodec(sofabolt.CodecProtobuf).SetContent(data)
		_, err = rw.Write()
		if err != nil {
			nt.logger.Errorf("failed to send response(%s): %s", class, err.Error())
		}
	})
}

func (nt *BOLTTransport) OnRecv(fn func(err error, class string, req, res proto.Message)) error {
	nt.setHandler(nt.accesslogHandler(nt.recvHandler(fn)))
	// block until see the error
	return <-nt.client.GetReadError()
}

func (hl *BOLTTransport) getHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   hl.config.boltTransport.timeout,
				KeepAlive: 0,
				DualStack: false,
			}).DialContext,
			DisableKeepAlives:   true, // short-lived connection
			MaxIdleConns:        0,
			TLSHandshakeTimeout: hl.config.boltTransport.timeout,
		},
		Timeout: hl.config.boltTransport.timeout,
	}
}

func (hl *BOLTTransport) getHandler() sofabolt.Handler {
	hl.handlerLock.RLock()
	handler := hl.handler
	hl.handlerLock.RUnlock()
	return handler
}

func (hl *BOLTTransport) setHandler(fn sofabolt.Handler) {
	hl.handlerLock.Lock()
	hl.handler = fn
	hl.handlerLock.Unlock()
}
