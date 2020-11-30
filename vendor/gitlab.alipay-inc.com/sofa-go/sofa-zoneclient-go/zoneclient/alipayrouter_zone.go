package zoneclient

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gitlab.alipay-inc.com/sofa-go/sofa-zoneclient-go/zoneclient/model"
)

func (al *AlipayRouter) QueryZoneInfo(zoneName string) (model.ZoneInfo, error) {
	var zoneInfo model.ZoneInfo
	zoneName = strings.ToUpper(strings.TrimSpace(zoneName))
	if zoneName != "" {
		if zone := al.getZoneRouteInfo().GetZoneByName(zoneName); zone == nil {
			return model.ZoneInfo{}, errors.New("zoneclient:  no zone found")
		} else {
			return *zone, nil
		}
	}
	return zoneInfo, nil
}

func (al *AlipayRouter) Route(rv model.RouterValue) (model.ZoneInfo, error) {
	uid, isMark, err := al.uidConvert(rv.Uid)
	if err != nil {
		return model.ZoneInfo{}, err
	}

	if !isMark && isPressFlow(rv.FlowType) {
		isMark = true
	}

	inGray := al.recomputeInGray(rv.InGray)
	if uid == -1 {
		// if uid equals -1, return defaultGzone
		return al.QueryDefaultGzone(isMark, inGray)
	}

	// else compute target zone by uid
	return al.compute(model.RouteContext{
		Uid:    uid,
		IsMark: isMark,
		InGray: inGray,
	})
}

func (al *AlipayRouter) ElasticRoute(rv model.RouterValue) (model.ZoneInfo, error) {
	eid := strings.TrimSpace(rv.Eid)
	if eid == "" {
		return model.ZoneInfo{}, errors.New("zoneclient:  eid cannot be empty")
	}

	uid, isMark, err := al.uidConvert(rv.Uid)
	if err != nil {
		return model.ZoneInfo{}, err
	}

	if !isMark && isPressFlow(rv.FlowType) {
		isMark = true
	}

	inGray := al.recomputeInGray(rv.InGray)
	elasticUIDMap := al.getZoneRouteInfo().GetUIDMap(inGray, true, isMark)
	elasticInfo := al.getZoneRouteInfo().GetElasticInfo(inGray)
	elasticStatus := elasticInfo.GetElasticStatus(isMark)
	groupMap := al.getZoneRouteInfo().GetGroupMap(inGray)

	switch elasticStatus {
	case model.ES_ELASTIC:
		if elasticUIDMap[uid] != "" && (elasticInfo.IsContainsEid(eid, isMark) || eid == model.DEFAULT_EID) &&
			eid != model.NO_ELASTIC_EID {
			group := groupMap[elasticUIDMap[uid]]
			return al.getZoneRouteInfo().RouteByColor(group, isMark)
		}
	default:
	}

	return al.compute(model.RouteContext{
		Uid:    uid,
		IsMark: isMark,
		InGray: inGray,
	})
}

func (al *AlipayRouter) GetDefaultGzone(inGray int32) (model.ZoneInfo, error) {
	defaultGzone := al.getZoneRouteInfo().GetDefaultGzone(inGray, false)
	if defaultGzone == nil {
		return model.ZoneInfo{}, errors.New("zoneclient:  no default gzone")
	}
	return *defaultGzone, nil
}

func (al *AlipayRouter) GetDefaultCzone() (model.ZoneInfo, error) {
	defaultCzone := al.getZoneRouteInfo().GetDefaultCzone(false)
	if defaultCzone == nil {
		return model.ZoneInfo{}, errors.New("zoneclient: no default czone")
	}
	return *defaultCzone, nil
}

func (al *AlipayRouter) IsRpcElastic(interfaze, version,
	uniqueId, method string, isMark bool) (model.BizElastic, error) {
	newVersion := strings.TrimSpace(version)
	if newVersion == "" {
		newVersion = "1.0"
	}
	service := fmt.Sprintf("%v:%v", interfaze, newVersion)
	if uniqueId != "" {
		service = fmt.Sprintf("%v:%v", service, uniqueId)
	}
	conditions := []string{"RPC", service, method}
	return al.getElasticRuleInfo().ElasticRuleTree.Contains(conditions, isMark), nil
}

func (al *AlipayRouter) IsAppElastic(appName string, isMark bool) (bool, error) {
	appName = strings.TrimSpace(appName)
	if appName == "" {
		return false, errors.New("zoneclient:  appname cannot be empty")
	}

	conditions := []string{"APP", appName}
	bizElastic := al.getElasticRuleInfo().ElasticRuleTree.Contains(conditions, isMark)
	return bizElastic.IsElastic, nil
}

func (al *AlipayRouter) IsGzRpc(service string) (bool, error) {
	if strings.TrimSpace(service) == "" {
		return false, errors.New("zoneclient:  service is empty")
	}
	if al.elasticRuleInfo == nil {
		return false, errors.New("elasticRule init failed")
	}
	conditions := []string{"GZ_RPC", strings.TrimSpace(service)}
	return al.elasticRuleInfo.ElasticRuleTree.Contains(conditions, false).IsElastic, nil
}

func (al *AlipayRouter) QueryAllZoneInfo() ([]model.ZoneInfo, error) {
	return al.getZoneRouteInfo().GetAllZoneInfo()
}

func (al *AlipayRouter) GetCurZone() (model.ZoneInfo, error) {
	return al.QueryZoneInfo(al.config.GetZone())
}

func (al *AlipayRouter) compute(ctx model.RouteContext) (model.ZoneInfo, error) {
	zoneInfo, _ := al.getZoneRouteInfo().DoRoute(ctx)
	if zoneInfo.ZoneName == "" {
		return al.GetDefaultGzone(ctx.InGray)
	}
	return zoneInfo, nil
}

func (al *AlipayRouter) recomputeInGray(inGray int32) int32 {
	newInGray := inGray
	localZone := al.getZoneRouteInfo().GetCurZone()
	if localZone.IsGray && inGray == model.ROUTE_IN_GRAY {
		newInGray = model.ROUTE_DEFAULT
	} else if !localZone.IsGray && inGray == model.ROUTE_OUT_GRAY {
		newInGray = model.ROUTE_DEFAULT
	}
	return newInGray
}

func (al *AlipayRouter) uidConvert(uid string) (int32, bool, error) {
	if uid == "" {
		return -1, false, errors.New("zoneclient:  uid is empty")
	}

	end := strings.ToUpper(uid[len(uid)-1:])
	isMark := strings.Contains(model.CONVERTUID_STRING, end)
	if isMark {
		index := strings.Index(model.CONVERTUID_STRING, end)
		uid = fmt.Sprintf("%v%d", uid[:len(uid)-1], index)
	}
	ret, err := strconv.ParseInt(uid, 10, 64)
	return int32(ret), isMark, err
}

func (al *AlipayRouter) GetDisasterStatus(uid string, isElastic bool) (model.DisasterStatusEnum, error) {
	uidInt, isMark, err := al.uidConvert(uid)
	if err != nil {
		return model.NORMAL, err
	}

	// check uid elastic
	if isElastic {
		UIDMap := al.getZoneRouteInfo().GetUIDMap(model.ROUTE_DEFAULT, isElastic, isMark)
		if _, ok := UIDMap[uidInt]; !ok {
			return model.NORMAL, errors.New("uid[" + uid + "] is not elastic")
		}
	}

	// get uid disaster status
	disasterStatusMap := al.getZoneRouteInfo().GetDisasterStatusMap(isMark)
	if _, ok := disasterStatusMap[isElastic]; ok {
		if _, ok = disasterStatusMap[isElastic][uidInt]; ok {
			return disasterStatusMap[isElastic][uidInt], nil
		}
	}

	return model.NORMAL, nil
}

func (al *AlipayRouter) GetElasticStatus(isMark bool) (model.ElasticStatusEnum, error) {
	elasticInfo := al.getZoneRouteInfo().GetElasticInfo(model.ROUTE_DEFAULT)

	flowType := model.FT_ONLINE
	if isMark {
		flowType = model.FT_PRESS
	}

	var status string
	if elasticInfo.ElasticValueMap != nil {
		if elasticValue, ok := elasticInfo.ElasticValueMap[string(flowType)]; ok {
			status = elasticValue.Status
		}
	}
	if status != "" {
		return model.ElasticStatusEnum(status), nil
	}
	return model.ES_NORMAL, nil
}

func (al *AlipayRouter) GetProdUidZone(uid string) (model.UidZone, error) {
	uidInt, isMark, err := al.uidConvert(uid)
	if err != nil {
		return model.UidZone{}, err
	}

	curZone := al.getZoneRouteInfo().GetCurZone()
	if curZone == nil {
		return model.UidZone{}, errors.New("current zone is nil")
	}

	zoneRouteInfo := al.zoneRouteInfo
	if curZone.IsGray {
		zoneRouteInfo = al.getZoneRouteInfo().GetExtraZoneRouteInfo()
	}

	groupMap := zoneRouteInfo.GetGroupMap(model.ROUTE_DEFAULT)

	defaultZones := make([]model.ZoneInfo, 0)
	defaultUIDMap := zoneRouteInfo.GetDefaultUIDMap()
	if groupName, ok := defaultUIDMap[uidInt]; ok {
		if group, isContains := groupMap[groupName]; isContains {
			defaultZones = group.Zones
		}
	}
	elasticZones := make([]model.ZoneInfo, 0)
	elasticUIDMap := zoneRouteInfo.GetUIDMap(model.ROUTE_DEFAULT, true, isMark)
	if groupName, ok := elasticUIDMap[uidInt]; ok {
		if group, isContains := groupMap[groupName]; isContains {
			elasticZones = group.Zones
		}
	}
	return model.UidZone{
		DefaultZones: defaultZones,
		ElasticZones: elasticZones,
	}, nil
}

func (al *AlipayRouter) GetSrcZone(zoneName string) (model.ZoneInfo, error) {
	return al.getZoneRouteInfo().GetSrcZone(zoneName, false)
}

func (al *AlipayRouter) GetUidZone(uid string) (model.UidZone, error) {
	uidInt, isMark, err := al.uidConvert(uid)
	if err != nil {
		return model.UidZone{}, err
	}

	zoneRouteInfo := al.getZoneRouteInfo()
	groupMap := zoneRouteInfo.GetGroupMap(model.ROUTE_DEFAULT)

	defaultZones := make([]model.ZoneInfo, 0)
	defaultUIDMap := zoneRouteInfo.GetDefaultUIDMap()
	if groupName, ok := defaultUIDMap[uidInt]; ok {
		if group, isContains := groupMap[groupName]; isContains {
			defaultZones = group.Zones
		}
	}
	elasticZones := make([]model.ZoneInfo, 0)
	elasticUIDMap := zoneRouteInfo.GetUIDMap(model.ROUTE_DEFAULT, true, isMark)
	if groupName, ok := elasticUIDMap[uidInt]; ok {
		if group, isContains := groupMap[groupName]; isContains {
			elasticZones = group.Zones
		}
	}
	return model.UidZone{
		DefaultZones: defaultZones,
		ElasticZones: elasticZones,
	}, nil
}

func (al *AlipayRouter) GetZoneNames() []string {
	return al.getZoneRouteInfo().GetAllZoneNames()
}

func (al *AlipayRouter) IsElastic(uid string) (bool, error) {
	uidInt, isMark, err := al.uidConvert(uid)
	if err != nil {
		return false, err
	}

	elasticUIDMap := al.getZoneRouteInfo().GetUIDMap(model.ROUTE_DEFAULT, true, isMark)
	if _, ok := elasticUIDMap[uidInt]; ok {
		return true, nil
	}

	return false, nil
}

func (al *AlipayRouter) QueryOriginZone(zoneName, flowType string) (model.ZoneInfo, error) {
	isMark := isPressFlow(flowType)
	return al.getZoneRouteInfo().GetSrcZone(zoneName, isMark)
}

func (al *AlipayRouter) IsUidElastic(uid string, flowType string) (bool, error) {
	uidInt, isMark, err := al.uidConvert(uid)
	if err != nil {
		return false, err
	}

	if !isMark && isPressFlow(flowType) {
		isMark = true
	}

	elasticUIDMap := al.getZoneRouteInfo().GetUIDMap(model.ROUTE_DEFAULT, true, isMark)
	if _, ok := elasticUIDMap[uidInt]; ok {
		return true, nil
	}

	return false, nil
}

func (al *AlipayRouter) QueryDefaultGzone(isMark bool, inGray int32) (model.ZoneInfo, error) {
	defaultGzone := al.getZoneRouteInfo().GetDefaultGzone(inGray, isMark)
	if defaultGzone == nil {
		return model.ZoneInfo{}, errors.New("zoneclient:  no default gzone")
	}
	return *defaultGzone, nil
}

func (al *AlipayRouter) QueryDefaultCzone(flowType string) (model.ZoneInfo, error) {
	isMark := isPressFlow(flowType)
	defaultCzone := al.getZoneRouteInfo().GetDefaultCzone(isMark)
	if defaultCzone == nil {
		return model.ZoneInfo{}, errors.New("zoneclient: no default czone")
	}
	return *defaultCzone, nil
}

func (al *AlipayRouter) GetCloselyRzone(filterElasticZone bool, flowType string) (model.ZoneInfo, error) {
	var isPress = isPressFlow(flowType)
	return al.zoneRouteInfo.GetCloselyRzone(filterElasticZone, isPress)
}

func (al *AlipayRouter) GetElasticLabel(eid, interfaze, version, uniqueId, method string, isMark bool) (string, error) {
	elasticInfo := al.getZoneRouteInfo().GetElasticInfo(0)
	elasticStatus := elasticInfo.GetElasticStatus(isMark)
	if elasticStatus != model.ES_ELASTIC {
		return model.NO_ELASTIC, nil
	}

	bizElastic, err := al.IsRpcElastic(interfaze, version, uniqueId, method, isMark)
	if err != nil {
		al.logger.Errorf(fmt.Sprintf("zoneclient: failed to get elastic label, %v", err))
		return model.NO_ELASTIC, err
	}

	if bizElastic.IsElastic {
		if bizElastic.IsUseDefaultEid {
			return model.ELASTIC_DEFAULT, nil
		} else if strings.TrimSpace(eid) != "" {
			return model.ELASTIC_UNDEFAULT, nil
		} else {
			return model.NO_CONSISTENT, nil
		}
	}

	return model.NO_ELASTIC, nil
}
