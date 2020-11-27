package http

import (
	"encoding/json"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/log"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/rule"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/util"
)

type InjectRuleHandler struct {
}

func (*InjectRuleHandler) Handle(req *RequestEntity, resp *ResponseEntity) {
	opt := req.Opt

	if opt == "insert" {
		insertRule(req, resp)
	}
	if opt == "remove" {
		removeRule(req, resp)
	}
	if opt == "clear" {
		clearRules(req, resp)
	}
	if opt == "get" {
		getRule(req, resp)
	}
	if opt == "getAll" {
		getAll(req, resp)
	}
}

func insertRule(req *RequestEntity, resp *ResponseEntity) {
	// injectId校验
	injectId := req.InjectId
	if injectId == "" {
		resp.Success = false
		resp.ErrorMsg = "injectId is blank"
		return
	}

	// 反序列化为InjectRule
	var injectRule rule.InjectRule
	err := json.Unmarshal([]byte(req.RuleJson), &injectRule)
	if err != nil {
		resp.Success = false
		resp.ErrorMsg = "insert rule, json unmarshal error, injectId=" + injectId
		return
	}

	// 保存注入规则
	rule.StoreRule(injectId, &injectRule)
	resp.Success = true
	resp.Msg = "insert rule success, injectId=" + injectId

	log.AgentLogger.Info("insert inject rule success,rule=" + util.ToJsonString(injectRule))
}

func removeRule(req *RequestEntity, resp *ResponseEntity) {
	// injectId校验
	injectId := req.InjectId
	if injectId == "" {
		resp.Success = false
		resp.ErrorMsg = "injectId is blank"
		return
	}

	// 删除注入规则
	rule.RemoveRule(injectId)
	resp.Success = true
	resp.Msg = "remove rule success, injectId=" + injectId

	log.AgentLogger.Info("remove inject rule success,injectId=" + injectId)
}

func clearRules(req *RequestEntity, resp *ResponseEntity) {
	// 删除所有注入规则
	rule.ClearRules()
	resp.Success = true
	resp.Msg = "clear rules success"

	log.AgentLogger.Info("clear rules success")
}

func getRule(req *RequestEntity, resp *ResponseEntity) {
	// injectId校验
	injectId := req.InjectId
	if injectId == "" {
		resp.Success = false
		resp.ErrorMsg = "injectId is blank"
		return
	}

	// 获取injectId对应的注入规则
	injectRule, exists := rule.LoadRuleById(injectId)
	if !exists {
		resp.Success = false
		resp.ErrorMsg = "injectRule not exist, injectId=" + injectId
		return
	}

	// 返回序列化的注入规则
	result, error := json.Marshal(injectRule)
	if error != nil {
		resp.Success = false
		resp.ErrorMsg = "get rule, json marshal error, injectId=" + injectId
		return
	}

	resp.Success = true
	resp.Msg = "get rule success:" + string(result)

	log.AgentLogger.Info("get inject rule success,rule=" + string(result))
}

func getAll(req *RequestEntity, resp *ResponseEntity) {
	// 返回当前所有规则的injectId
	allInjectId := rule.LoadAllRuleId()
	resp.Success = true
	resp.Msg = "all injectId:" + allInjectId

	log.AgentLogger.Info("get all inject rule success,id=" + allInjectId)
}
