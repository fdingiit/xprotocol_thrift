package common

type AwatchConfig struct {
	// awatch日志目录，例如 /home/admin/a/b/
	AwatchLogDirPath string
	// awatch磁盘注入目标目录，例如 /home/admin/c/d/
	AwatchDiskInjectDirPath string
}
