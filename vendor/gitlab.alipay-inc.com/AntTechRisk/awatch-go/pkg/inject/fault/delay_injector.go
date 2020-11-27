package fault

import (
	"fmt"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/log"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/rule"
	"time"
)

// 接口实现
type DelayInjector struct {
}

func (*DelayInjector) Inject(injectContext map[string]string, injectRule *rule.InjectRule) *InjectResult {
	log.DefaultLogger.Info(fmt.Sprintf("delay inject start,injectId=%v,traceId=%v", injectRule.InjectId, injectContext[rule.RpcKeyTraceId]))

	injectResult := new(InjectResult)

	delayTime := injectRule.DelayTime
	if delayTime <= 0 {
		injectResult.Success = false
		injectResult.ErrorMsg = "delay time must be greater than 0"
		return injectResult
	}

	// 延时单位是秒
	time.Sleep(time.Duration(delayTime) * time.Second)

	injectResult.Success = true
	return injectResult
}
