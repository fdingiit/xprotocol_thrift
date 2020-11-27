package tloggger

import (
	"os/user"

	"mosn.io/pkg/log"
)

var (
	TBaseLogger log.ErrorLogger
)

func InitTBaseLogger() {
	if TBaseLogger == nil {
		usr, err := user.Current()
		if err != nil {
			log.DefaultLogger.Warnf("[LOGGER_MANAGER] get or create tbase logger error, TBaseLogger is nil. error: %v", err)
			TBaseLogger = log.DefaultLogger
			return
		}
		userHome := usr.HomeDir
		TBaseLogger, err = CreateDefaultErrorLogger(userHome+"/logs/tbase-go-client/tbase-go-client.log", log.INFO)
		if TBaseLogger == nil {
			log.DefaultLogger.Warnf("[LOGGER_MANAGER] get or create tbase logger error, TBaseLogger is nil. error: %v", err)
			TBaseLogger = log.DefaultLogger
		}
	}
}

func CreateDefaultErrorLogger(output string, level log.Level) (log.ErrorLogger, error) {
	lg, err := log.GetOrCreateLogger(output, nil)
	if err != nil {
		return nil, err
	}
	return &log.SimpleErrorLog{
		Logger:    lg,
		Formatter: log.DefaultFormatter,
		Level:     level,
	}, nil
}
