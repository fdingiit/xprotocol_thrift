package mist

import (
	"fmt"
	"mosn.io/pkg/utils"
	"os"
	"sync"
	"time"

	"gitlab.alipay-inc.com/infrasec/api/mist/types"
)

type Issuer struct {
	lock   sync.RWMutex
	config *types.Issue

	_svid_ string
	stop   chan struct{}
}

func NewIssuer(config *types.Issue) (*Issuer, error) {
	if err := checkIssueConfig(config); err != nil {
		return nil, err
	}
	issuer := &Issuer{
		config: config,
		stop:   make(chan struct{}),
	}
	return issuer, nil
}

func (issuer *Issuer) UpdateConfig(config *types.Issue) error {
	if err := checkIssueConfig(config); err != nil {
		return err
	}
	issuer.lock.Lock()
	issuer.config = config
	issuer.lock.Unlock()
	return nil
}

func (issuer *Issuer) GenSvid(exp int64) (string, error) {
	return issuer.issueSvid(exp)
}

func (issuer *Issuer) GetSvidFromCache() (string, error) {
	issuer.lock.RLock()
	defer issuer.lock.RUnlock()
	if len(issuer._svid_) <= 0 {
		return "", fmt.Errorf("svid is empty")
	}
	return issuer._svid_, nil
}

func (issuer *Issuer) Start() error {
	issuer.run()
	return nil
}

func (issuer *Issuer) Stop() error {
	issuer.stop <- struct{}{}
	return nil
}

func (issuer *Issuer) run() {
	utils.GoWithRecover(func() {
		for {
			svid, err := issuer.issueSvid(0)
			if err == nil && svid != "" {
				issuer.lock.Lock()
				issuer._svid_ = svid
				issuer.lock.Unlock()
			}

			select {
			case <-issuer.stop:
				return
			case <-time.After(time.Duration(issuer.getCacheInterval()/2+1) * time.Second):
			}
		}
	}, nil)
}

func (issuer *Issuer) issueSvid(exp int64) (svid string, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "issuer error: %v", r)
			err = fmt.Errorf("issuer error: %v", r)
		}
	}()
	return IssueJWTSVID(issuer.getUrl(), exp)
}

func (issuer *Issuer) getUrl() string {
	issuer.lock.RLock()
	defer issuer.lock.RUnlock()
	return issuer.config.Url
}

func (issuer *Issuer) getCacheInterval() int64 {
	issuer.lock.RLock()
	defer issuer.lock.RUnlock()
	return issuer.config.CacheInterval
}

func checkIssueConfig(config *types.Issue) error {
	if config == nil {
		return fmt.Errorf("config is nil")
	}
	if len(config.Url) <= 0 {
		return fmt.Errorf("url is empty")
	}
	if config.CacheInterval <= 0 {
		return fmt.Errorf("CacheInterval error[%d]", config.CacheInterval)
	}
	return nil
}
