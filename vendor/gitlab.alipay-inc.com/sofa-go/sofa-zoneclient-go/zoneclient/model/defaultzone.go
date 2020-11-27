package model

type DefaultZone struct {
	DefaultZone            *ZoneInfo
	WeightedRoundRobin     WeightedRoundRobin
	MarkWeightedRoundRobin WeightedRoundRobin
	ZoneMap                map[string]*ZoneInfo
}

func (dz *DefaultZone) GetZoneByWeight(isGzoneColor bool, isMark bool) *ZoneInfo {
	if dz.DefaultZone == nil ||
		dz.DefaultZone.GetRouteWeight(isMark) == 0 || (isGzoneColor && dz.DefaultZone.IsGZone()) {
		weightedRoundRobin := dz.WeightedRoundRobin
		if isMark {
			weightedRoundRobin = dz.MarkWeightedRoundRobin
		}
		zoneName := weightedRoundRobin.GetServerAsPerAlgo()
		return dz.ZoneMap[zoneName]
	} else {
		return dz.DefaultZone
	}
}

func (dz *DefaultZone) GetZoneByColor(isMark,
	isGzoneColor bool, localZone string, zoneColorMap map[string]string) *ZoneInfo {
	localZoneColor := zoneColorMap[localZone]
	// if localZone has't color, route by weight
	if localZoneColor == "" {
		return dz.GetZoneByWeight(isGzoneColor, isMark)
	}

	// if defaultZone equals localZone, or has the same color with localZone, return defaultZone
	if dz.DefaultZone != nil {
		defaultZoneColor := zoneColorMap[dz.DefaultZone.ZoneName]
		if dz.DefaultZone.ZoneName == localZone || defaultZoneColor == localZoneColor {
			return dz.DefaultZone
		}
	}

	// else return the same color zone in zoneMap
	for _, zone := range dz.ZoneMap {
		if localZoneColor == zoneColorMap[zone.ZoneName] {
			return zone
		}
	}

	// else route by weight
	return dz.GetZoneByWeight(isGzoneColor, isMark)
}
