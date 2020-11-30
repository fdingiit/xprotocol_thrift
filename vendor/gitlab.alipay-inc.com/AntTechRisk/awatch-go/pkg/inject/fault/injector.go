package fault

import (
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/rule"
)

type InjectResult struct {
	Success         bool   // 注入是否成功
	ErrorMsg        string // 注入失败的错误信息
	InvokeInterrupt bool   // 是否需要中断调用，对于错误和跳过执行的故障行为，这个值为true
	InjectedError   error  // 错误类的故障行为制造的错误
}

// 注入接口
type Injector interface {
	Inject(injectContext map[string]string, injectRule *rule.InjectRule) *InjectResult
}

// 注册injector
var injectorRepo = make(map[string]Injector)

func init() {
	injectorRepo[rule.FaultActionDelay] = new(DelayInjector)
	injectorRepo[rule.FaultActionException] = new(ErrorInjector)
	injectorRepo[rule.FaultActionBurnCpu] = new(CpuInjector)
	injectorRepo[rule.FaultActionMemLeak] = new(MemoryInjector)
	injectorRepo[rule.FaultActionFloodDisk] = new(DiskInjector)
	injectorRepo[rule.FaultActionSkip] = new(SkipInjector)
}

func GetInjector(faultAction string) (Injector, bool) {
	injector, exists := injectorRepo[faultAction]
	if !exists {
		return nil, false
	}

	return injector, true
}
