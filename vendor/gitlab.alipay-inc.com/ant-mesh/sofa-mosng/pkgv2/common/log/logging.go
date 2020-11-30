package log

import (
	mlog "mosn.io/mosn/pkg/log"
)

func ErrorLogger() Logger {
	if errorLogger == nil {
		return mlog.DefaultLogger
	}
	return errorLogger
}

func ProxyLogger() Logger {
	if proxyLogger == nil {
		return mlog.DefaultLogger
	}
	return proxyLogger
}

func DefaultLogger() Logger {
	if defaultLogger == nil {
		return mlog.DefaultLogger
	}
	return defaultLogger
}

func ConfigLogger() Logger {
	if configLogger == nil {
		return mlog.DefaultLogger
	}
	return configLogger
}

func StartLogger() Logger {
	return mlog.StartLogger
}

func MosnDefaultLogger() Logger {
	return mlog.DefaultLogger
}
