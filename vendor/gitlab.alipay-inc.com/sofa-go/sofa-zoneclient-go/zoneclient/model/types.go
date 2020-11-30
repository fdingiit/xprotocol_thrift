package model

import "strings"

type DrMetaInfo struct {
	DrType DrCommonCode // dr type, 'ldr' or 'rdr'
	Level  DrCommonCode // dr level, 'group' or 'zone'
	DrName string       // dr zone/group name
}

// dr group/zone config info
type DrRuleInfo struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Ldr      string `json:"ldr"`
	Rdr      string `json:"rdr"`
	UidRange string `json:"uidRange"`
	IdcName  string `json:"idcName"`
}

// route rule
type FlowRule struct {
	DisasterInfo    DisasterInfo      `json:"disasterInfo"`
	Domain          map[string]string `json:"domain"`
	DrRule          []DrRuleInfo      `json:"drRule"`
	ElasticInfo     ElasticInfo       `json:"elasticInfo"`
	Enable          bool              `json:"enable"`
	FlowRuleGroupes []FlowRuleGroup   `json:"flowRuleGroupes"`
	GrayInfo        GrayInfo          `json:"grayInfo"`
	Ldc             []CityIdcInfo     `json:"ldc"`
	Rule            string            `json:"rule"`
	RuleName        string            `json:"ruleName"`
	Version         int64             `json:"version"`
	ExtraInfo       ExtraInfo         `json:"extraInfo"`
}

type FlowRules []FlowRule

type CityIdcInfo struct {
	Name string    `json:"name,omitempty"`
	Type string    `json:"type,omitempty"`
	Idcs []IdcInfo `json:"idcs"`
}

type IdcInfo struct {
	Name        string        `json:"name,omitempty"`
	Type        string        `json:"type,omitempty"`
	Zones       []IdcZoneInfo `json:"zones"`
	CityIdcInfo *CityIdcInfo  `json:"-"`
}

type IdcZoneInfo struct {
	Name   string `json:"name,omitempty"`
	Type   string `json:"type,omitempty"`
	Status string `json:"status"`
}

type ExtraInfo struct {
	Groups       []FlowRuleGroup `json:"groups"`
	DrRule       []DrRuleInfo    `json:"drRule"`
	ElasticInfo  ElasticInfo     `json:"elasticInfo"`
	DisasterInfo DisasterInfo    `json:"disasterInfo"`
	Ldc          []CityIdcInfo   `json:"ldc"`
}

type VipZoneInfo struct {
	City          string
	Idc           string
	Zone          string
	TargetAppName string
}

type GrayRouteInfo struct {
	ElasticInfo           ElasticInfo
	GrayUidMap            map[int32]string
	GrayMarkUidMap        map[int32]string
	GrayElasticUidMap     map[int32]string
	GrayElasticMarkUidMap map[int32]string
	GrayDefaultGzone      DefaultZone
}

type DisasterValue struct {
	UidDisasterStatusMap map[string]DisasterStatusEnum `json:"uidDisasterStatusMap"`
}

type DisasterItem struct {
	DisasterValueMap map[string]DisasterValue `json:"disasterValueMap"`
}

type DisasterInfo struct {
	DisasterItemMap map[string]DisasterItem `json:"disasterItemMap"`
}

type ElasticValue struct {
	ElasticValues []string `json:"elasticValues"`
	Status        string   `json:"status"`
}

type ElasticInfo struct {
	ElasticValueMap map[string]ElasticValue `json:"elasticValueMap"` // key is 'ONLINE' or 'PRESS'
}

func (ei *ElasticInfo) GetElasticStatus(isMark bool) ElasticStatusEnum {
	var status string
	if isMark {
		status = ei.ElasticValueMap[string(FT_PRESS)].Status
	} else {
		status = ei.ElasticValueMap[string(FT_ONLINE)].Status
	}

	if status == "" {
		return ES_NORMAL
	} else {
		return ElasticStatusEnum(status)
	}
}

func (ei *ElasticInfo) IsContainsEid(eid string, isMark bool) bool {
	var elasticValues []string
	if isMark {
		elasticValues = ei.ElasticValueMap[string(FT_PRESS)].ElasticValues
	} else {
		elasticValues = ei.ElasticValueMap[string(FT_ONLINE)].ElasticValues
	}

	if len(elasticValues) > 0 {
		for i := 0; i < len(elasticValues); i++ {
			if elasticValues[i] == strings.TrimSpace(eid) {
				return true
			}
		}
	}
	return false
}

type GrayInfo struct {
	Groups      []FlowRuleGroup `json:"groups"`
	ElasticInfo ElasticInfo     `json:"elasticInfo"`
}

type RouterValue struct {
	Uid      string
	Eid      string
	InGray   int32
	FlowType string // ONLINE,PRESS, identity uid first
}

type RouteContext struct {
	Uid    int32
	IsMark bool
	InGray int32
}

type Server struct {
	Url     string
	Ip      string
	Port    int32
	Timeout int32
}

type UidZone struct {
	DefaultZones []ZoneInfo
	ElasticZones []ZoneInfo
}

type DrmInfo struct {
	AppName string
	DataId  string
	Attr    string
	Value   string
	Version int
}
