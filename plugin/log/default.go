package log

import (
	nglog "github.com/diycoder/elf/kit/log"
)

var defaultLog nglog.Logger = NewPLogger()

var (
	Debug      = defaultLog.Debug
	Debugf     = defaultLog.Debugf
	Trace      = defaultLog.Trace
	Tracef     = defaultLog.Tracef
	Info       = defaultLog.Info
	Infof      = defaultLog.Infof
	Warn       = defaultLog.Warn
	Warnf      = defaultLog.Warnf
	Error      = defaultLog.Error
	Errorf     = defaultLog.Errorf
	Panic      = defaultLog.Panic
	Panicf     = defaultLog.Panicf
	Fatal      = defaultLog.Fatal
	Fatalf     = defaultLog.Fatalf
	WithField  = defaultLog.WithField
	WithFields = defaultLog.WithFields
)

func SetLogger(logger nglog.Logger) {
	defaultLog = logger
	Debug = defaultLog.Debug
	Debugf = defaultLog.Debugf
	Trace = defaultLog.Trace
	Tracef = defaultLog.Tracef
	Info = defaultLog.Info
	Infof = defaultLog.Infof
	Warn = defaultLog.Warn
	Warnf = defaultLog.Warnf
	Error = defaultLog.Error
	Errorf = defaultLog.Errorf
	Panic = defaultLog.Panic
	Panicf = defaultLog.Panicf
	Fatal = defaultLog.Fatal
	Fatalf = defaultLog.Fatalf
	WithField = defaultLog.WithField
	WithFields = defaultLog.WithFields
}
