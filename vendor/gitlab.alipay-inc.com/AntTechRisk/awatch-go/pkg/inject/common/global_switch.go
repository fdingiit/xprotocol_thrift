package common

import (
	"encoding/json"
)

// debug日志开关
var DebugLog = false
// 注入总控开关
var InjectOff = true
// 测试模式开关(测试模式中故障行为不执行，便于进行单元测试)
var TestMode = false

func BuildGlobalSwitchStr() string {
	m := make(map[string]interface{})

	m["debugLog"] = DebugLog
	m["injectOff"] = InjectOff
	m["testMode"] = TestMode

	json, _ := json.Marshal(m)
	return string(json)
}
