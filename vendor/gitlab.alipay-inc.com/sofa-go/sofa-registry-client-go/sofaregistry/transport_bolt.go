package sofaregistry

import (
	"bytes"
	"errors"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gogo/protobuf/proto"
	"gitlab.alipay-inc.com/sofa-go/sofa-bolt-go/sofabolt"
	"gitlab.alipay-inc.com/sofa-go/sofa-logger-go/sofalogger"
	sofaregistryproto "gitlab.alipay-inc.com/sofa-go/sofa-registry-proto-go/proto"
)

var _ Transport = (*BOLTTransport)(nil)

const (
	defaultRedialTimeout     = 1 * time.Second
	defaultHeartbeatInterval = 15 * time.Second
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

type BOLTTransport struct {
	handlerLock       sync.RWMutex
	handler           sofabolt.Handler
	client            *sofabolt.Client
	config            *Config
	redialCh          chan net.Conn
	logger            sofalogger.Logger
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

	if nt.dialer == nil {
		nt.dialer = sofabolt.DialerFunc(nt.dial)
	}
	redialer := newRedialer(nt.dialer, nt.redialCh)
	nt.redialCh = redialer.redialCh

	// do dial to get the connection whatever if it's failed
	conn, _ := nt.dialer.Dial()

	client, err := sofabolt.NewClient(
		sofabolt.WithClientConn(conn),
		sofabolt.WithClientRedial(redialer), // always redial
		sofabolt.WithClientTimeout(
			0,                                     // no read timeout
			nt.config.boltTransportConfig.timeout, // write timeout
			0,                                     // no idle timeout
			0,                                     // no flush timeout
		),
		sofabolt.WithClientHandler(nt),
		sofabolt.WithClientHeartbeat(defaultHeartbeatInterval,
			nt.config.boltTransportConfig.timeout,
			0,
			nt.onHeartbeat,
		),
		sofabolt.WithClientMaxPendingCommands( // set max pending commands
			nt.config.GetBOLTTransportConfig().GetMaxPendingCommands()),
	)
	if err != nil {
		return err
	}

	nt.client = client

	return nil
}

func (nt *BOLTTransport) onHeartbeat(success bool) {
	if success {
		nt.logger.Infof("heartbeat ping success")
		atomic.StoreUint64(&nt.heartbeatFailures, 0)
		return
	}

	hb := atomic.AddUint64(&nt.heartbeatFailures, 1)
	if int(hb) >= nt.config.boltTransportConfig.maxHeartbeatAttempts { // close the connection and wait redial
		atomic.StoreUint64(&nt.heartbeatFailures, 0)
		nt.logger.Errorf("heartbeat ping failure overflow and try close connection to redial: %d >= %d", hb,
			nt.config.boltTransportConfig.maxHeartbeatAttempts)
		_ = nt.Close()

	} else {
		nt.logger.Errorf("heartbeat ping failure: %d", hb)
	}
}

func (nt *BOLTTransport) Close() error {
	return nt.client.GetConn().Close()
}

func (nt *BOLTTransport) OnRedial(fn func(net.Conn)) {
	for conn := range nt.redialCh {
		go fn(conn)
	}
}

func (nt *BOLTTransport) dial() (net.Conn, error) {
	srv, err := nt.GetRandomServer()
	if err != nil {
		nt.logger.Errorf("redial get %s %s", nt.config.GetBOLTTransportLocatorURL(), errstring(err))
		return nil, err
	}

	conn, err := net.DialTimeout("tcp4", srv.GetAddress(), nt.config.GetBOLTTransportConfig().GetTimeout())
	if err != nil {
		nt.logger.Errorf("redial choose %s %s", srv.GetAddress(), errstring(err))
	} else {
		nt.logger.Infof("redial choose %s", srv.GetAddress())
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

func (nt *BOLTTransport) Send(class string, req proto.Message, res proto.Message) error {
	breq := sofabolt.AcquireRequest()
	bres := sofabolt.AcquireResponse()
	defer func() {
		sofabolt.ReleaseRequest(breq)
		sofabolt.ReleaseResponse(bres)
	}()
	breq.SetCodec(sofabolt.CodecProtobuf)

	data, err := proto.Marshal(req)
	if err != nil {
		return err
	}
	breq.SetClassString(class).SetContent(data)

	if err = nt.client.DoTimeout(breq, bres, nt.config.GetBOLTTransportConfig().GetTimeout()); err != nil {
		return err
	}

	return proto.Unmarshal(bres.GetContent(), res)
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

func (nt *BOLTTransport) subscribeHandler(fn func(err error, req proto.Message)) sofabolt.HandlerFunc {
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
			message string
			success = true
			preq    = new(sofaregistryproto.ReceivedDataPb)
		)

		err := proto.Unmarshal(breq.GetContent(), preq)
		if err != nil {
			success = false
			message = err.Error()
		}

		d, err := proto.Marshal(&sofaregistryproto.ResultPb{
			Success: success,
			Message: message,
		})
		if err == nil {
			if rw.GetResponse().GetType() != sofabolt.TypeBOLTRequestOneWay {
				rw.GetResponse().SetCodec(sofabolt.CodecProtobuf)
				rw.GetResponse().SetClassString(RESPONSEPbClass).SetContent(d)
				_, err = rw.Write()
			}
		}

		fn(err, preq)
	})
}

func (nt *BOLTTransport) OnRecv(fn func(err error, req proto.Message)) error {
	nt.setHandler(nt.accesslogHandler(nt.subscribeHandler(fn)))
	// block until see the error
	return <-nt.client.GetReadError()
}

func (hl *BOLTTransport) GetServers() ([]*Server, error) {
	var (
		res        *http.Response
		err        error
		data       []byte
		statuscode int
		url        string
	)

	started := time.Now()
	defer func() {
		logf := hl.logger.Infof
		if err != nil {
			logf = hl.logger.Errorf
		}
		logf("redial get %s %d %s <%s> %s", url,
			statuscode, string(data), errstring(err), time.Since(started).String())
	}()

	url = hl.config.GetBOLTTransportLocatorURL()
	res, err = hl.getHTTPClient().Get(url)
	if err != nil {
		return nil, err
	}
	// nolint
	defer res.Body.Close()
	statuscode = res.StatusCode

	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	servers := make([]*Server, 0, len(data))
	split := bytes.Split(bytes.TrimSpace(data), []byte(";"))
	for i := range split {
		address := string(split[i])
		// should be host:port
		_, _, err = net.SplitHostPort(address)
		if err != nil {
			continue
		}
		servers = append(servers, &Server{
			address: address,
		})
	}

	return servers, nil
}

func (hl *BOLTTransport) GetRandomServer() (*Server, error) {
	servers, err := hl.GetServers()
	if err != nil {
		return nil, err
	}

	if len(servers) == 0 {
		return nil, errors.New("bolttransport: no available server")
	}

	return servers[rand.Intn(len(servers))], nil
}

func (hl *BOLTTransport) getHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   hl.config.boltTransportConfig.timeout,
				KeepAlive: 0,
				DualStack: false,
			}).DialContext,
			DisableKeepAlives:   true, // short-lived connection
			MaxIdleConns:        0,
			TLSHandshakeTimeout: hl.config.boltTransportConfig.timeout,
		},
		Timeout: hl.config.boltTransportConfig.timeout,
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

type redialer struct {
	dialer   sofabolt.Dialer
	redialCh chan net.Conn
}

func newRedialer(dialer sofabolt.Dialer, redialCh chan net.Conn) *redialer {
	return &redialer{
		dialer:   dialer,
		redialCh: redialCh,
	}
}

func (d *redialer) Dial() (net.Conn, error) {
	conn, err := d.dialer.Dial()
	if err == nil {
		select {
		case d.redialCh <- conn:
		case <-time.After(defaultRedialTimeout):
		}
	}
	return conn, err
}
