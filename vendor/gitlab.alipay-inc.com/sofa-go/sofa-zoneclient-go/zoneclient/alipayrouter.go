package zoneclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"gitlab.alipay-inc.com/sofa-go/sofa-zoneclient-go/zoneclient/model"

	"gitlab.alipay-inc.com/sofa-go/sofa-drm-client-go/sofadrm"
	"gitlab.alipay-inc.com/sofa-go/sofa-logger-go/sofalogger"
)

const (
	defaultAlipayRouterAntvipPort = 80
)

// nolint
const (
	ZoneInfoDRMDataIdFormat           = "Alipay.%v:name=com.alipay.routeclient.DefaultRouteCoordinator.zoneInfo,version=3.0@DRM"
	ElasticRuleVersionDRMDataIdFormat = "Alipay.%v:name=com.alipay.routeclient.DefaultRouteCoordinator.elasticRuleVersion,version=3.0@DRM"
	ZoneColorDRMDataIdFormat          = "Alipay.%v:name=com.alipay.routeclient.DefaultRouteCoordinator.zoneColor,version=3.0@DRM"
	WhiteListRPCLoadTestDataIDFormat  = "Alipay.%v:name=com.alipay.routeclient.DefaultRouteCoordinator.whiteListRPCLoadTest,version=3.0@DRM" // R->G load test white list
)

type AlipayRouter struct {
	sync.RWMutex
	elasticRuleInfo         *model.ElasticRuleInfo
	zoneRouteInfo           *model.ZoneRouteInfo
	whitelistRPCLoadTestMap map[string]bool
	zoneColorMap            map[string]string
	locator                 Locator

	config    *Config
	logger    sofalogger.Logger
	drm       *sofadrm.Client
	dataIdMap map[string]*model.DrmInfo
}

func NewAlipayRouter(options ...AlipayRouterOptionSetter) (*AlipayRouter, error) {
	al := &AlipayRouter{
		zoneRouteInfo: &model.ZoneRouteInfo{},
		elasticRuleInfo: &model.ElasticRuleInfo{
			Version: -1,
		},
		zoneColorMap: map[string]string{},
		dataIdMap:    make(map[string]*model.DrmInfo),
	}

	for i := range options {
		options[i].set(al)
	}

	if err := al.polyfill(); err != nil {
		return nil, err
	}

	err := al.locator.RefreshServers()
	if err != nil {
		return nil, err
	}

	if err = al.getAndUpdateZoneRouteInfo(); err != nil {
		return nil, err
	}

	if err = al.getAndUpdateElasticRuleInfo(""); err != nil {
		return nil, err
	}

	al.doDRMListen()

	return al, nil
}

func (al *AlipayRouter) OnDRMPush(dataID, value string) {
	if len(value) >= 256 {
		al.logger.Infof("drm push dataID=%s value=%s...", dataID, value[:256])
	} else {
		al.logger.Infof("drm push dataID=%s value=%s", dataID, value)
	}

	needUpdate, drmInfo := al.needUpdate(dataID, value)
	if !needUpdate {
		al.logger.Infof("don't need update drm dataId %s", dataID)
		return
	} else {
		al.logger.Infof("drm %s enabled, version = %v, value = %s", drmInfo.Attr, drmInfo.Version, drmInfo.Value)
	}

	var err error

	switch drmInfo.Attr {
	case model.ZONE_INFO:
		err = al.updateZoneRouteInfo(drmInfo.Value)
	case model.ZONE_COLOR:
		err = al.updateZoneColorInfo(drmInfo.Value)
	case model.ELASTIC_RULE:
		err = al.getAndUpdateElasticRuleInfo(drmInfo.Value)
	case model.WHITE_LIST_RPC_LOAD_TEST:
		al.updateWhiteListRPCLoadTestInfo(drmInfo.Value)
	}

	if err != nil {
		al.logger.Errorf("failed to upgrade from drm dataID=%s: %v", dataID, err)
	}
}

func (al *AlipayRouter) doDRMListen() {
	appName := al.config.appName
	// init app drm resource
	if strings.TrimSpace(appName) != "" {
		al.initDrmDataId(model.ZONE_INFO, appName)
		al.initDrmDataId(model.ELASTIC_RULE, appName)
		al.initDrmDataId(model.ZONE_COLOR, appName)
		al.initDrmDataId(model.WHITE_LIST_RPC_LOAD_TEST, appName)
	}
	// init global drm resource
	al.initDrmDataId(model.ZONE_INFO, model.DRM_DOMAIN)
	al.initDrmDataId(model.ELASTIC_RULE, model.DRM_DOMAIN)
	al.initDrmDataId(model.ZONE_COLOR, model.ZONE_COLOR)
	al.initDrmDataId(model.WHITE_LIST_RPC_LOAD_TEST, model.WHITE_LIST_RPC_LOAD_TEST)

	// add drm listener
	for dataID := range al.dataIdMap {
		al.drm.AddListener(dataID, al)
		_, _, err := al.drm.GetValue(dataID)
		if err != nil {
			al.logger.Errorf("failed to get %s from drm: %s", dataID, err.Error())
		} else {
			al.logger.Infof("fetch dataID=%s from drm", dataID)
		}
	}
}

func (al *AlipayRouter) needUpdate(dataID, value string) (bool, *model.DrmInfo) {
	if drmInfo := al.dataIdMap[dataID]; drmInfo != nil {
		// 是否默认drm
		globalDrmDataId := buildGlobalDrmDataId(drmInfo.Attr)
		appDrmDataId, _ := buildDrmDataId(drmInfo.Attr, al.config.appName)
		globalDrmInfo := al.dataIdMap[globalDrmDataId]
		if globalDrmInfo == nil {
			al.logger.Errorf("failed to get global drm info, dataId = %s, global dataId = %s", dataID, globalDrmDataId)
			return false, nil
		}
		appDrmInfo := al.dataIdMap[appDrmDataId]

		switch dataID {
		case globalDrmDataId:
			globalDrmInfo.Value = value

			if appDrmInfo == nil || strings.TrimSpace(appDrmInfo.Value) == "" {
				return true, globalDrmInfo
			} else {
				return false, appDrmInfo
			}
		case appDrmDataId:
			if appDrmInfo != nil {
				appDrmInfo.Value = value
			}

			if appDrmInfo == nil || strings.TrimSpace(appDrmInfo.Value) == "" {
				return true, globalDrmInfo
			} else {
				return true, appDrmInfo
			}
		}
	}

	al.logger.Errorf("unregistered drm dataID=%s", dataID)
	return false, nil
}

func (al *AlipayRouter) initDrmDataId(attr, appName string) {
	drmDataId, err := buildDrmDataId(attr, appName)
	if err != nil {
		al.logger.Errorf("[zoneclient] build zoneInfo dataId %s failed,%v", drmDataId, err)
	} else {
		al.dataIdMap[drmDataId] = &model.DrmInfo{
			DataId:  drmDataId,
			Attr:    attr,
			AppName: appName,
			Version: -1,
		}
	}
}

func (al *AlipayRouter) polyfill() error {
	if al.logger == nil {
		al.logger = sofalogger.StdoutLogger
	}

	if al.config == nil {
		return errors.New("zoneclient: config is nil")
	}

	if al.locator == nil {
		return errors.New("zoneclient: locator is nil")
	}

	if al.drm == nil {
		return errors.New("zoneclient: drm client is nil")
	}

	return nil
}

func (al *AlipayRouter) updateZoneColorInfo(value string) error {
	zoneColorMap := make(map[string]string)
	zoneColorStr := strings.TrimSpace(value)
	if zoneColorStr != "" {
		for _, zoneColor := range strings.Split(zoneColorStr, ";") {
			if zoneColor == "" {
				continue
			}
			zcArray := strings.Split(zoneColor, ":")
			if len(zcArray) != 2 {
				return errors.New("zoneclient: invalid zone color format")
			}

			color := strings.TrimSpace(zcArray[0])
			zoneArray := strings.Split(zcArray[1], ",")
			for i := range zoneArray {
				zone := strings.TrimSpace(zoneArray[i])
				zoneColorMap[zone] = color
			}
		}
	}

	al.Lock()
	al.zoneColorMap = zoneColorMap
	al.zoneRouteInfo.SetZoneColorMap(zoneColorMap)
	al.Unlock()

	return nil
}

func (al *AlipayRouter) MarshalJSON() ([]byte, error) {
	type Status struct {
		ElasticRuleInfo         *model.ElasticRuleInfo `json:"elastic_rule_info"`
		ZoneRouteInfo           *model.ZoneRouteInfo   `json:"zone_route_info"`
		WhitelistRPCLoadTestMap map[string]bool        `json:"whitelist_rpcloadtest_map"`
		ZoneColorMap            map[string]string      `json:"zone_color_map"`
		Servers                 []string               `json:"servers"`
		Config                  *Config                `json:"config"`
	}

	al.RLock()
	defer al.RUnlock()

	s := &Status{
		ElasticRuleInfo:         al.elasticRuleInfo,
		ZoneRouteInfo:           al.zoneRouteInfo,
		WhitelistRPCLoadTestMap: al.whitelistRPCLoadTestMap,
		ZoneColorMap:            al.zoneColorMap,
		Config:                  al.config,
	}

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(s); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (al *AlipayRouter) GetZoneColorMap() map[string]string {
	al.RLock()
	defer al.RUnlock()
	return al.zoneColorMap
}

func (al *AlipayRouter) updateWhiteListRPCLoadTestInfo(value string) {
	value = strings.TrimSpace(value)
	whitelistPRCLoadTestMap := make(map[string]bool, 32)
	if value != "" {
		whiteList := strings.Split(value, ",")
		for i := range whiteList {
			service := strings.ToLower(strings.TrimSpace(whiteList[i]))
			if service != "" {
				whitelistPRCLoadTestMap[service] = true
			}
		}
	}
	al.Lock()
	al.whitelistRPCLoadTestMap = whitelistPRCLoadTestMap
	al.Unlock()
}

func (al *AlipayRouter) getAndUpdateZoneRouteInfo() error {
	var (
		server model.Server
		ok     bool
	)

	server, ok = al.locator.GetRandomServer()
	if !ok {
		return errors.New("zoneclient: no available servers")
	}

	data, err := al.getZoneRouteInfoFromHTTPOrFile(fmt.Sprintf("%v:%v", server.Ip, server.Port),
		al.config.GetAlipayRouterConfig().GetZoneRoutePath())
	if err != nil {
		return err
	}
	return al.updateZoneRouteInfo(string(data))
}

func (al *AlipayRouter) getAndUpdateElasticRuleInfo(version string) error {
	var (
		intv int64
		err  error
	)
	if version != "" {
		intv, err = strconv.ParseInt(version, 10, 64)
		if err != nil {
			return err
		}
	}

	al.Lock()
	if al.elasticRuleInfo.Version == intv {
		al.Unlock()
		return nil
	}
	al.Unlock()

	var (
		server model.Server
		ok     bool
	)

	server, ok = al.locator.GetRandomServer()
	if !ok {
		return errors.New("zoneclient: no available servers")
	}

	data, err := al.getElasticRuleInfoFromHTTPOrFile(version, fmt.Sprintf("%v:%v", server.Ip, server.Port),
		al.config.GetAlipayRouterConfig().GetElasticRulePath())
	if err != nil {
		return err
	}
	return al.updateElasticRuleInfo(string(data))
}

func (al *AlipayRouter) getZoneRouteInfoFromHTTPOrFile(server, path string) (data []byte, err error) {
	url := fmt.Sprintf("http://%s%s%s", server, model.HttpZoneInfoUrl, al.config.GetZone())
	return al.getFromHTTPOrFile(url, path, false)
}

func (al *AlipayRouter) getElasticRuleInfoFromHTTPOrFile(version, server, path string) (data []byte, err error) {
	url := fmt.Sprintf("http://%s%s%s", server, model.HttpElasticRuleUrl, version)
	return al.getFromHTTPOrFile(url, path, true)
}

func (al *AlipayRouter) getFromHTTPOrFile(url, path string, gzip bool) (data []byte, err error) {
	defer func() {
		if path == "" {
			return
		}

		if err != nil { // try load from file when occurs error
			data, err = readFile(path)
			if err == nil {
				al.logger.Infof("load data from file: %s", path)
			}

		} else { // write to file
			err = writeFile(path, data)
			if err == nil {
				al.logger.Infof("write data to file: %s", path)
			}
		}
	}()

	data, err = al.getHTTPRequestBodyFromURL(url)
	if err != nil {
		return nil, err
	}

	if gzip {
		data, err = gunzip(data)
		if err != nil {
			return nil, err
		}
	}

	return data, err
}

func (al *AlipayRouter) updateElasticRuleInfo(elasticRule string) error {
	elasticInfo, err := model.ParseElasticRule(elasticRule)
	if err != nil {
		return err
	}

	al.RLock()
	al.elasticRuleInfo = elasticInfo
	al.RUnlock()
	return nil
}

func (al *AlipayRouter) updateZoneRouteInfo(zoneInfo string) error {
	zoneRouteInfo, err := model.ParseZoneRouteInfo(al.config.GetAppName(), al.config.GetZone(), zoneInfo)
	if err != nil {
		return err
	}

	al.Lock()
	zoneRouteInfo.SetZoneColorMap(al.zoneColorMap)
	al.zoneRouteInfo = zoneRouteInfo
	al.Unlock()

	return nil
}

func (al *AlipayRouter) getElasticRuleInfo() *model.ElasticRuleInfo {
	al.RLock()
	defer al.RUnlock()
	return al.elasticRuleInfo
}

func (al *AlipayRouter) getZoneRouteInfo() *model.ZoneRouteInfo {
	al.RLock()
	defer al.RUnlock()
	return al.zoneRouteInfo
}

func (al *AlipayRouter) getHTTPRequestBodyFromURL(url string) ([]byte, error) {
	resp, err := al.getHTTPClient().Get(url)
	if err != nil {
		return nil, err
	}
	// nolint
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (al *AlipayRouter) getHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   al.config.GetAlipayRouterConfig().GetTimeout(),
				KeepAlive: 0,
				DualStack: false,
			}).DialContext,
			DisableKeepAlives:   true, // short-lived connection
			MaxIdleConns:        0,
			TLSHandshakeTimeout: al.config.GetAlipayRouterConfig().GetTimeout(),
		},
		Timeout: al.config.GetAlipayRouterConfig().GetTimeout(),
	}
}

func buildGlobalDrmDataId(attr string) string {
	var globalDrmDataId string
	switch attr {
	case model.ZONE_INFO:
		globalDrmDataId, _ = buildDrmDataId(model.ZONE_INFO, model.DRM_DOMAIN)
	case model.ELASTIC_RULE:
		globalDrmDataId, _ = buildDrmDataId(model.ELASTIC_RULE, model.DRM_DOMAIN)
	case model.ZONE_COLOR:
		globalDrmDataId, _ = buildDrmDataId(model.ZONE_COLOR, model.ZONE_COLOR)
	case model.WHITE_LIST_RPC_LOAD_TEST:
		globalDrmDataId, _ = buildDrmDataId(model.WHITE_LIST_RPC_LOAD_TEST, model.WHITE_LIST_RPC_LOAD_TEST)
	}

	return globalDrmDataId
}

func buildDrmDataId(attr, appName string) (string, error) {
	if strings.TrimSpace(attr) == "" {
		return "", fmt.Errorf("[zoneclient] attr is blank")
	}

	switch attr {
	case model.ZONE_INFO:
		return fmt.Sprintf(ZoneInfoDRMDataIdFormat, appName), nil
	case model.ELASTIC_RULE:
		return fmt.Sprintf(ElasticRuleVersionDRMDataIdFormat, appName), nil
	case model.ZONE_COLOR:
		return fmt.Sprintf(ZoneColorDRMDataIdFormat, appName), nil
	case model.WHITE_LIST_RPC_LOAD_TEST:
		return fmt.Sprintf(WhiteListRPCLoadTestDataIDFormat, appName), nil
	}
	return "", fmt.Errorf("[zoneclient] build drm dataId failed, %v, %v", attr, appName)
}
