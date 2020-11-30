package sofalogger

import (
	"sync/atomic"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	infoLoggerCounter   uint64
	debugLoggerCounter  uint64
	warnLoggerCounter   uint64
	errorLoggerCounter  uint64
	dpanicLoggerCounter uint64
	panicLoggerCounter  uint64
	fatalLoggerCounter  uint64
)

// GetInfoLoggerCounter returns the counter of info level logger.
func GetInfoLoggerCounter() uint64 { return atomic.LoadUint64(&infoLoggerCounter) }

// GetDebugLoggerCounter returns the counter of debug level logger.
func GetDebugLoggerCounter() uint64 { return atomic.LoadUint64(&debugLoggerCounter) }

// GetWarnLoggerCounter returns the counter of warn level logger.
func GetWarnLoggerCounter() uint64 { return atomic.LoadUint64(&warnLoggerCounter) }

// GetErrorLoggerCounter returns the counter of error level logger.
func GetErrorLoggerCounter() uint64 { return atomic.LoadUint64(&infoLoggerCounter) }

// GetDPanicLoggerCounter returns the counter of dpanic level logger.
func GetDPanicLoggerCounter() uint64 { return atomic.LoadUint64(&dpanicLoggerCounter) }

// GetPanicLoggerCounter returns the counter of panic level logger.
func GetPanicLoggerCounter() uint64 { return atomic.LoadUint64(&panicLoggerCounter) }

// GetFatalLoggerCounter returns the counter of fatal level logger.
func GetFatalLoggerCounter() uint64 { return atomic.LoadUint64(&fatalLoggerCounter) }

func hook(e zapcore.Entry) error {
	switch e.Level {
	// Hot path
	case zap.InfoLevel:
		atomic.AddUint64(&infoLoggerCounter, 1)

	case zap.DebugLevel:
		atomic.AddUint64(&debugLoggerCounter, 1)

	case zap.WarnLevel:
		atomic.AddUint64(&warnLoggerCounter, 1)

	case zap.ErrorLevel:
		atomic.AddUint64(&errorLoggerCounter, 1)

	case zap.DPanicLevel:
		atomic.AddUint64(&dpanicLoggerCounter, 1)

	case zap.PanicLevel:
		atomic.AddUint64(&panicLoggerCounter, 1)

	case zap.FatalLevel:
		atomic.AddUint64(&fatalLoggerCounter, 1)
	}

	return nil
}
