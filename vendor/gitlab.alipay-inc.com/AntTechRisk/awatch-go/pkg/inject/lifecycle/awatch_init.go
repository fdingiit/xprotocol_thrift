package lifecycle

import (
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/common"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/fault"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/http"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/log"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/rule"
	"sync"
)

var lock sync.Mutex

var awatchLogDirPath string
var awatchDiskInjectDirPath string

// awatch配置初始化，mosn解析完配置后，将配置传给awatch-go
func InitAwatchConfig(awatchConfig common.AwatchConfig) {
	awatchLogDirPath = awatchConfig.AwatchLogDirPath
	awatchDiskInjectDirPath = awatchConfig.AwatchDiskInjectDirPath
}

// awatch初始化
func InitAwatch() {
	lock.Lock()
	defer lock.Unlock()

	// 初始化日志，只做一次
	log.InitAwatchLogger(awatchLogDirPath)

	// 初始化磁盘注入目标目录
	fault.InitDiskInjector(awatchDiskInjectDirPath)

	// 初始化http server
	http.InitAwatchHttp()

	// 打开注入总控开关
	common.InjectOff = false

	log.DefaultLogger.Info("awatch initialized")
}

// awatch销毁
func DestroyAwatch() {
	lock.Lock()
	defer lock.Unlock()

	// 关闭注入总控开关，此时正在进行的注入会停止
	common.InjectOff = true

	// 清除所有的注入规则
	rule.ClearRules()

	// 关闭http server
	http.DestroyAwatchHttp()

	log.DefaultLogger.Info("awatch destroyed")
}
