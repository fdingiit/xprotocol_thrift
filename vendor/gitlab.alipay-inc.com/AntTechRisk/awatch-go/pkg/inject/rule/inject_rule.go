package rule

import (
	"sync"
)

// 枚举: 注入类型
const (
	InjectTypeRpc     = "rpc"
	InjectTypeCpu     = "cpu"
	InjectTypeMemory  = "mem"
	InjectTypeDisk    = "disk"
	InjectTypeProcess = "proc"
)

// 枚举: 故障行为
const (
	FaultActionDelay       = "delay"
	FaultActionException   = "exception"
	FaultActionSkip        = "skip"
	FaultActionBurnCpu     = "burn_cpu"
	FaultActionMemLeak     = "mem_leak"
	FaultActionFloodDisk   = "flood_disk"
	FaultActionKillProcess = "kill_proc"
)

type InjectRule struct {
	// 注入全局唯一ID
	InjectId string `json:"injectId"`
	// 注入规则owner
	Owner string `json:"owner"`
	// 注入类型
	InjectType string `json:"injectType"`
	// 故障行为
	FaultAction string `json:"faultAction"`
	// 签名信息,json
	Signature string `json:"signature"`

	// 延时类故障的延时，单位秒
	DelayTime int `json:"delayTime"`
	// 异常类故障的异常类型
	ExceptionType string `json:"exceptionType"`
	// 异常类故障的异常信息
	ExceptionMsg string `json:"exceptionMsg"`
	// cpu故障的燃烧核数
	CpuBurnSize int `json:"cpuBurnSize"`

	// 注入条件
	InjectCondition *InjectCondition `json:"injectCondition"`

	// 是否生效
	IsActive bool `json:"-"`
	// 生效时间戳
	ActiveTime int64 `json:"-"`
	// 已注入次数
	InjectedCount int `json:"-"`

	// 控制并发注入的锁
	InjectLock sync.Mutex `json:"-"`
}

var instantlyInjectTypes = []string{InjectTypeCpu, InjectTypeMemory, InjectTypeDisk, InjectTypeProcess}

func IsInstantlyInjectType(injectType string) bool {
	for _, val := range instantlyInjectTypes {
		if val == injectType {
			return true
		}
	}

	return false
}
