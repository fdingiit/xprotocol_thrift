package model

import (
	"strings"
)

const ALL = "*"

type ElasticRuleInfo struct {
	Version         int64                  `json:"version"`    // elastic rule version
	ElasticRuleMap  map[string]interface{} `json:"elasticBiz"` // elastic rule tree
	ElasticRuleTree ElasticRuleTree        `json:"-"`
}

type BizElastic struct {
	IsElastic       bool
	IsUseDefaultEid bool
}

type ElasticRuleTree struct {
	IsUseDefaultEid bool
	Status          ElasticSubRuleStatus
	ElasticRuleMap  map[string]ElasticRuleTree
}

// determine whether contains elastic rules
func (ert *ElasticRuleTree) Contains(conditions []string, isMark bool) BizElastic {
	ruleMap := ert.ElasticRuleMap
	isUseDefaultEid := false
	var status ElasticSubRuleStatus

	for _, condition := range conditions {

		// contains all ?
		if all, ok := ruleMap[ALL]; ok {
			isUseDefaultEid = all.IsUseDefaultEid
			status = all.Status
			break
		}

		key := strings.TrimSpace(condition)
		if elasticInfo, ok := ruleMap[key]; ok {
			if len(elasticInfo.ElasticRuleMap) > 0 {
				// elasticInfo is not empty
				ruleMap = elasticInfo.ElasticRuleMap
			} else {
				// elasticInfo is empty
				ruleMap = make(map[string]ElasticRuleTree)
				isUseDefaultEid = elasticInfo.IsUseDefaultEid
				status = elasticInfo.Status
			}
		} else {
			return BizElastic{}
		}
	}

	return BizElastic{
		IsElastic:       (isMark && status == PRESS) || status == VALID,
		IsUseDefaultEid: isUseDefaultEid,
	}
}

func (ert *ElasticRuleTree) Put(key string, elasticInfo ElasticRuleTree) {
	ert.ElasticRuleMap[key] = elasticInfo
}
