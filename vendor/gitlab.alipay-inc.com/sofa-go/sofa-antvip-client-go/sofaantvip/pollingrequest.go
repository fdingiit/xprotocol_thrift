package sofaantvip

import (
	"strconv"
	"time"

	"gitlab.alipay-inc.com/sofa-go/sofa-antvip-client-go/sofaantvip/protobuf"
)

const (
	ExtensionNotExistDomains = "EXTENSION_NOT_EXIST_DOMAINS"
	ExtensionZoneInfoList    = "EXTENSION_ZONE_INFO_LIST"
)

type pollingRequest struct {
	From                      string                 `json:"from"`
	ClientVersion             string                 `json:"clientVersion"`
	AllowPolling              bool                   `json:"allowPolling"`
	DataCenter                string                 `json:"datacenter"`
	VipDomainName2ChecksumMap map[string]string      `json:"vipDomainName2ChecksumMap"`
	StartTime                 int64                  `json:"startTime"`
	AcceptTime                int64                  `json:"acceptTime"`
	ExtensionParams           map[string]interface{} `json:"extensionParams"`
}

type pollingResponse struct {
	ErrorCode       int                    `json:"errorCode"`
	ErrorMsg        string                 `json:"errorMsg"`
	NameList        []string               `json:"nameList"`
	VipDomains      []VipDomain            `json:"vipDomains"`
	StartTime       int64                  `json:"startTime"`
	AcceptTime      int64                  `json:"acceptTime"`
	ExtensionParams map[string]interface{} `json:"extensionParams"`
	zoneInfoList    *ZoneInfoList
}

func newPollingRequest(config *Config) pollingRequest {
	extParams := make(map[string]interface{}, 16)
	extParams["AppName"] = config.appName
	extParams["checksumForCompatibleSign"] = "true"
	extParams[ExtensionZoneInfoList] = "N"
	extParams = paddingExtParams(config, extParams)

	request := pollingRequest{
		From:            config.trFrom,
		ClientVersion:   config.version,
		AllowPolling:    false,
		StartTime:       time.Now().UnixNano() / 1000000,
		ExtensionParams: extParams,
		DataCenter:      config.datacenter,
	}

	return request
}

func (pr *pollingRequest) ToProtobuf() *protobuf.PollingRequestMsg {
	extensionParams := make(map[string]string)
	for k, v := range pr.ExtensionParams {
		if vs, ok := v.(string); ok {
			extensionParams[k] = vs
		}
	}

	return &protobuf.PollingRequestMsg{
		From:                      pr.From,
		ClientVersion:             pr.ClientVersion,
		AllowPolling:              pr.AllowPolling,
		VipDomainName2ChecksumMap: pr.VipDomainName2ChecksumMap,
		ExtensionParams:           extensionParams,
		StartTime:                 pr.StartTime,
		AcceptTime:                pr.AcceptTime,
	}
}

func paddingExtParams(config *Config, extParams map[string]interface{}) map[string]interface{} {
	accessKey := config.accessKey
	accessSecret := config.accessSecret
	id := config.instanceID

	// only sign when access key and secret are not empty
	if accessKey == "" || accessSecret == "" || id == "" {
		return extParams
	}

	cacheTime := int64(5 * 60 * 1000)
	timestamp := time.Now().UnixNano() / 1000000 / cacheTime * cacheTime
	content := id + strconv.FormatInt(timestamp, 10)
	extParams["AccessKey"] = accessKey
	extParams["AccessInstanceId"] = id
	extParams["AccessAlgorithm"] = "HmacSHA256"
	extParams["AccessContent"] = content
	extParams["AccessTimestamp"] = strconv.FormatInt(timestamp, 10)
	extParams["AccessSignature"] = doHMACSha256AndBase64([]byte(content), []byte(accessSecret))

	return extParams
}

func (pr *pollingResponse) Polyfill() {
	vipdomains := pr.VipDomains
	for i := range vipdomains {
		vipdomains[i].Polyfill(vipdomains[i].HealthCheckDefaultPort)
	}
}

func (pr *pollingResponse) LoadDomains() []VipDomain {
	return append(pr.VipDomains, pr.LoadDeletedDomains()...)
}

func (pr *pollingResponse) LoadDeletedDomains() []VipDomain {
	params, ok := pr.ExtensionParams[ExtensionNotExistDomains]
	if !ok {
		return nil
	}

	// ExtensionNotExistDomains maybe []stirng or string
	deletedDomains, ok := params.([]string)
	if !ok {
		domain, pok := params.(string)
		if !pok {
			return nil
		}
		deletedDomains = append(deletedDomains, domain)
	}

	vipdoamins := make([]VipDomain, 0, len(deletedDomains))
	for i := range deletedDomains {
		vipdoamins = append(vipdoamins, VipDomain{
			Name:      deletedDomains[i],
			IsDeleted: true,
		})
	}
	return vipdoamins
}

func (pr *pollingResponse) FromProtobuf(msg *protobuf.PollingResponseMsg) {
	var (
		vipDomains []VipDomain
		zi         *ZoneInfoList
	)

	params := make(map[string]interface{})

	for k, v := range msg.ExtensionParams {
		if k == ExtensionZoneInfoList {
			zi = NewZoneInfoList(msg.ExtensionParams[k])
			pr.zoneInfoList = zi
			params[k] = zi
		} else {
			params[k] = v
		}
	}

	if len(msg.VipDomains) > 0 {
		for _, vipDomainMsg := range msg.GetVipDomains() {
			var vipDomain VipDomain
			vipDomain.FromProtobuf(vipDomainMsg, zi)
			vipDomains = append(vipDomains, vipDomain)
		}
	}

	pr.VipDomains = vipDomains
	pr.ExtensionParams = params
	pr.StartTime = msg.GetStartTime()
	pr.AcceptTime = msg.GetAcceptTime()
}
