package fault

import (
	"fmt"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/common"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/log"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/rule"
	"regexp"
	"runtime"
	"runtime/debug"
)

var cpuBurnRegex string

func init() {
	for i := 0; i < 20; i++ {
		cpuBurnRegex += "zxcvbnmasdf123#456"
	}
}

type CpuInjector struct {
}

func (*CpuInjector) Inject(injectContext map[string]string, injectRule *rule.InjectRule) *InjectResult {
	log.DefaultLogger.Info("cpu inject start,injectId=" + injectRule.InjectId)

	injectResult := new(InjectResult)

	// 此类故障必须有持续时间作为注入条件
	injectDurationTime := injectRule.InjectCondition.DurationTime
	if injectDurationTime <= 0 {
		injectResult.Success = false
		injectResult.ErrorMsg = "duration time condition not exist"
		return injectResult
	}

	// 设置注入的CPU核数
	cpuBurnCoreCount := injectRule.CpuBurnSize
	if cpuBurnCoreCount <= 0 {
		cpuBurnCoreCount = runtime.NumCPU()
	}

	log.DefaultLogger.Info(fmt.Sprintf("burn cpu with %d goroutine,injectId=%v", cpuBurnCoreCount, injectRule.InjectId))
	// 开启goroutine消耗cpu
	for i := 0; i < cpuBurnCoreCount; i++ {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.DefaultLogger.Error("goroutine recovered, stack:" + string(debug.Stack()))
				}
			}()

			// 无限循环解析正则表达式
			for {
				// 注入开关已关闭
				if common.InjectOff {
					break
				}

				// 注入规则已失效，退出无限循环
				if !injectRule.IsActive {
					break
				}

				// 注入持续时间已超过，退出无限循环
				if !rule.IsDurationTimeConditionMatched(injectRule) {
					break
				}

				// 解析正则表达式
				regexp.MatchString("^([a-zA-Z]+)+$", cpuBurnRegex)

				// 手动调度其他goroutine的执行
				runtime.Gosched()
			}
		}()
	}

	injectResult.Success = true
	return injectResult
}
