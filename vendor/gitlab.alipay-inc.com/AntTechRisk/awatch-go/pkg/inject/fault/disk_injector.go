package fault

import (
	"fmt"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/common"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/log"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/rule"
	"gitlab.alipay-inc.com/AntTechRisk/awatch-go/pkg/inject/util"
	"os"
	"runtime"
	"runtime/debug"
)

var fooFileDirPath = "/home/admin/logs/.chaos/"

type DiskInjector struct {
}

// 显式初始化磁盘注入的写入文件的目录，该目录由mosn配置指定
func InitDiskInjector(floodDiskDirPath string) {
	if floodDiskDirPath == "" {
		floodDiskDirPath = fooFileDirPath
	}
	os.Mkdir(floodDiskDirPath, os.ModePerm)
}

func (*DiskInjector) Inject(injectContext map[string]string, injectRule *rule.InjectRule) *InjectResult {
	log.DefaultLogger.Info("disk inject start,injectId=" + injectRule.InjectId)

	injectResult := new(InjectResult)

	// 此类故障必须有持续时间作为注入条件
	injectDurationTime := injectRule.InjectCondition.DurationTime
	if injectDurationTime <= 0 {
		injectResult.Success = false
		injectResult.ErrorMsg = "duration time condition not exist"
		return injectResult
	}

	count := runtime.NumCPU()

	log.DefaultLogger.Info(fmt.Sprintf("ocuupy disk with %d goroutine,injectId=%v", count, injectRule.InjectId))
	for i := 0; i < count; i++ {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.DefaultLogger.Error("goroutine recovered, stack:" + string(debug.Stack()))
				}
			}()

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

				// 创建并写入文件，文件目录固定，文件名随机，写入1M字节数据
				var content = make([]byte, 1024*1024)
				var fileName = "chaos." + util.RandString(6)
				floodDisk(fileName, string(content))

				// 手动调度其他goroutine的执行
				runtime.Gosched()
			}
		}()
	}

	injectResult.Success = true
	return injectResult
}

func floodDisk(fileName string, content string) {
	// 创建文件
	file, err := os.Create(fooFileDirPath + fileName)
	defer file.Close()
	if err != nil {
		return
	}

	// 写文件
	file.WriteString(content)
}
