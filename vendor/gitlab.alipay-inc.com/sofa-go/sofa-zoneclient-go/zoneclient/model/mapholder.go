package model

type GroupMapHolder struct {
	GroupMap      map[string]*FlowRuleGroup
	GrayGroupMap  map[string]*FlowRuleGroup
	ExtraGroupMap map[string]*FlowRuleGroup
}

func (gmh *GroupMapHolder) GetGroupMap(inGray int32) map[string]*FlowRuleGroup {
	switch inGray {
	case ROUTE_IN_GRAY:
		return gmh.GrayGroupMap
	case ROUTE_OUT_GRAY:
		return gmh.ExtraGroupMap
	default:
		return gmh.GroupMap
	}
}

type ZoneMapHolder struct {
	ZoneMap      map[string]*ZoneInfo
	GrayZoneMap  map[string]*ZoneInfo
	ExtraZoneMap map[string]*ZoneInfo
}

func (zmh *ZoneMapHolder) GetZoneMap(inGray int32) map[string]*ZoneInfo {
	switch inGray {
	case ROUTE_IN_GRAY:
		return zmh.GrayZoneMap
	case ROUTE_OUT_GRAY:
		return zmh.ExtraZoneMap
	default:
		return zmh.ZoneMap
	}
}

func (zmh *ZoneMapHolder) GetByZoneName(zoneName string) *ZoneInfo {
	zoneInfo := zmh.ZoneMap[zoneName]
	if zoneInfo == nil {
		zoneInfo = zmh.GrayZoneMap[zoneName]
	}
	if zoneInfo == nil {
		zoneInfo = zmh.ExtraZoneMap[zoneName]
	}

	return zoneInfo
}
