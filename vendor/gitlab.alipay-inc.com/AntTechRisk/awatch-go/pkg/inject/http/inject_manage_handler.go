package http

import (
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/executor"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/log"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/rule"
	"strings"
	"sync"
	"time"
)

type InjectManageHandler struct {
}

func (*InjectManageHandler) Handle(req *RequestEntity, resp *ResponseEntity) {
	opt := req.Opt

	if opt == "start" {
		startInject(req, resp)
	}
	if opt == "stop" {
		stopInject(req, resp)
	}
	if opt == "stopAll" {
		stopAllInject(req, resp)
	}
}

var lock sync.Mutex

func startInject(req *RequestEntity, resp *ResponseEntity) {
	lock.Lock()
	defer lock.Unlock()

	// injectId校验
	injectId := req.InjectId
	if injectId == "" {
		resp.Success = false
		resp.ErrorMsg = "injectId is blank"
		return
	}

	// 查找注入规则
	injectRule, exists := rule.LoadRuleById(injectId)
	if !exists {
		resp.Success = false
		resp.ErrorMsg = "injectRule not exist, injectId=" + injectId
		return
	}
	// 幂等控制
	if injectRule.IsActive {
		resp.Success = true
		resp.Msg = "inject started, injectId=" + injectId
		return
	}

	// 生效注入规则
	injectRule.IsActive = true
	injectRule.ActiveTime = time.Now().Unix() // 单位秒

	resp.Success = true
	resp.Msg = "inject started, injectId=" + injectId

	log.AgentLogger.Info("start inject success,injectId=" + injectId)

	// 全局类型的故障直接开始注入
	if rule.IsInstantlyInjectType(injectRule.InjectType) {
		executor.InjectInstantly(injectRule)
	}
}

func stopInject(req *RequestEntity, resp *ResponseEntity) {
	lock.Lock()
	defer lock.Unlock()

	// injectId校验
	injectId := req.InjectId
	if injectId == "" {
		resp.Success = false
		resp.ErrorMsg = "injectId is blank"
		return
	}

	// 查找注入规则
	injectRule, exists := rule.LoadRuleById(injectId)
	if !exists {
		resp.Success = false
		resp.ErrorMsg = "injectRule not exist, injectId=" + injectId
		return
	}

	// 失效注入规则
	injectRule.IsActive = false
	resp.Success = true
	resp.Msg = "inject stopped, injectId=" + injectId

	log.AgentLogger.Info("stop inject success,injectId=" + injectId)
}

func stopAllInject(req *RequestEntity, resp *ResponseEntity) {
	lock.Lock()
	defer lock.Unlock()

	allInjectId := rule.LoadAllRuleId()
	injectIdArray := strings.Split(allInjectId, ",")
	for _, injectId := range injectIdArray {
		injectRule, exists := rule.LoadRuleById(injectId)
		if exists {
			injectRule.IsActive = false
		}
	}

	resp.Success = true
	resp.Msg = "all inject stopped, injectId=" + allInjectId

	log.AgentLogger.Info("stop all inject success,injectId=" + allInjectId)
}
