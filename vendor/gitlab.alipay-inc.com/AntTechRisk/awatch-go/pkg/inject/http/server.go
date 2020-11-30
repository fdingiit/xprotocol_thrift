package http

import (
	"fmt"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/log"
	"net/http"
	"net/url"
	"runtime/debug"
)

var startMark bool

var httpServer *http.Server

func init() {
	// panic: http: multiple registrations for /awatch-go
	// 只能注册一次
	http.HandleFunc("/awatch-go", HandlerFunc)
}

// 显式初始化http server, 需要mosn调用
func InitAwatchHttp() {
	if startMark == true {
		return
	}

	httpServer = &http.Server{Addr: ":18800"}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.DefaultLogger.Error("goroutine recovered, stack:" + string(debug.Stack()))
			}
		}()

		startMark = true

		// 这里会阻塞，监听http连接请求
		err := httpServer.ListenAndServe()
		if err != nil {
			// 该server被shutdown后，ListenAndServe报错返回
			errMsg := fmt.Sprint(err)
			log.DefaultLogger.Error("http server listen error:" + errMsg)
		}
	}()
}

func DestroyAwatchHttp() {
	if httpServer != nil {
		err := httpServer.Shutdown(nil)
		if err != nil {
			errMsg := fmt.Sprint(err)
			log.DefaultLogger.Error("http server shutdown error:" + errMsg)
		} else {
			startMark = false
			httpServer = nil
			log.DefaultLogger.Info("http server shutdown success")
		}
	}
}

func HandlerFunc(w http.ResponseWriter, r *http.Request) {
	// 解析get请求参数
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		httpRespond(w, "parse request query error", http.StatusInternalServerError)
		return
	}

	reqParamMap := make(map[string]string)
	for k, v := range queryForm {
		if len(v) > 0 {
			reqParamMap[k] = v[0]
		}
	}

	// 组装RequestEntity
	var req = new(RequestEntity)
	if val, exists := reqParamMap["topic"]; exists {
		req.Topic = val
	}
	if val, exists := reqParamMap["opt"]; exists {
		req.Opt = val
	}
	if val, exists := reqParamMap["ruleJson"]; exists {
		req.RuleJson = val
	}
	if val, exists := reqParamMap["injectId"]; exists {
		req.InjectId = val
	}
	if val, exists := reqParamMap["appName"]; exists {
		req.AppName = val
	}
	if val, exists := reqParamMap["operator"]; exists {
		req.Operator = val
	}

	// 获取处理该topic的handler
	awatchHandler, exists := HandlerRepo[req.Topic]
	if !exists {
		httpRespond(w, "invalid topic:"+req.Topic, http.StatusInternalServerError)
		return
	}

	// 处理请求
	var resp = new(ResponseEntity)
	awatchHandler.Handle(req, resp)

	// 返回处理结果
	if !resp.Success {
		httpRespond(w, resp.ErrorMsg, http.StatusInternalServerError)
		return
	}
	httpRespond(w, resp.Msg, http.StatusOK)
}

func httpRespond(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintln(w, msg)
}
