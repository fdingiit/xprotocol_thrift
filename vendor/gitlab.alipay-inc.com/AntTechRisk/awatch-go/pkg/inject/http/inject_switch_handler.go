package http

import (
	"encoding/json"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/common"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/log"
)

type InjectSwitchHandler struct {
}

func (*InjectSwitchHandler) Handle(req *RequestEntity, resp *ResponseEntity) {
	opt := req.Opt

	if opt == "set" {
		setSwitch(req, resp)
	}
	if opt == "get" {
		getSwitch(resp)
	}
}

func setSwitch(req *RequestEntity, resp *ResponseEntity) {
	jsonMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(req.RuleJson), &jsonMap)
	if err != nil {
		resp.Success = false
		resp.ErrorMsg = "json unmarshal error"
		return
	}

	// debug日志开关
	if val, exists := jsonMap["debugLog"]; exists {
		common.DebugLog = val.(bool)
	}
	// 全局注入控制开关
	if val, exists := jsonMap["injectOff"]; exists {
		common.InjectOff = val.(bool)
	}

	resp.Success = true
	resp.Msg = "success, current=" + common.BuildGlobalSwitchStr()

	log.AgentLogger.Info("set switch success,current=" + common.BuildGlobalSwitchStr())
}

func getSwitch(resp *ResponseEntity) {
	resp.Success = true
	resp.Msg = "current switch=" + common.BuildGlobalSwitchStr()

	log.AgentLogger.Info("get switch success,result=" + common.BuildGlobalSwitchStr())
}
