package rollingwriter

import "github.com/natefinch/lumberjack"

type Option struct {
	maxsize    int
	maxbackups int
	maxAge     int
	localTime  bool
	compress   bool
}

func NewOption() *Option {
	return &Option{}
}

func (o *Option) SetMaxSize(i int) *Option    { o.maxsize = i; return o }
func (o *Option) SetMaxBackups(i int) *Option { o.maxbackups = i; return o }
func (o *Option) SetMaxAge(i int) *Option     { o.maxAge = i; return o }
func (o *Option) EnableLocalTime() *Option    { o.localTime = true; return o }
func (o *Option) EnableCompress() *Option     { o.compress = true; return o }

type RollingWriter struct {
	logger *lumberjack.Logger
}

func New(filename string, option *Option) *RollingWriter {
	return &RollingWriter{
		logger: &lumberjack.Logger{
			Filename:   filename,
			MaxSize:    option.maxsize,
			MaxBackups: option.maxbackups,
			MaxAge:     option.maxAge,
			LocalTime:  option.localTime,
			Compress:   option.compress,
		},
	}
}

func (rw *RollingWriter) Write(b []byte) (int, error) {
	return rw.logger.Write(b)
}

func (rw *RollingWriter) Close() error {
	return rw.logger.Close()
}
