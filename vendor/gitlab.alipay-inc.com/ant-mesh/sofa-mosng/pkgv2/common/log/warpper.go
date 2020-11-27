package log

import (
	mlog "mosn.io/pkg/log"
	"strconv"
	"sync/atomic"
	"time"
)

type Level mlog.Level

const (
	FATAL Level = iota
	ERROR       = Level(mlog.ERROR)
	WARN        = Level(mlog.WARN)
	INFO        = Level(mlog.INFO)
	DEBUG       = Level(mlog.DEBUG)
	TRACE       = Level(mlog.TRACE)
	RAW         = Level(mlog.RAW)
)

func (l Level) Level() mlog.Level {
	return mlog.Level(l)
}

func ConvertLevel(level mlog.Level) Level {
	return Level(level)
}

func GetOrCreateMosngLogger(loggerName, output string, roller *Roller) (Logger, error) {
	if roller == nil {
		roller = &Roller{}
	}
	if ml, err := mlog.GetOrCreateLogger(output, roller.Roller); err == nil {
		logger := &MosngLogger{
			name:   loggerName,
			Logger: ml,
			level:  INFO.Level(),
		}
		return logger, nil
	} else {
		return nil, err
	}

}

type Logger interface {
	mlog.ErrorLogger
}

type Roller struct {
	*mlog.Roller
}

type MosngLogger struct {
	name string
	*mlog.Logger
	level mlog.Level
}

var (
	// lastTime is used to cache time
	lastTime atomic.Value
)

// timeCache is used to reduce format
type timeCache struct {
	t int64
	s string
}

func logTime() string {
	var s string
	t := time.Now()
	now := t.Unix()
	value := lastTime.Load()
	if value != nil {
		last := value.(*timeCache)
		if now <= last.t {
			s = last.s
		}
	}
	if s == "" {
		s = t.Format("2006-01-02 15:04:05")
		lastTime.Store(&timeCache{now, s})
	}
	mi := t.UnixNano() % 1e9 / 1e6
	s = s + "," + strconv.Itoa(int(mi))
	return s
}

// default logger format:
// {time} [{level}] {content}
func (m *MosngLogger) formatter(lvPre string, format string) string {
	return logTime() + " - " + format
}

func (l *MosngLogger) codeFormatter(lvPre, errCode, format string) string {
	return logTime() + " " + lvPre + " [" + errCode + "] " + format
}

func (m *MosngLogger) Infof(fmt string, args ...interface{}) {
	if m.level >= mlog.INFO {
		s := m.formatter(mlog.InfoPre, fmt)
		m.Logger.Printf(s, args...)
	}
}

func (m *MosngLogger) Warnf(fmt string, args ...interface{}) {
	if m.level >= mlog.WARN {
		s := m.formatter(mlog.WarnPre, fmt)
		m.Logger.Printf(s, args...)
	}
}

func (m *MosngLogger) Errorf(fmt string, args ...interface{}) {
	if m.level >= mlog.ERROR {
		s := m.formatter(mlog.ErrorPre, fmt)
		m.Logger.Printf(s, args...)
	}
}

func (m *MosngLogger) Debugf(fmt string, args ...interface{}) {
	if m.level >= mlog.DEBUG {
		s := m.formatter(mlog.DebugPre, fmt)
		m.Logger.Printf(s, args...)
	}
}

// Alertf is a wrapper of Errorf
func (m *MosngLogger) Alertf(alert string, format string, args ...interface{}) {
	if m.level >= mlog.ERROR {
		s := m.codeFormatter(mlog.ErrorPre, alert, format)
		m.Logger.Printf(s, args...)
	}
}

func (m *MosngLogger) Tracef(format string, args ...interface{}) {
	if m.level >= mlog.TRACE {
		s := m.formatter(mlog.TracePre, format)
		m.Logger.Printf(s, args...)
	}
}

func (m *MosngLogger) Fatalf(format string, args ...interface{}) {
	s := m.formatter(mlog.FatalPre, format)
	m.Logger.Fatalf(s, args...)
}

func (m *MosngLogger) Fatal(args ...interface{}) {
	args = append([]interface{}{
		m.formatter(mlog.FatalPre, ""),
	}, args...)
	m.Logger.Fatal(args...)
}

func (m *MosngLogger) Fatalln(args ...interface{}) {
	args = append([]interface{}{
		m.formatter(mlog.FatalPre, ""),
	}, args...)
	m.Logger.Fatalln(args...)
}

// SetLogLevel updates the log level
func (m *MosngLogger) SetLogLevel(l mlog.Level) {
	m.level = l
}

// GetLogLevel returns the logger's level
func (m *MosngLogger) GetLogLevel() mlog.Level {
	return m.level
}
