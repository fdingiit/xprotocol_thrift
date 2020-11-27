package model

import (
	"errors"
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// parse json string to zoneRouteInfo
func ParseZoneRouteInfo(appName, localZone, zoneInfoStr string) (*ZoneRouteInfo, error) {
	zoneRouteInfo := NewZoneRouteInfo(appName)

	rules := FlowRules{}
	if err := json.UnmarshalFromString(zoneInfoStr, &rules); err != nil {
		return &zoneRouteInfo, fmt.Errorf("failed to json unmarshal: %v", err)
	}
	if len(rules) == 0 {
		return &zoneRouteInfo, errors.New("zone route rule is empty, zoneClient init failed")
	}

	flowRule := rules[0]
	zoneRouteInfo.SetFlowRules(rules)
	zoneRouteInfo.SetVersion(flowRule.Version)
	zoneRouteInfo.SetIsEnable(flowRule.Enable)
	zoneRouteInfo.SetMod(flowRule.Rule)

	// parse domain
	zoneRouteInfo.ParseDomain(flowRule.Domain)
	// parse ldc
	zoneRouteInfo.ParseLdc(flowRule.Ldc)
	// parse flowRuleGroups
	if err := zoneRouteInfo.ParseFlowRuleGroup(localZone, flowRule.FlowRuleGroupes); err != nil {
		return nil, err
	}
	// parse drRule
	zoneRouteInfo.ParseDrRule(flowRule.DrRule)
	// parse elasticInfo
	zoneRouteInfo.ParseElasticInfo(flowRule.ElasticInfo)
	// set defaultCzone
	if err := zoneRouteInfo.SetDefaultCZone(); err != nil {
		return nil, err
	}
	// set defaultGzone
	if err := zoneRouteInfo.SetDefaultGZone(); err != nil {
		return nil, err
	}
	// parse disasterInfo
	zoneRouteInfo.ParseDisasterInfo(flowRule.DisasterInfo)
	// parse grayInfo
	zoneRouteInfo.ParseGrayInfo(flowRule.GrayInfo)
	// parse extraInfo
	zoneRouteInfo.ParseExtraInfo(flowRule.ExtraInfo)

	// filter building zones
	zoneRouteInfo.Filter()
	zoneRouteInfo.InitWeightedRoundRobin()

	return &zoneRouteInfo, nil
}

// parse json string to elasticRuleInfo
func ParseElasticRule(elasticInfoStr string) (*ElasticRuleInfo, error) {
	var ruleInfo ElasticRuleInfo
	err := json.Unmarshal([]byte(elasticInfoStr), &ruleInfo)
	if err != nil {
		return &ruleInfo, err
	}

	// parse elastic biz to tree
	elasticRuleTree, err := parseToElasticTree(ruleInfo.ElasticRuleMap)
	if err != nil {
		return &ruleInfo, err
	}

	ruleInfo.ElasticRuleTree = elasticRuleTree

	return &ruleInfo, nil
}

// parse elastic rule to tree
func parseToElasticTree(ruleMap map[string]interface{}) (ElasticRuleTree, error) {
	ruleTree := ElasticRuleTree{
		IsUseDefaultEid: false,
		Status:          INVALID,
		ElasticRuleMap:  make(map[string]ElasticRuleTree),
	}

	if len(ruleMap) == 0 {
		return ruleTree, nil
	}

	var subTree ElasticRuleTree
	var err error = nil
	for k, v := range ruleMap {
		if _, ok := v.(map[string]interface{}); ok {
			subTree, err = parseToElasticTree(v.(map[string]interface{}))
		} else {
			subTree, err = parseEmptyTree(v.(string))
		}

		if err == nil {
			ruleTree.Put(k, subTree)
		} else {
			break
		}
	}

	return ruleTree, err
}

// get empty biz tree
func parseEmptyTree(value string) (ElasticRuleTree, error) {
	elasticRuleTree := ElasticRuleTree{
		IsUseDefaultEid: false,
		Status:          INVALID,
	}
	if len(value) != 2 {
		return elasticRuleTree, errors.New("status format is not correct")
	}

	elasticRuleTree.Status = ElasticSubRuleStatus(value[0:1])
	elasticRuleTree.IsUseDefaultEid = value[1:2] != "0"
	return elasticRuleTree, nil
}
