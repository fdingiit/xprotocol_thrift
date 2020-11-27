package sofalogger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewAtomicLevel returns a new AtomicLevel.
func NewAtomicLevel(level Level) *AtomicLevel {
	al := zap.NewAtomicLevelAt(level)
	return &al
}

// Option represents an option.
type Option = zap.Option

// AddCaller return a enable caller option.
func AddCaller() Option {
	return zap.AddCaller()
}

// AddCallerSkip return a caller skip option.
func AddCallerSkip(skip int) Option {
	return zap.AddCallerSkip(skip)
}

func Hooks(hooks ...func(Entry) error) Option {
	return zap.Hooks(hooks...)
}

type PrimitiveArrayEncoder = zapcore.PrimitiveArrayEncoder

type TimeEncoder = zapcore.TimeEncoder

func DummyTimeEncoder(t time.Time, enc PrimitiveArrayEncoder) {
}

// Config represents the configuration for logger.
type Config struct {
	name        string
	level       *AtomicLevel
	options     []Option
	timeEncoder TimeEncoder
	hostname    bool
	pid         bool
}

// NewCallerConfig returns a Config with enable caller.
func NewCallerConfig() *Config {
	return &Config{
		options: []Option{
			AddCaller(),
		},
	}
}

// NewConfig returns a clean Config.
func NewConfig() *Config { return &Config{} }

// EnablePID enables the pid for logger.
func (c *Config) EnablePID() *Config {
	c.pid = true
	return c
}

// EnableHostname enables the hostname for logger.
func (c *Config) EnableHostname() *Config {
	c.hostname = true
	return c
}

// SetTimeEncoder sets the time encoder to use.
func (c *Config) SetTimeEncoder(te TimeEncoder) *Config {
	c.timeEncoder = te
	return c
}

// SetName sets the logger name.
func (c *Config) SetName(name string) *Config { c.name = name; return c }

// SetLevel sets the logger level.
func (c *Config) SetLevel(al *AtomicLevel) *Config { c.level = al; return c }

// AddOption adds a new option to config.
func (c *Config) AddOption(option Option) *Config {
	c.options = append(c.options, option)
	return c
}

// AddOptions adds options to config.
func (c *Config) AddOptions(options ...Option) *Config {
	c.options = append(c.options, options...)
	return c
}
