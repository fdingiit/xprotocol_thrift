package main

import (
	"go.uber.org/zap"
)

var (
	ProtocolLogger *zap.SugaredLogger
)

func init() {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		"stdout",
		"./log",
	}
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	logger, _ := cfg.Build()
	ProtocolLogger = logger.Sugar()
}
