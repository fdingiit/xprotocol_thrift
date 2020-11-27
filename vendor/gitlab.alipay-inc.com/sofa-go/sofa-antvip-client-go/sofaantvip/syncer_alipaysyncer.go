package sofaantvip

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
	"time"

	"gitlab.alipay-inc.com/sofa-go/sofa-antvip-client-go/sofaantvip/protobuf"

	"github.com/gogo/protobuf/proto"
	"gitlab.alipay-inc.com/sofa-go/sofa-logger-go/sofalogger"
)

// AlipaySyncerOptionSetter configures a AlipaySyncer.
type AlipaySyncerOptionSetter interface {
	set(*AlipaySyncer)
}

type AlipaySyncerOptionSetterFunc func(*AlipaySyncer)

func (f AlipaySyncerOptionSetterFunc) set(c *AlipaySyncer) {
	f(c)
}

func WithAlipaySyncerConfig(config *Config) AlipaySyncerOptionSetterFunc {
	return AlipaySyncerOptionSetterFunc(func(c *AlipaySyncer) {
		c.config = config
	})
}

func WithAlipaySyncerLocator(locator Locator) AlipaySyncerOptionSetterFunc {
	return AlipaySyncerOptionSetterFunc(func(c *AlipaySyncer) {
		c.locator = locator
	})
}

func WithAlipaySyncerLogger(logger sofalogger.Logger) AlipaySyncerOptionSetterFunc {
	return AlipaySyncerOptionSetterFunc(func(c *AlipaySyncer) {
		c.logger = logger
	})
}

type AlipaySyncer struct {
	sync.RWMutex
	config       *Config
	locator      Locator
	logger       sofalogger.Logger
	zoneInfoList *ZoneInfoList
}

func NewAlipaySyncer(options ...AlipaySyncerOptionSetter) (*AlipaySyncer, error) {
	hs := &AlipaySyncer{}
	for i := range options {
		options[i].set(hs)
	}

	if err := hs.polyfill(); err != nil {
		return nil, err
	}

	return hs, nil
}

func (hs *AlipaySyncer) polyfill() error {
	if hs.config == nil {
		return errors.New("alipaysyncer: config cannot be nil")
	}

	hs.SetZoneInfoList(hs.config.GetAlipaySyncerConfig().GetZoneInfoList())

	if hs.locator == nil {
		return errors.New("alipaysyncer: config cannot be nil")
	}

	if hs.logger == nil {
		hs.logger = sofalogger.StdoutLogger
	}

	return nil
}

func (hs *AlipaySyncer) GetVipDomainList(domains map[string]string, polling bool,
	timeout time.Duration) (vipdomains []VipDomain, err error) {
	started := time.Now()
	defer func() {
		if err != nil {
			hs.config.metrics.alipaysyncer.addFailure()
		} else {
			hs.config.metrics.alipaysyncer.addSuccess()
		}
	}()

	var (
		url        string
		statuscode int
	)

	vipdomains, url, statuscode, err = hs.getVipDomainList(domains, polling, timeout)

	if hs.config.alipaySyncer.accesslog {
		var p func(format string, a ...interface{})
		if err != nil {
			p = hs.logger.Errorf
		} else {
			p = hs.logger.Infof
		}

		p("AlipaySyncer %s(polling=%t) %s %d %s[%d] <%s> %s timeout=%s", "POST", polling, url,
			statuscode, prettyDomains(domains), len(vipdomains),
			errstring(err), time.Since(started).String(), timeout.String())
	}

	if err != nil {
		// If see the error, refresh the server list.
		return nil, merror(err, hs.locator.RefreshServers())
	}

	return vipdomains, nil
}

func (hs *AlipaySyncer) getVipDomainList(
	domains map[string]string,
	polling bool,
	timeout time.Duration) (
	vipdomains []VipDomain,
	url string,
	statuscode int,
	err error) {
	if len(domains) == 0 {
		return nil, "", 0, nil
	}

	var pres *pollingResponse

	statuscode = -1

	server, ok := hs.locator.GetRandomServer()
	if !ok {
		return nil, "", statuscode, errors.New("alipaysyncer: no servers")
	}

	url = hs.config.GetAlipaySyncerConfig().GetURL(server.Host)

	preq := hs.buildPollingRequest(hs.config, domains, polling)
	body, err := proto.Marshal(preq.ToProtobuf())
	if err != nil {
		return nil, url, statuscode, err
	}

	pres, statuscode, err = hs.loadVipDomainsFromURL(url, body,
		timeout)
	if err != nil {
		return nil, url, statuscode, err
	}

	if pres.zoneInfoList != nil {
		hs.SetZoneInfoList(pres.zoneInfoList)
	}

	vipdomains = pres.LoadDomains()

	return vipdomains, url, statuscode, nil
}

func (hs *AlipaySyncer) SetZoneInfoList(zi *ZoneInfoList) {
	if zi != nil {
		hs.Lock()
		hs.zoneInfoList = zi
		hs.Unlock()
	}
}

func (hs *AlipaySyncer) GetZoneInfoListCheckSum() string {
	hs.Lock()
	defer hs.Unlock()
	if hs.zoneInfoList != nil {
		return hs.zoneInfoList.ChecksumAlipay()
	}
	return "N"
}

func (hs *AlipaySyncer) loadVipDomainsFromURL(url string, body []byte,
	timeout time.Duration) (*pollingResponse, int, error) {
	var (
		res        *http.Response
		err        error
		statuscode int
		data       []byte
	)

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, -1, err
	}

	req.Header.Add("Serialize-type", "pb3")
	res, err = hs.getHTTPClient(timeout).Do(req)
	if err != nil {
		return nil, -1, err
	}
	// nolint
	defer res.Body.Close()
	statuscode = res.StatusCode

	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, statuscode, err
	}

	var responseMsg protobuf.PollingResponseMsg
	err = proto.Unmarshal(data, &responseMsg)
	if err != nil {
		return nil, statuscode, err
	}

	var pr pollingResponse
	pr.FromProtobuf(&responseMsg)
	pr.Polyfill()

	return &pr, statuscode, nil
}

func (hl *AlipaySyncer) getHTTPClient(timeout time.Duration) *http.Client {
	if timeout == 0 {
		timeout = hl.config.GetSyncTimeout()
	}

	return &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   timeout,
				KeepAlive: 15 * time.Second,
				DualStack: false,
			}).DialContext,
			DisableKeepAlives:   true, // short-lived connection
			MaxIdleConns:        0,
			TLSHandshakeTimeout: timeout,
		},
		Timeout: timeout,
	}
}

func (hs *AlipaySyncer) buildPollingRequest(config *Config, domains map[string]string, polling bool) pollingRequest {
	pr := newPollingRequest(config)
	pr.ExtensionParams[ExtensionZoneInfoList] = hs.GetZoneInfoListCheckSum()
	pr.VipDomainName2ChecksumMap = domains
	pr.AllowPolling = polling
	return pr
}
