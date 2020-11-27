package sofaantvip

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/detailyang/keymutex-go"
	"gitlab.alipay-inc.com/sofa-go/sofa-logger-go/sofalogger"
)

var ErrDomainNotFound = errors.New("sofaantvip: not found domain")

const (
	ProbeDomain  = "probe.antvip"
	ProbeTimeout = 5 * time.Second
)

//go:generate syncmap -pkg sofaantvip -o domainlistenermap_generated.go -name domainListenerMap map[string]*domainListener

type AntvipClient struct {
	context       context.Context
	config        *Config
	locks         *keymutex.KeyMutex
	domains       domainListenerMap
	stickydomains sync.Map // map[string]struct{}
	logger        sofalogger.Logger
	syncer        Syncer
}

func NewAntvipClient(options ...AntvipClientOptionSetter) (*AntvipClient, error) {
	ac := &AntvipClient{
		locks: keymutex.New(256),
	}

	for _, option := range options {
		option.set(ac)
	}

	if err := ac.polyfill(); err != nil {
		return nil, err
	}

	if ac.config.syncInterval > 0 {
		go ac.DoIntervalSync()
	}

	return ac, nil
}

// NewDummyAntvipClient build a dummy client
//
// FYI: hack for mosn silly usage
func NewDummyAntvipClient() *AntvipClient {
	return &AntvipClient{
		locks:   keymutex.New(256),
		config:  new(Config),
		logger:  sofalogger.StdoutLogger,
		context: context.TODO(),
		syncer:  NewDirectSyncer(make(map[string]VipDomain)),
	}
}

func (ac *AntvipClient) DoSync(timeout time.Duration) error {
	domains := make(map[string]string, 32)

	ac.domains.Range(func(k string, ln *domainListener) bool {
		d := ln.getDomain()
		if d == nil || ln.getName() == "" {
			domains[k] = "N"
		} else {
			domains[k] = ac.getDomainChecksum(d)
		}
		return true
	})

	vipdomains, err := ac.getVipDomain(domains, true, timeout)
	ac.BroadcastDomainWithError(err, vipdomains)

	return err
}

// DoIntervalSync does the work of sync in the given interval.
func (ac *AntvipClient) DoIntervalSync() {
	timer := time.NewTimer(ac.config.syncInterval)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			if err := ac.DoSync(ac.config.syncTimeout); err != nil {
				ac.logger.Errorf("failed to sync antvip: %s", err.Error())
			}
		case <-ac.context.Done():
			return
		}
		timer.Reset(ac.config.syncInterval)
	}
}

func (ac *AntvipClient) polyfill() error {
	if ac.logger == nil {
		ac.logger = sofalogger.StdoutLogger
	}

	if ac.context == nil {
		ac.context = context.TODO()
	}

	if ac.syncer == nil {
		return errors.New("sofaantvip: syncer cannot be nil")
	}

	if ac.config == nil {
		ac.config = NewConfig()
	}

	return nil
}

func (ac *AntvipClient) BroadcastDomain(domains []VipDomain) {
	ac.BroadcastDomainWithError(nil, domains)
}

func (ac *AntvipClient) BroadcastDomainWithError(err error, domains []VipDomain) {
	for i := range domains {
		domain := domains[i]
		dl, ok := ac.domains.Load(domain.Name)
		if !ok {
			continue
		}
		listeners := dl.getListeners()
		for j := range listeners {
			listeners[j].OnDomainChanged(err, domain.Name, &domain)
		}
	}
}

func (ac *AntvipClient) AddListener(domain string, listener Listener) {
	ac.locks.LockKey(domain)
	defer ac.locks.UnlockKey(domain)

	dl, ok := ac.domains.Load(domain)
	if !ok {
		ac.domains.Store(domain, newDomainListener(domain, listener))
	} else {
		dl.addListener(listener)
	}
}

func (ac *AntvipClient) FetchListeners(domain string) ([]Listener, bool) {
	dl, ok := ac.domains.Load(domain)
	if !ok {
		return nil, false
	}
	return dl.getListeners(), true
}

func (ac *AntvipClient) GetConfig() *Config {
	return ac.config
}

func (ac *AntvipClient) GetLocalZone() string {
	return ac.config.GetZone()
}

func (ac *AntvipClient) GetVipDomain(domainName string, timeout time.Duration) (*VipDomain, error) {
	ac.locks.LockKey(domainName)
	defer ac.locks.UnlockKey(domainName)

	ld, ok := ac.domains.Load(domainName)
	if ok { // if hit the cache
		vd := ld.getDomain()
		if vd != nil {
			if vd.IsDeleted {
				return nil, ErrDomainNotFound
			}
			return vd, nil
		}
	}

	// Store the domain to next poll if need
	ac.domains.LoadOrStore(domainName, newDomainListener(domainName))

	vds, err := ac.getVipDomain(map[string]string{domainName: "N"}, false, timeout)
	if err != nil {
		return nil, err
	}
	if len(vds) < 1 {
		return nil, ErrDomainNotFound
	}

	return &vds[0], nil
}

func (ac *AntvipClient) ReserveDomain(domain string) {
	ac.stickydomains.Store(domain, struct{}{})
}

func (ac *AntvipClient) ResetDomains() {
	ac.domains.Range(func(key string, value *domainListener) bool {
		_, ok := ac.stickydomains.Load(key)
		if !ok { // Cleanup
			ac.domains.Delete(key)
		}
		return true
	})
}

func (ac *AntvipClient) getVipDomain(domains map[string]string, polling bool,
	timeout time.Duration) ([]VipDomain, error) {
	vds, err := ac.syncer.GetVipDomainList(domains, polling, timeout)
	if err != nil {
		return nil, err
	}

	for i := range vds {
		domain := vds[i]
		dl, ok := ac.domains.Load(domain.Name)
		if !ok { // stale domain
			continue
		}
		dl.setDomain(&domain)
	}
	return vds, nil
}

func (ac *AntvipClient) getDomainChecksum(domain *VipDomain) string {
	if ac.config.domainChecksumMode == DomainAlipayChecksumMode {
		return domain.ChecksumAlipay()
	}
	return domain.ChecksumCloud()
}

func (ac *AntvipClient) GetLocalVipDomains() []*VipDomain {
	var domains []*VipDomain
	ac.domains.Range(func(key string, dl *domainListener) bool {
		if d := dl.getDomain(); d != nil {
			domains = append(domains, dl.getDomain())
		}
		return true
	})
	return domains
}

func (ac *AntvipClient) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type status struct {
		Config  *Config `json:"config"`
		Metrics struct {
			Success    int64 `json:"success"`
			Failure    int64 `json:"failure"`
			LastSyntAt int64 `json:"last_sync_at"`
		}
		Domains []VipDomain `json:"domains"`
		Health  string      `json:"health"`
	}

	s := &status{}
	s.Metrics.Success = ac.config.metrics.GetAlipaySyncerMetric().LoadSuccess()
	s.Metrics.Failure = ac.config.metrics.GetAlipaySyncerMetric().LoadFailure()
	s.Metrics.LastSyntAt = ac.config.metrics.GetAlipaySyncerMetric().LoadLastSyncAt()

	if err := ac.IsHealth(); err != nil {
		s.Health = err.Error()
	} else {
		s.Health = "ok"
	}

	if err := json.NewEncoder(w).Encode(&s); err != nil {
		ac.logger.Errorf("failed to write response")
	}
}

func (ac *AntvipClient) IsHealth() error {
	_, err := ac.GetVipDomain(ProbeDomain, ProbeTimeout)
	if err == ErrDomainNotFound { // allow domain not found
		return nil
	}

	return err
}
