package fault

import (
	"errors"
	"fmt"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/log"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/rule"
)

type ErrorInjector struct {
}

func (*ErrorInjector) Inject(injectContext map[string]string, injectRule *rule.InjectRule) *InjectResult {
	log.DefaultLogger.Info(fmt.Sprintf("error inject start,injectId=%v,traceId=%v", injectRule.InjectId, injectContext[rule.RpcKeyTraceId]))

	injectResult := new(InjectResult)

	errorMsg := injectRule.ExceptionMsg
	if errorMsg == "" {
		errorMsg = "injected error"
	}

	injectedError := errors.New(errorMsg)

	// 产生错误，并阻断调用
	injectResult.Success = true
	injectResult.InvokeInterrupt = true
	injectResult.InjectedError = injectedError

	return injectResult
}
