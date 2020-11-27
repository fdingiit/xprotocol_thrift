package rule

import "time"

type InjectCondition struct {
	// 最大注入次数
	MaxCount int `json:"maxCount"`
	// 注入持续时间，单位秒
	DurationTime int64 `json:"durationTime"`
	// 注入流量类型，包括正常流量，压测流量等
	FlowType string `json:"flowType"`
}

const (
	FlowTypeAll    = "all_flow"
	FlowTypeNormal = "normal_flow"
	FlowTypeShadow = "shadow_flow"
)

func IsInjectConditionMatched(injectContext map[string]string, injectRule *InjectRule) (bool, string) {
	// TODO 每增加一个注入条件，都要检查这里是不是要进行新增条件的判断
	if IsInstantlyInjectType(injectRule.InjectType) {
		return true, ""
	}

	injectCondition := injectRule.InjectCondition

	// 次数和持续时间至少存在其中一个
	if injectCondition.MaxCount <= 0 && injectCondition.DurationTime <= 0 {
		return false, "max count and duration time not exist"
	}

	// 1、次数
	if !IsCountConditionMatched(injectRule) {
		return false, "count condition not matched"
	}

	// 2、持续时间
	if !IsDurationTimeConditionMatched(injectRule) {
		return false, "duration time condition not matched"
	}

	// 3、流量类型
	if !IsFlowTypeConditionMatched(injectContext, injectRule) {
		return false, "flow type condition not matched"
	}

	return true, ""
}

func IsCountConditionMatched(injectRule *InjectRule) bool {
	injectCondition := injectRule.InjectCondition

	if injectCondition.MaxCount <= 0 {
		return true
	}

	return injectRule.InjectedCount < injectCondition.MaxCount
}

func IsDurationTimeConditionMatched(injectRule *InjectRule) bool {
	injectCondition := injectRule.InjectCondition

	if injectCondition.DurationTime <= 0 {
		return true
	}

	now := time.Now().Unix() // 当前时间戳，单位秒
	return now < injectRule.ActiveTime+injectCondition.DurationTime
}

func IsFlowTypeConditionMatched(injectContext map[string]string, injectRule *InjectRule) bool {
	injectCondition := injectRule.InjectCondition
	if injectCondition.FlowType == "" {
		return false
	}
	if injectCondition.FlowType == FlowTypeAll {
		return true
	}

	// 从注入上下文中获取压测标识
	loadTestMark := injectContext[RpcKeyMark]
	trafficType := FlowTypeAll
	if loadTestMark == "T" || loadTestMark == "t" {
		trafficType = FlowTypeShadow
	}
	if loadTestMark == "F" || loadTestMark == "f" {
		trafficType = FlowTypeNormal
	}
	// TODO 可能需要支持镜像流量

	return injectCondition.FlowType == trafficType
}
