package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	AwatchLogHome = "/home/admin/logs/awatch-go/"
)

var DefaultLogger *zap.Logger = zap.NewNop()
var AgentLogger *zap.Logger = zap.NewNop()
var InjectLogger *zap.Logger = zap.NewNop()

var initMark bool

// 显式初始化awatch的日志，日志目录由mosn配置指定
func InitAwatchLogger(logDirPath string) {
	if initMark == true {
		return
	}

	if logDirPath == "" {
		logDirPath = AwatchLogHome
	}
	DefaultLogger = initLogger(logDirPath + "awatch-default.log")
	AgentLogger = initLogger(logDirPath + "awatch-agent.log")
	InjectLogger = initLogger(logDirPath + "awatch-inject.log")

	initMark = true
}

func initLogger(logPath string) *zap.Logger {
	// rolling file
	hook := lumberjack.Logger{
		Filename: logPath,
		MaxSize:  256, // megabytes
		MaxAge:   14,
	}

	// zap encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	// log level
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(zapcore.InfoLevel)

	// log line number
	caller := zap.AddCaller()

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(&hook)),
		atomicLevel,
	)

	return zap.New(core, caller)
}
