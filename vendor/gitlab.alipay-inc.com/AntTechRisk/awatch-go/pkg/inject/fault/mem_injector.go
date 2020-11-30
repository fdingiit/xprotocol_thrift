package fault

import (
	"container/list"
	"fmt"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/common"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/log"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/rule"
	"runtime"
	"runtime/debug"
)

type Foo struct {
	b [1024 * 1024]byte
}

type MemoryInjector struct {
}

func (*MemoryInjector) Inject(injectContext map[string]string, injectRule *rule.InjectRule) *InjectResult {
	log.DefaultLogger.Info("memory inject start,injectId=" + injectRule.InjectId)

	injectResult := new(InjectResult)

	// 此类故障必须有持续时间作为注入条件
	injectDurationTime := injectRule.InjectCondition.DurationTime
	if injectDurationTime <= 0 {
		injectResult.Success = false
		injectResult.ErrorMsg = "duration time condition not exist"
		return injectResult
	}

	count := runtime.NumCPU()

	log.DefaultLogger.Info(fmt.Sprintf("occupy memory with %d goroutine,injectId=%v", count, injectRule.InjectId))
	for i := 0; i < count; i++ {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.DefaultLogger.Error("goroutine recovered, stack:" + string(debug.Stack()))
				}
			}()

			fooList := list.New()
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

				// 制造内存占用现象
				var bytes [1024 * 1024]byte
				foo := Foo{bytes}
				fooList.PushBack(foo)

				// 手动调度其他goroutine的执行
				runtime.Gosched()
			}
		}()
	}

	injectResult.Success = true
	return injectResult
}
