package log

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	nglog "github.com/diycoder/elf/kit/log"
)

var (
	nopLogger, _              = nglog.New(nglog.ZapLogger, nglog.WithWriter(ioutil.Discard))
	stdLogger, _              = nglog.New(nglog.ZapLogger, nglog.WithWriter(os.Stdout))
	_            nglog.Logger = (*pLogger)(nil)
)

type pLogger struct {
	levelEnabler nglog.LevelEnabler
	fields       sync.Map
	selector     sync.Map
}

func NewPLogger() *pLogger {
	return &pLogger{
		levelEnabler: nglog.DebugLevel,
	}
}

func (pl *pLogger) Debug(args ...interface{}) {
	pl.doSelect(nglog.DebugLevel).Debug(args...)
}

func (pl *pLogger) Debugf(format string, args ...interface{}) {
	pl.doSelect(nglog.DebugLevel).Debugf(format, args...)
}

func (pl *pLogger) Info(args ...interface{}) {
	pl.doSelect(nglog.InfoLevel).Info(args...)
}

func (pl *pLogger) Infof(format string, args ...interface{}) {
	pl.doSelect(nglog.InfoLevel).Infof(format, args...)
}

func (pl *pLogger) Warn(args ...interface{}) {
	pl.doSelect(nglog.WarnLevel).Warn(args...)
}

func (pl *pLogger) Warnf(format string, args ...interface{}) {
	pl.doSelect(nglog.WarnLevel).Warnf(format, args...)
}

func (pl *pLogger) Error(args ...interface{}) {
	pl.doSelect(nglog.ErrorLevel).Error(args...)
}

func (pl *pLogger) Errorf(format string, args ...interface{}) {
	pl.doSelect(nglog.ErrorLevel).Errorf(format, args...)
}

func (pl *pLogger) Trace(args ...interface{}) {
	pl.Debug(args...)
}

func (pl *pLogger) Tracef(format string, args ...interface{}) {
	pl.Debugf(format, args...)
}

// Error 以上级别全部替换成 Error
func (pl *pLogger) Panic(args ...interface{}) {
	pl.doSelect(nglog.PanicLevel).Panic(args...)
}

func (pl *pLogger) Panicf(format string, args ...interface{}) {
	pl.doSelect(nglog.PanicLevel).Panicf(format, args...)
}

func (pl *pLogger) Fatal(args ...interface{}) {
	pl.doSelect(nglog.FatalLevel).Fatal(args...)
}

func (pl *pLogger) Fatalf(format string, args ...interface{}) {
	pl.doSelect(nglog.FatalLevel).Fatalf(format, args...)
}

func (pl *pLogger) WithField(key string, value interface{}) nglog.Logger {
	if key == "" {
		return pl
	}
	clone := pl.clone()
	clone.fields.Store(key, value)
	return clone
}

func (pl *pLogger) WithFields(fields map[string]interface{}) nglog.Logger {
	if len(fields) == 0 {
		return pl
	}
	clone := pl.clone()
	for k, v := range fields {
		clone.fields.Store(k, v)
	}

	return clone
}

func (pl *pLogger) doSelect(lv nglog.Level) nglog.Logger {
	if !pl.levelEnabler.Enabled(lv) {
		return nopLogger
	}

	fields := make(map[string]interface{})
	pl.fields.Range(func(key, value interface{}) bool {
		fields[key.(string)] = value
		return true
	})

	if v, ok := pl.fields.Load("type"); ok {
		if typ, ok := v.(string); ok {
			if zl, ok := pl.selector.Load(typ); ok && zl != nil {
				// 如果是 tracing 就删除掉
				if typ == "tracing" {
					delete(fields, "type")
				}
				return zl.(nglog.Logger).WithFields(fields)
			}
		}
	}

	if zl, ok := pl.selector.Load(lv.String()); ok && zl != nil {
		return zl.(nglog.Logger).WithFields(fields)
	}

	return stdLogger.WithFields(fields)
}

func (pl *pLogger) Select(lv nglog.Level) nglog.Logger {
	return pl.doSelect(lv)
}

func (pl *pLogger) Register(typ string, l nglog.Logger) error {
	if typ == "" {
		return errors.New("The type cannot be empty ")
	}

	if l == nil {
		return errors.New("Logger cannot be nil ")
	}

	if _, ok := pl.selector.LoadOrStore(typ, l); ok {
		return fmt.Errorf("Logger with type %s already registered ", typ)
	}

	return nil
}

// Changing level on the fly without app restart
func (pl *pLogger) SetLogLevel(lv nglog.Level) error {
	lvs := nglog.AllLevels()
	for _, level := range lvs {
		if level == lv {
			pl.levelEnabler = lv
			return nil
		}
	}
	return fmt.Errorf("Invalid level: %v ", lv)
}

func (pl *pLogger) clone() *pLogger {
	clone := &pLogger{
		levelEnabler: pl.levelEnabler,
	}

	pl.fields.Range(func(key, value interface{}) bool {
		clone.fields.Store(key, value)
		return true
	})
	pl.selector.Range(func(key, value interface{}) bool {
		clone.selector.Store(key, value)
		return true
	})

	return clone
}
