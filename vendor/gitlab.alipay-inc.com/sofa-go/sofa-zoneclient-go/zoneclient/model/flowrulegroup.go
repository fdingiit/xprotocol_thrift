package model

type FlowRuleGroup struct {
	GroupName                      string     `json:"groupName"`
	UidStringMultiRange            string     `json:"uidStringMultiRange"` // such as 0~19,40~59
	MarkUidStringMultiRange        string     `json:"markUidStringMultiRange"`
	ElasticUidStringMultiRange     string     `json:"elasticUidStringMultiRange"`
	ElasticMarkUidStringMultiRange string     `json:"elasticMarkUidStringMultiRange"`
	Zones                          []ZoneInfo `json:"zones"`
	DefaultGzone                   string     `json:"defaultGzone"`
	GroupType                      string     `json:"groupType"`
	IsElastic                      bool       `json:"isElastic"`
	UidMultiRange                  UidMultiRange
	ElasticUidMultiRange           UidMultiRange
	MarkUidMultiRange              UidMultiRange
	ElasticMarkUidMultiRange       UidMultiRange
	DefaultUidRange                UidMultiRange
	// if curGroup is elastic, this is the src group, else this is nil
	SrcGroup *FlowRuleGroup
	// if curGroup isn't elastic, this is corresponding elastic groups, else this is nil
	ElasticGroups []*FlowRuleGroup
	ZoneMap       map[string]*ZoneInfo
}

func (frg *FlowRuleGroup) Init() error {
	frg.UidMultiRange = NewUidMultiRange()
	if err := frg.UidMultiRange.BuildFromString(frg.UidStringMultiRange); err != nil {
		return err
	}

	frg.MarkUidMultiRange = NewUidMultiRange()
	if err := frg.MarkUidMultiRange.BuildFromString(frg.MarkUidStringMultiRange); err != nil {
		return err
	}

	frg.ElasticUidMultiRange = NewUidMultiRange()
	if err := frg.ElasticUidMultiRange.BuildFromString(frg.ElasticUidStringMultiRange); err != nil {
		return err
	}

	frg.ElasticMarkUidMultiRange = NewUidMultiRange()
	if err := frg.ElasticMarkUidMultiRange.BuildFromString(frg.ElasticMarkUidStringMultiRange); err != nil {
		return err
	}

	frg.ZoneMap = make(map[string]*ZoneInfo, len(frg.Zones))
	for index := range frg.Zones {
		zone := &frg.Zones[index]
		frg.ZoneMap[zone.ZoneName] = zone
	}
	return nil
}

func (frg *FlowRuleGroup) ZoneFromGroup(zoneName string) ZoneInfo {
	zoneInfo := frg.ZoneMap[zoneName]
	if zoneInfo == nil {
		return ZoneInfo{}
	}
	return *zoneInfo
}

func (frg *FlowRuleGroup) isContains(zoneName string) bool {
	return frg.ZoneMap[zoneName] != nil
}

func (frg *FlowRuleGroup) FlowAvg(ifMark bool) bool {
	b := true
	for _, z := range frg.Zones {
		if ifMark {
			if z.MarkRouteWeight != 50 {
				b = false
				break
			}
		} else {
			if z.RouteWeight != 50 {
				b = false
				break
			}
		}
	}
	return b
}

func (frg *FlowRuleGroup) FindCurUidRange(isMark bool) UidMultiRange {
	if frg.IsElastic {
		if isMark {
			return frg.ElasticMarkUidMultiRange
		} else {
			return frg.ElasticUidMultiRange
		}
	} else {
		if isMark {
			return frg.MarkUidMultiRange
		} else {
			return frg.UidMultiRange
		}
	}
}

func (frg *FlowRuleGroup) IsCzoneGroup() bool {
	return frg.GroupType == "CZG"
}

func (frg *FlowRuleGroup) IsGzoneGroup() bool {
	return frg.GroupType == "GZG"
}

func (frg *FlowRuleGroup) IsRzoneGroup() bool {
	return frg.GroupType == "RZG"
}
