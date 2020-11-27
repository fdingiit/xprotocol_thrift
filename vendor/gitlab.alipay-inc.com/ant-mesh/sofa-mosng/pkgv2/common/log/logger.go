package log

import (
	"os/user"
)

var (
	configLogger  Logger
	errorLogger   Logger
	defaultLogger Logger
	proxyLogger   Logger
)

func init() {
	InitMosnLogger()
}

func InitMosnLogger() {
	usr, _ := user.Current()

	var err error
	logRoot := usr.HomeDir + "/logs/sofa-mosng/"
	if defaultLogger, err = GetOrCreateMosngLogger("mosng-default", logRoot+"mosng-default.log", nil); err != nil {
		DefaultLogger().Errorf("create %s logger error: %v", "mosng-default", err)
	}

	if errorLogger, err = GetOrCreateMosngLogger("mosng-error", logRoot+"mosng-error.log", nil); err != nil {
		DefaultLogger().Errorf("create %s logger error: %v", "mosng-error", err)
	}

	if proxyLogger, err = GetOrCreateMosngLogger("mosng-proxy", logRoot+"mosng-proxy.log", nil); err != nil {
		DefaultLogger().Errorf("create %s logger error: %v", "mosng-proxy", err)
	}

	if configLogger, err = GetOrCreateMosngLogger("mosng-config", logRoot+"mosng-config.log", nil); err != nil {
		DefaultLogger().Errorf("create %s logger error: %v", "mosng-config", err)
	}
}
