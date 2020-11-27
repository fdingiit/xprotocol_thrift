package model

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
)

type ZoneRouteInfo struct {
	flowRules         []FlowRule
	appName           string
	version           int64
	cities            []CityIdcInfo
	zoneNames         []string
	idcMap            map[string]*IdcInfo    // key is idcName, value is idcInfo
	vipZoneMap        map[string]VipZoneInfo // key is zoneName, value is vipZoneInfo
	domainInfoMap     map[string]string      // key is zoneName, value is domain
	uidMap            map[int32]string       // key is uid, value is groupName
	markUidMap        map[int32]string       // key is mark uid, value is groupName
	elasticUidMap     map[int32]string       // key is elastic uid, value is groupName
	elasticMarkUidMap map[int32]string       // key is elastic mark uid, value is groupName
	defaultUidMap     map[int32]string       // key is default uid, value is groupName
	drRuleMap         map[string]DrRuleInfo  // key is groupName, value is drRuleInfo

	// key1 is isElastic, key2 is uid, value2 is disasterStatus
	disasterStatusMap map[bool]map[int32]DisasterStatusEnum
	// key1 is isElastic, key2 is mark uid, value2 is disasterStatus
	markDisasterStatusMap       map[bool]map[int32]DisasterStatusEnum
	zoneStatusMap               map[string]string // key is zoneName, value is zoneStatus
	zoneColorMap                map[string]string // key is zoneName, value is color
	groupMapHolder              GroupMapHolder    // contains groups, grayGroups, extraGroups
	zoneMapHolder               ZoneMapHolder     // contains zones, grayZones, extraZones
	elasticInfo                 ElasticInfo
	curGroup                    *FlowRuleGroup
	curZone                     ZoneInfo
	defaultCZone                DefaultZone
	defaultGZone                DefaultZone
	isGzoneColor                bool
	grayRouteInfo               GrayRouteInfo  // gray zone route info
	extraRouteInfo              *ZoneRouteInfo // other zone route info
	isEnable                    bool
	mod                         int32
	groupWeightedRoundRobin     map[string]WeightedRoundRobin
	markGroupWeightedRoundRobin map[string]WeightedRoundRobin
}

func NewZoneRouteInfo(appName string) ZoneRouteInfo {
	return ZoneRouteInfo{
		appName:               appName,
		zoneNames:             make([]string, 0, 100),
		idcMap:                make(map[string]*IdcInfo),
		vipZoneMap:            make(map[string]VipZoneInfo),
		domainInfoMap:         make(map[string]string),
		uidMap:                make(map[int32]string),
		markUidMap:            make(map[int32]string),
		elasticUidMap:         make(map[int32]string),
		elasticMarkUidMap:     make(map[int32]string),
		defaultUidMap:         make(map[int32]string),
		drRuleMap:             make(map[string]DrRuleInfo),
		disasterStatusMap:     make(map[bool]map[int32]DisasterStatusEnum),
		markDisasterStatusMap: make(map[bool]map[int32]DisasterStatusEnum),
		zoneStatusMap:         make(map[string]string),
		zoneColorMap:          make(map[string]string),
		groupMapHolder: GroupMapHolder{
			GroupMap:      make(map[string]*FlowRuleGroup),
			GrayGroupMap:  make(map[string]*FlowRuleGroup),
			ExtraGroupMap: make(map[string]*FlowRuleGroup),
		},
		zoneMapHolder: ZoneMapHolder{
			ZoneMap:      make(map[string]*ZoneInfo),
			GrayZoneMap:  make(map[string]*ZoneInfo),
			ExtraZoneMap: make(map[string]*ZoneInfo),
		},
		groupWeightedRoundRobin:     make(map[string]WeightedRoundRobin),
		markGroupWeightedRoundRobin: make(map[string]WeightedRoundRobin),
	}
}

// ############################################
// parse functions
// ############################################

// parse domainMap
func (r *ZoneRouteInfo) ParseDomain(domainMap map[string]string) {
	if domainMap == nil {
		r.domainInfoMap = make(map[string]string)
	} else {
		r.domainInfoMap = domainMap
	}

	for zoneName := range r.domainInfoMap {
		r.zoneNames = append(r.zoneNames, zoneName)
	}
}

// parse ldc
// city -> idc -> zone
func (r *ZoneRouteInfo) ParseLdc(cities []CityIdcInfo) {
	if len(cities) == 0 {
		return
	}

	r.cities = cities
	idcMap := r.idcMap
	vipZoneMap := r.vipZoneMap
	zoneStatusMap := r.zoneStatusMap
	// cities
	for cityIndex := range cities {
		city := &cities[cityIndex]
		if len(city.Idcs) > 0 {
			// idcs
			for idcIndex := range city.Idcs {
				idc := &city.Idcs[idcIndex]

				idc.CityIdcInfo = city
				zones := idc.Zones
				if len(zones) > 0 {
					// zones
					for _, zone := range zones {
						zoneStatusMap[zone.Name] = zone.Status
						vipZoneMap[zone.Name] = VipZoneInfo{
							City: city.Name,
							Idc:  idc.Name,
							Zone: zone.Name,
						}
					}
				}

				idcMap[idc.Name] = idc
			}
		}
	}
}

// parse flowRuleGroups
func (r *ZoneRouteInfo) ParseFlowRuleGroup(localZone string, flowRuleGroups []FlowRuleGroup) error {
	if len(flowRuleGroups) == 0 {
		return errors.New("flowRuleGroups is empty")
	}

	uidMap := r.uidMap
	markUidMap := r.markUidMap
	elasticUidMap := r.elasticUidMap
	elasticMarkUidMap := r.elasticMarkUidMap

	idcMap := r.idcMap
	groupMap := r.groupMapHolder.GroupMap
	zoneMap := r.zoneMapHolder.ZoneMap
	domainMap := r.domainInfoMap
	zoneStatusMap := r.zoneStatusMap

	for gIndex := range flowRuleGroups {
		group := &flowRuleGroups[gIndex]
		groupName := group.GroupName
		if err := group.Init(); err != nil {
			// TODO: I don't know whether continue or not
			// keept it same as mosn.
		}
		groupMap[group.GroupName] = group

		// uid assign
		mapCopy(uidMap, group.UidMultiRange.CvtToUidMap(groupName))
		mapCopy(markUidMap, group.MarkUidMultiRange.CvtToUidMap(groupName))
		if group.IsElastic {
			mapCopy(elasticUidMap, group.ElasticUidMultiRange.CvtToUidMap(groupName))
			mapCopy(elasticMarkUidMap, group.ElasticMarkUidMultiRange.CvtToUidMap(groupName))
		}

		for zIndex := range group.Zones {
			zone := &group.Zones[zIndex]

			zoneName := zone.ZoneName
			zone.ZoneGroup = group
			zone.AppName = r.appName
			zone.IsElastic = group.IsElastic
			if zone.IdcName != "" {
				idcInfo := idcMap[zone.IdcName]
				zone.IdcInfo = idcInfo
				zone.ZoneDomain = domainMap[zoneName]
				zone.Status = zoneStatusMap[zoneName]
			}
			zoneMap[zoneName] = zone

			// whether cur zone
			if localZone == zoneName {
				r.curGroup = group
				r.curZone = *zone
			}
		}
	}

	if r.curGroup == nil {
		return errors.New("zoneclient: local zone is not in flowRuleGroups")
	}
	return nil
}

// parse drRules
func (r *ZoneRouteInfo) ParseDrRule(drRules []DrRuleInfo) {
	if len(drRules) == 0 {
		return
	}

	drRuleMap := r.drRuleMap
	defaultUidMap := r.defaultUidMap
	groupMap := r.groupMapHolder.GroupMap

	for index := range drRules {
		drRule := &drRules[index]
		level := DrCommonCode(drRule.Type)
		ldrName := drRule.Ldr
		rdrName := drRule.Rdr
		groupName := drRule.Name
		uidRange := drRule.UidRange
		drInfos := make([]DrMetaInfo, 0)

		if ldrName != "" {
			drMetaInfo := DrMetaInfo{
				DrName: ldrName,
				DrType: Ldr,
				Level:  level,
			}
			drInfos = append(drInfos, drMetaInfo)
		}

		if rdrName != "" {
			drMetaInfo := DrMetaInfo{
				DrName: rdrName,
				DrType: Rdr,
				Level:  level,
			}
			drInfos = append(drInfos, drMetaInfo)
		}

		if groupName != "" {
			drRuleMap[groupName] = *drRule

			if ruleGroup, ok := groupMap[groupName]; ok {
				// set elastic group
				r.setElasticGroups(ruleGroup, rdrName)

				// set drInfo
				for zIndex := range groupMap[groupName].Zones {
					ruleGroup.Zones[zIndex].DrInfo = drInfos
				}

				// set default uid range
				// nolint
				ruleGroup.DefaultUidRange.BuildFromString(uidRange)
				tempUidMap := ruleGroup.DefaultUidRange.CvtToUidMap(groupName)
				mapCopy(defaultUidMap, tempUidMap)
			}
		}
	}
}

// parse elasticInfo
func (r *ZoneRouteInfo) ParseElasticInfo(elasticInfo ElasticInfo) {
	if len(elasticInfo.ElasticValueMap) == 0 {
		r.elasticInfo = initElasticInfo()
	} else {
		r.elasticInfo = elasticInfo
	}
}

// set defaultCzone
func (r *ZoneRouteInfo) SetDefaultCZone() error {
	curZone := r.curZone
	curGroup := r.curGroup
	zoneMap := r.zoneMapHolder.ZoneMap
	groupMap := r.groupMapHolder.GroupMap

	var czone *ZoneInfo
	var czoneGroup *FlowRuleGroup
	if curZone.IsCZone() {
		// if curZone is czone, defualtCzone is curZone
		czone = &curZone
		czoneGroup = curGroup
	} else {
		curIdc := curZone.IdcInfo
		// if current idc has czone, defaultCzone is the czone in current idc (non-elastic is preferred)
		for index := range curIdc.Zones {
			idcZone := curIdc.Zones[index]
			if strings.EqualFold(idcZone.Type, "CZ") &&
				!(curZone.Status == string(ZS_RUNNING) && idcZone.Status == string(ZS_BUILDING)) {
				zone := zoneMap[idcZone.Name]
				if czone == nil || czone.IsElastic {
					czone = zone
					czoneGroup = groupMap[zone.ZoneGroup.GroupName]
				}
			}
		}
		if czone == nil {
			// if current idc has not czone, defaultCzone is the czone in current city (non-elastic is preferred)
			localCityName := curIdc.CityIdcInfo.Name
			groupMap := r.groupMapHolder.GroupMap
			for _, group := range groupMap {
				if group.IsCzoneGroup() && len(group.Zones) > 0 {
					zones := group.Zones
					for _, zone := range zones {
						if strings.EqualFold(zone.IdcInfo.CityIdcInfo.Name, localCityName) &&
							!(curZone.Status == string(ZS_RUNNING) && zone.Status == string(ZS_BUILDING)) &&
							(czoneGroup == nil || czoneGroup.IsElastic) {
							czoneGroup = group
						}
					}
				}
			}
		}
	}

	if czoneGroup == nil || czoneGroup.GroupName == "" {
		return errors.New("defaultCzone is not exist")
	}

	r.defaultCZone = DefaultZone{
		DefaultZone:            czone,
		WeightedRoundRobin:     NewWeightedRoundRobin(czoneGroup.Zones, false),
		MarkWeightedRoundRobin: NewWeightedRoundRobin(czoneGroup.Zones, true),
		ZoneMap:                czoneGroup.ZoneMap,
	}

	return nil
}

// set defaultGzone
func (r *ZoneRouteInfo) SetDefaultGZone() error {
	var defaultGzone *ZoneInfo
	var defaultGroup *FlowRuleGroup
	curZone := r.curZone
	curGroup := r.curGroup
	groupMap := r.groupMapHolder.GroupMap
	defaultGroup = groupMap[curGroup.DefaultGzone]
	// if curZone is gzone, defaultGzone is curZone
	if curZone.IsGZone() {
		defaultGzone = &curZone
	}

	if defaultGroup == nil {
		return errors.New("defaultGzone is not exist")
	}

	r.defaultGZone = DefaultZone{
		DefaultZone:            defaultGzone,
		WeightedRoundRobin:     NewWeightedRoundRobin(defaultGroup.Zones, false),
		MarkWeightedRoundRobin: NewWeightedRoundRobin(defaultGroup.Zones, true),
		ZoneMap:                defaultGroup.ZoneMap,
	}
	return nil
}

// set appName
func (r *ZoneRouteInfo) SetAppName(appName string) {
	r.appName = appName
}

// set version
func (r *ZoneRouteInfo) SetVersion(version int64) {
	r.version = version
}

// set isEnable
func (r *ZoneRouteInfo) SetIsEnable(isEnable bool) {
	r.isEnable = isEnable
}

// set zoneColorMap
func (r *ZoneRouteInfo) SetZoneColorMap(zoneColorMap map[string]string) {
	r.isGzoneColor = false
	for zoneName := range zoneColorMap {
		if zone := r.zoneMapHolder.GetByZoneName(zoneName); zone != nil && zone.IsGZone() {
			if zone.IsGZone() {
				r.isGzoneColor = true
			}
		}
	}
	r.zoneColorMap = zoneColorMap
}

// set mod
func (r *ZoneRouteInfo) SetMod(rule string) {
	var mod int64 = 100
	values := strings.Split(rule, "%")
	if len(values) == 2 {
		mod, _ = strconv.ParseInt(values[1], 10, 64)
	}
	r.mod = int32(mod)
}

// init groupWeightedRoundRobin
// this func must be invoked after ParseFlowRuleGroup
func (r *ZoneRouteInfo) InitWeightedRoundRobin() {
	for _, group := range r.groupMapHolder.GroupMap {
		r.groupWeightedRoundRobin[group.GroupName] = NewWeightedRoundRobin(group.Zones, false)
		r.markGroupWeightedRoundRobin[group.GroupName] = NewWeightedRoundRobin(group.Zones, true)
	}
	for _, group := range r.groupMapHolder.GrayGroupMap {
		r.groupWeightedRoundRobin[group.GroupName] = NewWeightedRoundRobin(group.Zones, false)
		r.markGroupWeightedRoundRobin[group.GroupName] = NewWeightedRoundRobin(group.Zones, true)
	}
	for _, group := range r.groupMapHolder.ExtraGroupMap {
		r.groupWeightedRoundRobin[group.GroupName] = NewWeightedRoundRobin(group.Zones, false)
		r.markGroupWeightedRoundRobin[group.GroupName] = NewWeightedRoundRobin(group.Zones, true)
	}
}

// set flowRules
func (r *ZoneRouteInfo) SetFlowRules(flowRules []FlowRule) {
	r.flowRules = flowRules
}

// parse disasterInfo
func (r *ZoneRouteInfo) ParseDisasterInfo(disasterInfo DisasterInfo) {
	if len(disasterInfo.DisasterItemMap) == 0 {
		return
	}

	disasterItemMap := disasterInfo.DisasterItemMap
	for flowType, disasterItem := range disasterItemMap {
		var statusMap map[bool]map[int32]DisasterStatusEnum
		if flowType == string(FT_PRESS) {
			statusMap = r.markDisasterStatusMap
		} else {
			statusMap = r.disasterStatusMap
		}

		disasterValueMap := disasterItem.DisasterValueMap
		for uidType, disasterValue := range disasterValueMap {
			isElastic := uidType == string(ELASTIC_UID)
			if _, ok := statusMap[isElastic]; !ok {
				statusMap[isElastic] = make(map[int32]DisasterStatusEnum)
			}
			uidStatusMap := disasterValue.UidDisasterStatusMap

			for uidStr, uidStatus := range uidStatusMap {
				uid, _ := strconv.ParseInt(uidStr, 10, 32)
				statusMap[isElastic][int32(uid)] = uidStatus
			}
		}
	}
}

// parse grayInfo
// grayInfo contains all gray zones flow rules
func (r *ZoneRouteInfo) ParseGrayInfo(grayInfo GrayInfo) {
	if len(grayInfo.Groups) == 0 {
		r.grayRouteInfo = GrayRouteInfo{
			ElasticInfo:           initElasticInfo(),
			GrayUidMap:            make(map[int32]string),
			GrayMarkUidMap:        make(map[int32]string),
			GrayElasticUidMap:     make(map[int32]string),
			GrayElasticMarkUidMap: make(map[int32]string),
		}
		return
	}

	crtIdcName := r.curZone.IdcInfo.Name
	var defaultGzone *ZoneInfo
	var gzoneGroup *FlowRuleGroup

	domainMap := r.domainInfoMap
	idcMap := r.idcMap
	grayZoneMap := r.zoneMapHolder.GrayZoneMap
	grayGroupMap := r.groupMapHolder.GrayGroupMap
	grayUidMap := make(map[int32]string)
	grayMarkUidMap := make(map[int32]string)
	grayElasticUidMap := make(map[int32]string)
	grayElasticMarkUidMap := make(map[int32]string)

	for gIndex := range grayInfo.Groups {
		group := &grayInfo.Groups[gIndex]
		groupName := group.GroupName
		if err := group.Init(); err != nil {
			// TODO: I don't know whether continue or not
			// keept it same as mosn.
		}

		// uid assign
		mapCopy(grayUidMap, group.UidMultiRange.CvtToUidMap(groupName))
		mapCopy(grayMarkUidMap, group.MarkUidMultiRange.CvtToUidMap(groupName))
		if group.IsElastic {
			mapCopy(grayElasticUidMap, group.ElasticUidMultiRange.CvtToUidMap(groupName))
			mapCopy(grayElasticMarkUidMap, group.ElasticMarkUidMultiRange.CvtToUidMap(groupName))
		}

		// default gzone
		if group.IsGzoneGroup() {
			gzoneGroup = group
		}

		for zIndex := range group.Zones {
			zone := &group.Zones[zIndex]
			zone.ZoneGroup = group
			zone.AppName = r.appName
			zone.IsElastic = group.IsElastic
			if zone.IdcName != "" {
				idcInfo := idcMap[zone.IdcName]
				zone.IdcInfo = idcInfo
				zone.ZoneDomain = domainMap[zone.ZoneName]
			}
			grayZoneMap[zone.ZoneName] = zone

			if zone.IsGZone() && zone.IdcName == crtIdcName {
				defaultGzone = zone
			}
		}

		grayGroupMap[groupName] = group
	}

	// gray elasticInfo
	elasticInfo := grayInfo.ElasticInfo
	if len(elasticInfo.ElasticValueMap) == 0 {
		elasticInfo = initElasticInfo()
	}

	r.grayRouteInfo = GrayRouteInfo{
		GrayUidMap:            grayUidMap,
		GrayMarkUidMap:        grayMarkUidMap,
		GrayElasticUidMap:     grayElasticUidMap,
		GrayElasticMarkUidMap: grayElasticMarkUidMap,
		ElasticInfo:           elasticInfo,
		GrayDefaultGzone: DefaultZone{
			DefaultZone:        defaultGzone,
			WeightedRoundRobin: NewWeightedRoundRobin(gzoneGroup.Zones, false),
			ZoneMap:            gzoneGroup.ZoneMap,
		},
	}
}

// parse extraInfo
// extraInfo contains all other zone flowRules,
// such as curZone is gray zone, other zone is non-gray zone
func (r *ZoneRouteInfo) ParseExtraInfo(extraInfo ExtraInfo) {
	extraRouteInfo := NewZoneRouteInfo(r.appName)
	if len(extraInfo.Groups) != 0 {
		extraRouteInfo.ParseLdc(extraInfo.Ldc)
		extraRouteInfo.parseExtraGroups(r, extraInfo.Groups)
		extraRouteInfo.ParseDrRule(extraInfo.DrRule)
		extraRouteInfo.ParseElasticInfo(extraInfo.ElasticInfo)
		extraRouteInfo.ParseDisasterInfo(extraInfo.DisasterInfo)
	}

	r.extraRouteInfo = &extraRouteInfo
}

// just use to parse extraRouteInfo.groups
func (r *ZoneRouteInfo) parseExtraGroups(zoneRouteInfo *ZoneRouteInfo, extraGroups []FlowRuleGroup) {
	if len(extraGroups) == 0 {
		return
	}

	var gzoneGroup *FlowRuleGroup
	uidMap := r.uidMap
	markUidMap := r.markUidMap
	elasticUidMap := r.elasticUidMap
	elasticMarkUidMap := r.elasticMarkUidMap
	groupMap := r.groupMapHolder.GroupMap
	zoneMap := r.zoneMapHolder.ZoneMap
	idcMap := r.idcMap

	extGroupMap := zoneRouteInfo.groupMapHolder.ExtraGroupMap
	extZoneMap := zoneRouteInfo.zoneMapHolder.ExtraZoneMap
	extDomainMap := zoneRouteInfo.domainInfoMap
	appName := zoneRouteInfo.appName

	for gIndex := range extraGroups {
		group := &extraGroups[gIndex]
		groupName := group.GroupName
		if err := group.Init(); err != nil {
			//
		}

		// uid assign
		mapCopy(uidMap, group.UidMultiRange.CvtToUidMap(groupName))
		mapCopy(markUidMap, group.MarkUidMultiRange.CvtToUidMap(groupName))
		if group.IsElastic {
			mapCopy(elasticUidMap, group.ElasticUidMultiRange.CvtToUidMap(groupName))
			mapCopy(elasticMarkUidMap, group.ElasticMarkUidMultiRange.CvtToUidMap(groupName))
		}

		for zIndex := range group.Zones {
			zone := &group.Zones[zIndex]
			zoneName := zone.ZoneName
			zone.ZoneGroup = group
			zone.AppName = appName
			zone.IsElastic = group.IsElastic
			if zone.IdcName != "" {
				idcInfo := idcMap[zone.IdcName]
				zone.IdcInfo = idcInfo
				zone.ZoneDomain = extDomainMap[zoneName]
			}

			zoneMap[zoneName] = zone
			extZoneMap[zoneName] = zone
		}

		// defaultGzone
		if group.IsGzoneGroup() {
			gzoneGroup = group
		}

		groupMap[groupName] = group
		extGroupMap[groupName] = group
	}

	r.defaultGZone = DefaultZone{
		WeightedRoundRobin: NewWeightedRoundRobin(gzoneGroup.Zones, false),
		ZoneMap:            gzoneGroup.ZoneMap,
	}
}

// set elastic groups and src group in flowRuleGroup
func (r *ZoneRouteInfo) setElasticGroups(flowRuleGroup *FlowRuleGroup, rdr string) {
	if flowRuleGroup.IsElastic {
		groupMap := r.groupMapHolder.GroupMap

		if srcGroup, ok := groupMap[rdr]; ok {
			elasticGroups := srcGroup.ElasticGroups
			if ok, _ := contains(flowRuleGroup, elasticGroups); !ok {
				srcGroup.ElasticGroups = append(elasticGroups, flowRuleGroup)
			}
			flowRuleGroup.SrcGroup = srcGroup
		}
	}
}

// filter 'BUILDING' zone related info when curZone status is 'RUNNING'
func (r *ZoneRouteInfo) Filter() {
	if r.curZone.Status != string(ZS_RUNNING) {
		return
	}

	// 1. filter zone
	zoneMap := r.zoneMapHolder.ZoneMap
	domainMap := r.domainInfoMap
	vipZoneMap := r.vipZoneMap
	for _, zone := range zoneMap {
		zoneName := zone.ZoneName
		if zone.Status == string(ZS_BUILDING) {
			delete(domainMap, zoneName)
			delete(vipZoneMap, zoneName)
			delete(zoneMap, zoneName)
		}
	}

	// 2. filter groups
	groupMap := r.groupMapHolder.GroupMap
	drRuleMap := r.drRuleMap
	for groupName, group := range groupMap {
		tempZones := make([]ZoneInfo, 0)
		for i, zone := range group.Zones {
			if _, ok := zoneMap[zone.ZoneName]; ok {
				tempZones = append(tempZones, group.Zones[i])
			}
		}
		group.Zones = tempZones

		if len(tempZones) == 0 {
			delete(groupMap, groupName)
			delete(drRuleMap, groupName)
		}
	}

	// 3. filter idc
	idcMap := r.idcMap
	for i := range idcMap {
		zones := idcMap[i].Zones

		tempZones := make([]IdcZoneInfo, 0)
		for j, zone := range idcMap[i].Zones {
			if _, ok := zoneMap[zone.Name]; ok {
				tempZones = append(tempZones, zones[j])
			}
		}
		idcMap[i].Zones = tempZones

		if len(tempZones) == 0 {
			delete(idcMap, i)
		}
	}

	// 4. filter city
	cities := r.cities
	tempCities := make([]CityIdcInfo, 0)
	for i := range cities {
		idcs := cities[i].Idcs
		tempIdcs := make([]IdcInfo, 0)
		for j := range idcs {
			tempZones := make([]IdcZoneInfo, 0)
			for k, zone := range idcs[j].Zones {
				if _, ok := zoneMap[zone.Name]; ok {
					tempZones = append(tempZones, idcs[j].Zones[k])
				}
			}
			idcs[j].Zones = tempZones

			if len(tempZones) != 0 {
				tempIdcs = append(tempIdcs, idcs[j])
			}
		}
		cities[i].Idcs = tempIdcs

		if len(tempIdcs) != 0 {
			tempCities = append(tempCities, cities[i])
		}
	}
	r.cities = tempCities
}

// ############################################
// get functions
// ############################################

func (r *ZoneRouteInfo) GetDefaultGzone(inGray int32, isMark bool) *ZoneInfo {
	var defaultGzone DefaultZone
	switch inGray {
	case ROUTE_IN_GRAY:
		defaultGzone = r.grayRouteInfo.GrayDefaultGzone
	case ROUTE_OUT_GRAY:
		defaultGzone = r.extraRouteInfo.defaultGZone
	default:
		defaultGzone = r.defaultGZone
	}

	return defaultGzone.GetZoneByColor(isMark, r.isGzoneColor, r.curZone.ZoneName, r.zoneColorMap)
}

func (r *ZoneRouteInfo) GetDefaultCzone(isMark bool) *ZoneInfo {
	return r.defaultCZone.GetZoneByColor(isMark, r.isGzoneColor, r.curZone.ZoneName, r.zoneColorMap)
}

func (r *ZoneRouteInfo) GetCurZone() *ZoneInfo {
	return &r.curZone
}

func (r *ZoneRouteInfo) GetAllZoneNames() []string {
	return r.zoneNames
}

func (r *ZoneRouteInfo) GetElasticInfo(inGray int32) ElasticInfo {
	switch inGray {
	case ROUTE_IN_GRAY:
		return r.grayRouteInfo.ElasticInfo
	case ROUTE_OUT_GRAY:
		return r.extraRouteInfo.elasticInfo
	default:
		return r.elasticInfo
	}
}

func (r *ZoneRouteInfo) GetGroupMap(inGray int32) map[string]*FlowRuleGroup {
	return r.groupMapHolder.GetGroupMap(inGray)
}

func (r *ZoneRouteInfo) GetFlowRule() FlowRule {
	if len(r.flowRules) > 0 {
		return r.flowRules[0]
	}
	return FlowRule{}
}

func (r *ZoneRouteInfo) GetUIDMap(inGray int32, isElastic bool, isMark bool) map[int32]string {
	switch inGray {
	case ROUTE_IN_GRAY:
		if isElastic {
			if isMark {
				return r.grayRouteInfo.GrayElasticMarkUidMap
			} else {
				return r.grayRouteInfo.GrayElasticUidMap
			}
		} else {
			if isMark {
				return r.grayRouteInfo.GrayMarkUidMap
			} else {
				return r.grayRouteInfo.GrayUidMap
			}
		}

	case ROUTE_OUT_GRAY:
		if isElastic {
			if isMark {
				return r.extraRouteInfo.elasticMarkUidMap
			} else {
				return r.extraRouteInfo.elasticUidMap
			}
		} else {
			if isMark {
				return r.extraRouteInfo.markUidMap
			} else {
				return r.extraRouteInfo.uidMap
			}
		}

	default:
		if isElastic {
			if isMark {
				return r.elasticMarkUidMap
			} else {
				return r.elasticUidMap
			}
		} else {
			if isMark {
				return r.markUidMap
			} else {
				return r.uidMap
			}
		}
	}
}

func (r *ZoneRouteInfo) GetDefaultUIDMap() map[int32]string {
	return r.defaultUidMap
}

func (r *ZoneRouteInfo) GetZoneByName(zoneName string) *ZoneInfo {
	return r.zoneMapHolder.GetByZoneName(zoneName)
}

func (r *ZoneRouteInfo) GetSrcZone(zoneName string, isMark bool) (ZoneInfo, error) {
	if zone := r.zoneMapHolder.GetByZoneName(zoneName); zone != nil {
		if !zone.IsElastic {
			return ZoneInfo{}, errors.New(zoneName + " is not elastic zone")
		}

		if srcGroup := zone.ZoneGroup.SrcGroup; srcGroup != nil {
			return r.routeByWeight(srcGroup, isMark), nil
		} else {
			return ZoneInfo{}, errors.New(" acquire src group failed," + zoneName)
		}
	} else {
		return ZoneInfo{}, errors.New(zoneName + " is not exist")
	}
}

func (r *ZoneRouteInfo) GetDisasterStatusMap(isMark bool) map[bool]map[int32]DisasterStatusEnum {
	if isMark {
		return r.markDisasterStatusMap
	}
	return r.disasterStatusMap
}

func (r *ZoneRouteInfo) GetExtraZoneRouteInfo() *ZoneRouteInfo {
	return r.extraRouteInfo
}

// ############################################
// route functions
// ############################################

func (r *ZoneRouteInfo) DoRoute(ctx RouteContext) (ZoneInfo, error) {
	zoneInfo := ZoneInfo{}
	if !r.isEnable {
		return zoneInfo, errors.New("enable is false")
	}

	uid := ctx.Uid % r.mod
	groupName := r.GetUIDMap(ctx.InGray, false, ctx.IsMark)[uid]
	if groupName == "" {
		return zoneInfo, fmt.Errorf("can't find uid groups, uid is %d", ctx.Uid)
	}
	group := r.groupMapHolder.GetGroupMap(ctx.InGray)[groupName]
	return r.RouteByColor(group, ctx.IsMark)
}

func (r *ZoneRouteInfo) RouteByColor(group *FlowRuleGroup, isMark bool) (ZoneInfo, error) {
	if group == nil {
		return ZoneInfo{}, errors.New("group is not exist")
	}

	curZoneColor := r.zoneColorMap[r.curZone.ZoneName]
	if curZoneColor != "" {
		for i := range group.Zones {
			zone := group.Zones[i]
			if curZoneColor == r.zoneColorMap[zone.ZoneName] {
				return zone, nil
			}
		}
	}

	if group.isContains(r.curZone.ZoneName) && group.FlowAvg(isMark) {
		return r.curZone, nil
	}
	return r.routeByWeight(group, isMark), nil
}

func (r *ZoneRouteInfo) routeByWeight(group *FlowRuleGroup, isMark bool) ZoneInfo {
	var weightedRoundRobin WeightedRoundRobin
	if isMark {
		weightedRoundRobin = r.markGroupWeightedRoundRobin[group.GroupName]
	} else {
		weightedRoundRobin = r.groupWeightedRoundRobin[group.GroupName]
	}
	zoneName := weightedRoundRobin.GetServerAsPerAlgo()
	return group.ZoneFromGroup(zoneName)
}

func (r *ZoneRouteInfo) GetAllZoneInfo() ([]ZoneInfo, error) {
	result := make([]ZoneInfo, 0)
	for _, val := range r.zoneMapHolder.ZoneMap {
		result = append(result, *val)
	}

	for _, val := range r.zoneMapHolder.GrayZoneMap {
		result = append(result, *val)
	}

	for _, val := range r.zoneMapHolder.ExtraZoneMap {
		result = append(result, *val)
	}
	return result, nil
}

func (r *ZoneRouteInfo) GetCloselyRzone(filterElasticZone, isPress bool) (ZoneInfo, error) {
	var curZone = r.GetCurZone()
	if curZone == nil {
		return ZoneInfo{}, fmt.Errorf("zoneclient: not found cur zone")
	}

	if curZone.IsRZone() && !(filterElasticZone && curZone.IsElastic) {
		var group = curZone.ZoneGroup
		if hasValidUid(group, isPress) {
			return *curZone, nil
		}
	}

	//local idc
	var groupMap = make(map[string]*FlowRuleGroup, 0)
	var checkedZoneMap = make(map[string]string, 0)
	checkedZoneMap[curZone.ZoneName] = curZone.ZoneName
	var idcInfo = curZone.IdcInfo
	if idcInfo == nil {
		return ZoneInfo{}, fmt.Errorf("zoneclient: not found idc info from zone %s", curZone.ZoneName)
	}

	if len(idcInfo.Zones) > 0 {
		for _, zone := range idcInfo.Zones {
			var zoneInfo = r.GetZoneByName(zone.Name)
			if zoneInfo == nil {
				return ZoneInfo{}, fmt.Errorf("zoneclient: not found zone %s", zone.Name)
			}

			var zoneGroup = zoneInfo.ZoneGroup
			if zoneInfo.IsRZone() && checkedZoneMap[zoneInfo.ZoneName] == "" {
				checkedZoneMap[zoneInfo.ZoneName] = zoneInfo.ZoneName

				if !(filterElasticZone && zoneInfo.IsElastic) {

					if hasValidUid(zoneGroup, isPress) {
						groupMap[zoneGroup.GroupName] = zoneGroup
					}
				}
			}
		}
	}

	var otherGroupMap = make(map[string]*FlowRuleGroup, 0)

	//local city
	var curCityInfo = idcInfo.CityIdcInfo
	if curCityInfo == nil {
		return ZoneInfo{}, fmt.Errorf("zoneclient: not found city info from idc %s", idcInfo.Name)
	}
	var allZones = r.zoneMapHolder.ZoneMap
	if len(groupMap) <= 0 {
		for _, zone := range allZones {
			if zone.IsRZone() && checkedZoneMap[zone.ZoneName] == "" && !(filterElasticZone && zone.IsElastic) {
				var idc = zone.IdcInfo
				if idc == nil {
					return ZoneInfo{}, fmt.Errorf("zoneclient: not found idc info from zone %s", zone.ZoneName)
				}

				var zoneCity = idc.CityIdcInfo
				if zoneCity == nil {
					return ZoneInfo{}, fmt.Errorf("zoneclient: not found city info from idc %s", idc.Name)
				}

				var zoneGroup = zone.ZoneGroup
				if hasValidUid(zoneGroup, isPress) {
					if curCityInfo.Name == zoneCity.Name {
						groupMap[zoneGroup.GroupName] = zoneGroup
					} else {
						otherGroupMap[zoneGroup.GroupName] = zoneGroup
					}
				}
			}
		}
	}

	//local city
	var closelyRzone ZoneInfo
	if len(groupMap) > 0 {
		closelyRzone = r.getCityCloselyZone(groupMap, isPress)
	}

	if closelyRzone.ZoneName == "" && len(otherGroupMap) > 0 {
		closelyRzone = r.getCityCloselyZone(otherGroupMap, isPress)
	}

	return closelyRzone, nil
}

func (r *ZoneRouteInfo) getCityCloselyZone(groupMap map[string]*FlowRuleGroup, isMark bool) ZoneInfo {
	var groups = make([]*FlowRuleGroup, 0, len(groupMap))
	for _, zoneGroup := range groupMap {
		groups = append(groups, zoneGroup)
	}

	if len(groups) > 0 {
		var targetGroup = groups[rand.Intn(len(groups))]
		return r.routeByWeight(targetGroup, isMark)
	}

	return ZoneInfo{}
}

// ############################################
// tool functions
// ############################################

//
func hasValidUid(group *FlowRuleGroup, isMark bool) bool {
	if group == nil {
		return false
	}

	var uidMultiRange = group.FindCurUidRange(isMark)
	if len(uidMultiRange.ListUidRange) > 0 {
		for _, uidRange := range uidMultiRange.ListUidRange {
			if uidRange.UidMaxValue >= 0 {
				return true
			}
		}
	}

	return false
}

// init elastic info
func initElasticInfo() ElasticInfo {
	elasticValueMap := make(map[string]ElasticValue)
	elasticValue := ElasticValue{
		Status:        string(ES_NORMAL),
		ElasticValues: make([]string, 0),
	}
	elasticValueMap["ONLINE"] = elasticValue
	elasticValueMap["PRESS"] = elasticValue

	return ElasticInfo{
		ElasticValueMap: elasticValueMap,
	}
}

// copy from src map to dst map
func mapCopy(dst, src interface{}) {
	dv, sv := reflect.ValueOf(dst), reflect.ValueOf(src)

	for _, k := range sv.MapKeys() {
		dv.SetMapIndex(k, sv.MapIndex(k))
	}
}

// check if obj contains target
func contains(obj interface{}, target interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	}

	return false, errors.New("not in array")
}
