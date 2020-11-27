package model

import (
	"strings"
)

// zone 信息
type ZoneInfo struct {
	AppName         string         `json:"-"`
	ZoneName        string         `json:"zoneName"`
	ZoneType        string         `json:"zoneType"`        // zoneType, such as (GZ|RZ|CZ)
	IdcName         string         `json:"idcName"`         // idcName of the zone
	RouteWeight     int32          `json:"routeWeigth"`     // weight in the group
	MarkRouteWeight int32          `json:"markRouteWeigth"` // press weight in the group
	IsGray          bool           `json:"isGray"`          // whether gray zone
	IsElastic       bool           `json:"-"`               // whether elastic zone
	Status          string         `json:"-"`               // running status
	ZoneDomain      string         `json:"-"`               // zone domain info, such as gz00a.alipay.com
	IdcInfo         *IdcInfo       `json:"-"`               // point to idc info
	ZoneGroup       *FlowRuleGroup `json:"-"`               // point to flow group
	DrInfo          []DrMetaInfo   `json:"-"`               // dr info
}

func (z *ZoneInfo) IsGZone() bool {
	return strings.EqualFold(z.ZoneType, "GZ")
}

func (z *ZoneInfo) IsCZone() bool {
	return strings.EqualFold(z.ZoneType, "CZ")
}

func (z *ZoneInfo) IsRZone() bool {
	return strings.EqualFold(z.ZoneType, "RZ")
}

func (z *ZoneInfo) GetRouteWeight(isMark bool) int32 {
	if isMark {
		return z.MarkRouteWeight
	} else {
		return z.RouteWeight
	}
}
