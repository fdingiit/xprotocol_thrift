package mist

import (
	"encoding/base64"
	"fmt"
	"github.com/golang/groupcache/lru"
	"github.com/json-iterator/go"
	"gitlab.alipay-inc.com/infrasec/api/mist/types"
	"math"
	"mosn.io/pkg/utils"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	preProcessSvidLen      = 128
	maxSvidCacheLen        = 4096
	defaultExpiration      = 15 * 60
	cleanupInterval        = 5 * 60
	svidExpireTolerantTime = 15 * 60
)

const DefaultSpiffeIDPrefix = "spiffe://alipay.com/"

type SVID struct {
	Sub string
	Exp int64

	Valid bool
}

type Verifier struct {
	lock           sync.RWMutex
	config         *types.Verify
	svidCache      *lru.Cache
	preProcessSvid chan string
	stop           chan struct{}
}

func NewVerifier(config *types.Verify) (*Verifier, error) {
	if err := checkVerifyConfig(config); err != nil {
		return nil, err
	}
	verifier := &Verifier{
		config:         config,
		stop:           make(chan struct{}),
		svidCache:      lru.New(maxSvidCacheLen),
		preProcessSvid: make(chan string, preProcessSvidLen),
	}
	return verifier, nil
}

func (verifier *Verifier) UpdateConfig(config *types.Verify) error {
	if err := checkVerifyConfig(config); err != nil {
		return err
	}
	verifier.lock.Lock()
	verifier.config = config
	verifier.svidCache = lru.New(int(config.SvidCacheCap))
	verifier.lock.Unlock()
	return nil
}

func (verifier *Verifier) Verify(svid string) (obj SVID, err error) {
	return verifier.VerifyWithCondition(svid, verifier.getSync())
}

func (verifier *Verifier) VerifyWithCondition(svid string, sync bool) (obj SVID, err error) {
	svid_struct, err := getSVIDStruct(svid)
	if err != nil {
		return SVID{}, err
	}

	v, ok := verifier.svidCache.Get(svid)
	if ok {
		obj = v.(SVID)
		if !obj.Valid {
			err = fmt.Errorf("svid is invalid")
		}
		return
	}

	if sync {
		obj, err := verifier.verifySvid(svid)
		verifier.svidCache.Add(svid, obj)
		return obj, err
	}

	select {
	case verifier.preProcessSvid <- svid:
	default:
	}

	// first time verify the svid
	svid_struct.Valid = true
	return svid_struct, nil
}

func (verifier *Verifier) Start() error {
	verifier.run()
	return nil
}

func (verifier *Verifier) Stop() error {
	verifier.stop <- struct{}{}
	return nil
}

func (verifier *Verifier) run() {
	utils.GoWithRecover(func() {
		for {
			select {
			case <-verifier.stop:
				return
			case svid := <-verifier.preProcessSvid:
				obj, _ := verifier.verifySvid(svid)
				verifier.svidCache.Add(svid, obj)
			}
		}
	}, nil)
}

func (verifier *Verifier) verifySvid(svid string) (obj SVID, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "verify error: %v", r)
			err = fmt.Errorf("issuer error: %v", r)
		}
	}()

	svid_struct, err := getSVIDStruct(svid)
	if err != nil {
		return
	}
	ok, err := VerifyJWTSVID(verifier.getUrl(), svid)
	svid_struct.Valid = ok
	return svid_struct, nil
}

func (verifier *Verifier) getUrl() string {
	verifier.lock.RLock()
	defer verifier.lock.RUnlock()
	return verifier.config.GetUrl()
}

func (verifier *Verifier) getSync() bool {
	verifier.lock.RLock()
	defer verifier.lock.RUnlock()
	return verifier.config.GetSync()
}

func checkVerifyConfig(config *types.Verify) error {
	if config == nil {
		return fmt.Errorf("config is nil")
	}
	if len(config.GetUrl()) <= 0 {
		return fmt.Errorf("url is empty")
	}
	if config.GetSvidCacheCap() <= 0 {
		return fmt.Errorf("SvidCacheCap error[%d]", config.GetSvidCacheCap())
	}
	return nil
}

func getSVIDStruct(svid string) (SVID, error) {
	sp := strings.Split(svid, ".")
	if len(sp) < 2 || len(sp[1]) <= 0 {
		return SVID{}, fmt.Errorf("svid part error")
	}
	claims, err := DecodeSegment(sp[1])
	if err != nil {
		return SVID{}, fmt.Errorf("svid claims error:%v", err)
	}
	exp := jsoniter.Get(claims, "exp").ToInt64()
	if exp <= 0 {
		return SVID{}, fmt.Errorf("svid exp error")
	}
	sub := jsoniter.Get(claims, "sub").ToString()
	if sub == "" {
		return SVID{}, fmt.Errorf("svid sub error")
	}

	diff := exp - time.Now().Unix()
	if math.Abs(float64(diff)) > float64(svidExpireTolerantTime) {
		return SVID{}, fmt.Errorf("svid expired")
	}
	return SVID{
		Exp: exp,
		Sub: sub,
	}, nil
}

func ParseFromSpiffeID(id string) map[string]string {
	m := make(map[string]string)
	if !strings.HasPrefix(id, DefaultSpiffeIDPrefix) {
		return m
	}
	attrsSeq := strings.Split(id[len(DefaultSpiffeIDPrefix):], "/")
	if len(attrsSeq)%2 != 0 {
		return m
	}
	for i := 0; i < len(attrsSeq); i += 2 {
		m[attrsSeq[i]] = attrsSeq[i+1]
	}
	return m
}

// Decode JWT specific base64url encoding with padding stripped
func DecodeSegment(seg string) ([]byte, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}

	return base64.URLEncoding.DecodeString(seg)
}
