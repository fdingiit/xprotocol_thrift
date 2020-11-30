package rule

import (
	"encoding/json"
	"strings"
	"sync"
)

var ruleRepo = make(map[string]*InjectRule)
var rwLock = new(sync.RWMutex)

func StoreRule(injectId string, rule *InjectRule) {
	rwLock.Lock()
	defer rwLock.Unlock()

	ruleRepo[injectId] = rule
}

func RemoveRule(injectId string) {
	rwLock.Lock()
	defer rwLock.Unlock()

	delete(ruleRepo, injectId)
}

func LoadRuleById(injectId string) (*InjectRule, bool) {
	rwLock.RLock()
	defer rwLock.RUnlock()

	if rule, exists := ruleRepo[injectId]; exists {
		return rule, true
	}
	return nil, false
}

func LoadFirstActiveRuleBySig(reqSigMap map[string]string) (*InjectRule, bool) {
	rwLock.RLock()
	defer rwLock.RUnlock()

	for injectId := range ruleRepo {
		if injectRule, exists := LoadRuleById(injectId); exists {
			if injectRule.IsActive && isRuleMatched(reqSigMap, injectRule) {
				return injectRule, true
			}
		}
	}

	return nil, false
}

func LoadAllRuleId() string {
	rwLock.RLock()
	defer rwLock.RUnlock()

	var ruleIdBuilder strings.Builder
	for injectId := range ruleRepo {
		ruleIdBuilder.WriteString(injectId)
		ruleIdBuilder.WriteString(",")
	}

	return ruleIdBuilder.String()
}

func ClearRules() {
	rwLock.Lock()
	defer rwLock.Unlock()

	for injectId := range ruleRepo {
		delete(ruleRepo, injectId)
	}
}

func isRuleMatched(reqSigMap map[string]string, injectRule *InjectRule) bool {
	injectType := injectRule.InjectType

	// 不同的注入类型，不同的签名匹配逻辑
	if injectType == InjectTypeRpc {
		return isRpcSigMatched(reqSigMap, injectRule.Signature)
	}

	// TODO 其他注入类型，mosn中只有RPC，后续扩展考虑接口方式

	return false
}

func isRpcSigMatched(reqSigMap map[string]string, ruleSignature string) bool {
	var ruleSigJson RpcSignature
	err := json.Unmarshal([]byte(ruleSignature), &ruleSigJson)
	if err != nil {
		return false
	}

	// 比较签名
	if ruleSigJson.InterfaceName != reqSigMap[RpcSigKeyInterface] {
		return false
	}
	if ruleSigJson.MethodName != reqSigMap[RpcSigKeyMethod] {
		return false
	}
	if ruleSigJson.UniqueId != "" && ruleSigJson.UniqueId != reqSigMap[RpcSigKeyUniqueId] {
		return false
	}
	if ruleSigJson.CallerAppName != "" && ruleSigJson.CallerAppName != reqSigMap[RpcSigKeyCaller] {
		return false
	}
	if ruleSigJson.TargetAppName != "" && ruleSigJson.TargetAppName != reqSigMap[RpcSigKeyTarget] {
		return false
	}

	return true
}
