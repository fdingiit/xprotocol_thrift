package executor

import (
	"fmt"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/common"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/fault"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/log"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/rule"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/util"
)

// 立即注入，指的是故障规则激活后，故障现象即刻产生，不依赖服务调用
// CPU，内存等类型的注入，属于这一类
// awatch在激活故障时自动调用，mosn不要调用这个函数
func InjectInstantly(injectRule *rule.InjectRule) *fault.InjectResult {
	injectResult := new(fault.InjectResult)

	// 注入总控开关
	if common.InjectOff {
		injectResult.Success = false
		injectResult.ErrorMsg = "inject switch off"
		return injectResult
	}

	return inject0(make(map[string]string), injectRule)
}

// 注入，故障行为产生依赖服务调用，mosn调用这个函数
// mosn需要在injectContext中传递如下信息(key):
// =============================
// 1、服务接口名称: interfaceName
// 2、服务方法名称: methodName
// 3、服务UniqueID: uniqueId(可选)
// 4、调用方: callerAppName(可选)
// 5、被调用方: targetAppName(可选)
// 6、流量类型标识: mark
// 7、链路ID: traceId
// =============================
// 根据返回结果InjectResult判定是否需要阻断mosn中的当前服务调用
// 当且仅当 (InjectResult.Success == true && InjectResult.InvokeInterrupt == true) 时，阻断当前服务调用
func Inject(injectContext map[string]string) *fault.InjectResult {
	// 入口处打印debug日志
	if common.DebugLog {
		entryLog := fmt.Sprintf("[entry]inject context:%v", injectContext)
		log.DefaultLogger.Info(entryLog)
	}

	injectResult := new(fault.InjectResult)

	// 注入总控开关
	if common.InjectOff {
		injectResult.Success = false
		injectResult.ErrorMsg = "inject switch off"
		return injectResult
	}

	// 匹配注入规则
	injectRule, exists := rule.LoadFirstActiveRuleBySig(injectContext)
	if !exists {
		injectResult.Success = false
		injectResult.ErrorMsg = "inject rule not exists"
		return injectResult
	}

	// 强制注入类型必须为RPC注入
	if injectRule.InjectType != rule.InjectTypeRpc {
		injectResult.Success = false
		injectResult.ErrorMsg = "inject type invalid, must be rpc"
		return injectResult
	}

	// 以注入流量数作为控制条件时，必须加锁保证注入流量数读写的原子性
	if injectRule.InjectCondition.MaxCount > 0 {
		injectRule.InjectLock.Lock()
		defer injectRule.InjectLock.Unlock()

		return inject0(injectContext, injectRule)
	} else {
		return inject0(injectContext, injectRule)
	}
}

func doInject(injectContext map[string]string, injectRule *rule.InjectRule) *fault.InjectResult {
	injectResult := new(fault.InjectResult)

	// 选择注入器进行注入
	injector, exists := fault.GetInjector(injectRule.FaultAction)
	if !exists {
		injectResult.Success = false
		injectResult.ErrorMsg = "injector not exist"
		return injectResult
	}

	return injector.Inject(injectContext, injectRule)
}

func postInject(injectContext map[string]string, injectRule *rule.InjectRule) {
	// 命中流量数增加
	injectRule.InjectedCount++

	injectLog := fmt.Sprintf("%v,%v,%v,%v", injectRule.InjectId, injectRule.InjectType, injectRule.FaultAction, injectContext[rule.RpcKeyTraceId])
	log.InjectLogger.Info(injectLog)
}

func inject0(injectContext map[string]string, injectRule *rule.InjectRule) *fault.InjectResult {
	log.DefaultLogger.Info("inject start,injectId=" + injectRule.InjectId + ",injectType=" + injectRule.InjectType +
		",context=" + util.ToJsonString(injectContext))

	injectResult := new(fault.InjectResult)
	// 测试模式
	if common.TestMode {
		injectResult.Success = true
		injectResult.ErrorMsg = "test mode"
		return injectResult
	}

	// 注入条件判断
	if matched, msg := rule.IsInjectConditionMatched(injectContext, injectRule); !matched {
		log.DefaultLogger.Info("inject condition not matched,injectId=" + injectRule.InjectId + ",reason=" + msg + ",traceId=" + injectContext[rule.RpcKeyTraceId])

		injectResult.Success = false
		injectResult.ErrorMsg = msg
		return injectResult
	}

	// 注入
	injectResult = doInject(injectContext, injectRule)
	if !injectResult.Success {
		return injectResult
	}

	// 注入后，打印日志等
	postInject(injectContext, injectRule)

	injectResult.Success = true
	return injectResult
}
