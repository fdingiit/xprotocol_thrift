package sofalogger

type GeminiLogger struct {
	errl  Logger
	infol Logger
}

func NewGeminiLogger(infol, errl Logger) *GeminiLogger {
	return &GeminiLogger{
		infol: infol,
		errl:  errl,
	}
}

func (gl *GeminiLogger) Infof(format string, a ...interface{}) {
	gl.infol.Infof(format, a...)
}

func (gl *GeminiLogger) Debugf(format string, a ...interface{}) {
	gl.infol.Debugf(format, a...)
}

func (gl *GeminiLogger) Errorf(format string, a ...interface{}) {
	gl.errl.Errorf(format, a...)
}
