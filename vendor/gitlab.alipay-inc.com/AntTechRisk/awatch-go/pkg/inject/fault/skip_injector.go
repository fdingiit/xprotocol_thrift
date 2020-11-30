package fault

import (
	"fmt"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/log"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/rule"
)

type SkipInjector struct {
}

func (*SkipInjector) Inject(injectContext map[string]string, injectRule *rule.InjectRule) *InjectResult {
	log.DefaultLogger.Info(fmt.Sprintf("skip inject start,injectId=%v,traceId=%v", injectRule.InjectId, injectContext[rule.RpcKeyTraceId]))

	injectResult := new(InjectResult)

	// 直接阻断调用
	injectResult.Success = true
	injectResult.InvokeInterrupt = true

	return injectResult
}
