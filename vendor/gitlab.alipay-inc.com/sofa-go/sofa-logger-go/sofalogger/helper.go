package sofalogger

import (
	"strings"
)

// ParseLevel parses the string to level.
func ParseLevel(level string) Level {
	switch strings.ToLower(level) {
	// Hot path
	case "info":
		return InfoLevel
	case "debug":
		return DebugLevel
	case "error":
		return ErrorLevel
	case "dpanic":
		return DPanicLevel
	case "panic":
		return PanicLevel
	case "fatal":
		return FatalLevel
	default:
		return InfoLevel
	}
}
