package sofaantvip

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"gitlab.alipay-inc.com/sofa-go/sofa-logger-go/sofalogger"
)

// CloudSyncerOptionSetter configures a CloudSyncer.
type CloudSyncerOptionSetter interface {
	set(*CloudSyncer)
}

type CloudSyncerOptionSetterFunc func(*CloudSyncer)

func (f CloudSyncerOptionSetterFunc) set(c *CloudSyncer) {
	f(c)
}

func WithCloudSyncerConfig(config *Config) CloudSyncerOptionSetterFunc {
	return CloudSyncerOptionSetterFunc(func(c *CloudSyncer) {
		c.config = config
	})
}

func WithCloudSyncerLogger(logger sofalogger.Logger) CloudSyncerOptionSetterFunc {
	return CloudSyncerOptionSetterFunc(func(c *CloudSyncer) {
		c.logger = logger
	})
}

type CloudSyncer struct {
	config *Config
	logger sofalogger.Logger
}

func NewCloudSyncer(options ...CloudSyncerOptionSetter) (*CloudSyncer, error) {
	hs := &CloudSyncer{}
	for i := range options {
		options[i].set(hs)
	}

	if err := hs.polyfill(); err != nil {
		return nil, err
	}

	return hs, nil
}

func (hs *CloudSyncer) polyfill() error {
	if hs.config == nil {
		return errors.New("CloudSyncer: config cannot be nil")
	}

	if hs.logger == nil {
		hs.logger = sofalogger.StdoutLogger
	}

	return nil
}

func (hs *CloudSyncer) GetVipDomainList(domains map[string]string,
	polling bool, timeout time.Duration) ([]VipDomain, error) {
	started := time.Now()
	v, url, statuscode, err := hs.getVipDomainList(domains, polling, timeout)

	if hs.config.cloudSyncer.accesslog {
		var p func(format string, a ...interface{})
		if err != nil {
			p = hs.logger.Errorf
		} else {
			p = hs.logger.Infof
		}

		p("AlipaySyncer %s(polling=%t) %s %d %s[%d] <%s> %s timeout=%s", "POST", polling, url,
			statuscode, prettyDomains(domains), len(v),
			errstring(err), time.Since(started).String(), timeout.String())
	}

	if err != nil {
		hs.config.metrics.CloudSyncer.addFailure()
	} else {
		hs.config.metrics.CloudSyncer.addSuccess()
	}

	return v, err
}

func (hs *CloudSyncer) getVipDomainList(
	domains map[string]string,
	polling bool,
	timeout time.Duration) (
	vipdomains []VipDomain,
	url string,
	statuscode int,
	err error) {
	url = hs.config.GetCloudSyncerConfig().GetURL()
	if len(domains) == 0 {
		return nil, url, 0, nil
	}

	var pres *pollingResponse

	statuscode = -1

	body, err := json.Marshal(hs.buildPollingRequest(hs.config, domains, polling))
	if err != nil {
		return nil, url, statuscode, err
	}

	pres, statuscode, err = hs.loadVipDomainsFromURL(url, body, timeout)
	if err != nil {
		return nil, url, statuscode, err
	}

	return pres.LoadDomains(), url, statuscode, nil
}

func (hs *CloudSyncer) loadVipDomainsFromURL(url string, body []byte,
	timeout time.Duration) (*pollingResponse, int, error) {
	var (
		res        *http.Response
		err        error
		data       []byte
		statuscode int
	)

	res, err = hs.getHTTPClient(timeout).Post(url, "application/json", bytes.NewReader(body))
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

	var pr pollingResponse
	if err := json.Unmarshal(data, &pr); err != nil {
		return nil, statuscode, err
	}
	pr.Polyfill()

	return &pr, statuscode, nil
}

func (hl *CloudSyncer) getHTTPClient(timeout time.Duration) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   timeout,
				KeepAlive: 0,
				DualStack: false,
			}).DialContext,
			DisableKeepAlives:   true, // short-lived connection
			MaxIdleConns:        0,
			TLSHandshakeTimeout: timeout,
		},
		Timeout: timeout,
	}
}

func (hs *CloudSyncer) buildPollingRequest(config *Config, domains map[string]string, polling bool) pollingRequest {
	preq := newPollingRequest(config)
	preq.AllowPolling = polling
	extensionMap := preq.ExtensionParams
	extensionMap["ChecksumNewVersionSign"] = "true"
	extensionMap["EXTENSION_ZONE_INFO_LIST"] = "N"
	preq.DataCenter = hs.config.datacenter
	preq.VipDomainName2ChecksumMap = domains
	return preq
}
