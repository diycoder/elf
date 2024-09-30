package log

import (
	"io"
	slog "log"
	"os"
	"strings"

	nglog "github.com/diycoder/elf/kit/log"
	"github.com/diycoder/elf/kit/log/writer/rotate"
)

func Init(cfg map[string]string) error {
	if fn, ok := cfg[FileNameKey]; ok && fn != "" {
		_ = SetLogFilename(fn)
	}

	if path, ok := cfg[PathKey]; ok && path != "" {
		if !strings.HasSuffix(path, "/") {
			path += "/"
		}
		_ = SetLogDir(path)
	}

	builders := []builder{
		newDebugZapLogger,
		newTraceZapLogger,
		newInfoZapLogger,
		newWarnZapLogger,
		newErrorZapLogger,
		newPanicZapLogger,
		newFatalZapLogger,
		newAccessZapLogger,
		newTracingZapLogger,
		newMicroZapLogger,
	}

	pl := NewPLogger()
	for _, builder := range builders {
		var (
			err error
			typ string
			zl  nglog.Logger
		)

		if typ, zl, err = builder(defaultLogMod); err != nil {
			return err
		}
		if err = pl.Register(typ, zl); err != nil {
			return err
		}
	}
	SetLogger(pl)
	replaceSysLogger()
	return nil
}

type builder func(int int) (string, nglog.Logger, error)

func newWriter(mod int, subDir, filename string) (io.Writer, error) {
	if mod == outTerminal {
		r := io.MultiWriter(os.Stdout)
		return r, nil
	}
	if mod == outFile {
		r, err := rotate.NewWriter(
			rotate.WithLogDir(defaultLogDir),
			rotate.WithLogSubDir(subDir),
			rotate.WithFileMode(defaultFileMode),
			rotate.WithFilename(filename),
			rotate.WithDuration(defaultRoateDuration),
		)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
	if mod == outTerminalAndFile {
		r, err := rotate.NewWriter(
			rotate.WithLogDir(defaultLogDir),
			rotate.WithLogSubDir(subDir),
			rotate.WithFileMode(defaultFileMode),
			rotate.WithFilename(filename),
			rotate.WithDuration(defaultRoateDuration),
		)
		if err != nil {
			return nil, err
		}
		r = io.MultiWriter(os.Stdout, r)
		return r, nil
	}
	return nil, nil
}

func newAccessZapLogger(mod int) (string, nglog.Logger, error) {
	const (
		typ      = "access"
		subDir   = "info/"
		filename = "gw-access.log"
	)

	r, err := newWriter(mod, subDir, filename)
	if err != nil {
		return "", nil, err
	}

	zl, err := nglog.New(
		defaultLogType,
		[]nglog.Option{
			nglog.WithWriter(r),
			nglog.Fields(defaultZapFields()),
			nglog.WithEncoder(defaultLogEncoder),
			nglog.WithEncoderCfg(defaultZapEncoderCfg()),
			nglog.WithLevelEnabler(defaultLogLevel),
			nglog.AddCaller(),
			nglog.AddCallerSkip(defaultLogCallerSkip),
		}...,
	)
	return typ, zl, err
}

func newDebugZapLogger(mod int) (string, nglog.Logger, error) {
	const (
		typ    = "debug"
		subDir = "debug/"
	)

	r, err := newWriter(mod, subDir, getDefaultLogFilename())
	if err != nil {
		return "", nil, err
	}

	zl, err := nglog.New(
		defaultLogType,
		[]nglog.Option{
			nglog.WithWriter(r),
			nglog.Fields(defaultZapFields()),
			nglog.WithEncoder(defaultLogEncoder),
			nglog.WithEncoderCfg(defaultZapEncoderCfg()),
			nglog.WithLevelEnabler(defaultLogLevel),
			nglog.AddCaller(),
			nglog.AddCallerSkip(defaultLogCallerSkip),
		}...,
	)
	return typ, zl, err
}

func newInfoZapLogger(mod int) (string, nglog.Logger, error) {
	const (
		typ    = "info"
		subDir = "info/"
	)

	r, err := newWriter(mod, subDir, getDefaultLogFilename())
	if err != nil {
		return "", nil, err
	}

	zl, err := nglog.New(
		defaultLogType,
		[]nglog.Option{
			nglog.WithWriter(r),
			nglog.Fields(defaultZapFields()),
			nglog.WithEncoder(defaultLogEncoder),
			nglog.WithEncoderCfg(defaultZapEncoderCfg()),
			nglog.WithLevelEnabler(defaultLogLevel),
			nglog.AddCaller(),
			nglog.AddCallerSkip(defaultLogCallerSkip),
		}...,
	)
	return typ, zl, err
}

func newWarnZapLogger(mod int) (string, nglog.Logger, error) {
	const (
		typ    = "warn"
		subDir = "warn/"
	)

	r, err := newWriter(mod, subDir, getDefaultLogFilename())
	if err != nil {
		return "", nil, err
	}

	zl, err := nglog.New(
		defaultLogType,
		[]nglog.Option{
			nglog.WithWriter(r),
			nglog.Fields(defaultZapFields()),
			nglog.WithEncoder(defaultLogEncoder),
			nglog.WithEncoderCfg(defaultZapEncoderCfg()),
			nglog.AddCaller(),
			nglog.WithLevelEnabler(defaultLogLevel),
			nglog.AddCallerSkip(defaultLogCallerSkip),
		}...,
	)
	return typ, zl, err
}

func newTraceZapLogger(mod int) (string, nglog.Logger, error) {
	const (
		typ    = "trace"
		subDir = "debug/"
	)

	r, err := newWriter(mod, subDir, getDefaultLogFilename())
	if err != nil {
		return "", nil, err
	}

	zl, err := nglog.New(
		defaultLogType,
		[]nglog.Option{
			nglog.WithWriter(r),
			nglog.Fields(defaultZapFields()),
			nglog.WithEncoder(defaultLogEncoder),
			nglog.WithEncoderCfg(defaultZapEncoderCfg()),
			nglog.WithLevelEnabler(defaultLogLevel),
			nglog.AddCaller(),
			nglog.AddCallerSkip(defaultLogCallerSkip),
		}...,
	)
	return typ, zl, err
}

func newErrorZapLogger(mod int) (string, nglog.Logger, error) {
	const (
		typ    = "error"
		subDir = "error/"
	)

	r, err := newWriter(mod, subDir, getDefaultLogFilename())
	if err != nil {
		return "", nil, err
	}

	zl, err := nglog.New(
		defaultLogType,
		[]nglog.Option{
			nglog.WithWriter(r),
			nglog.Fields(defaultZapFields()),
			nglog.WithEncoder(defaultLogEncoder),
			nglog.WithEncoderCfg(defaultZapEncoderCfg()),
			nglog.WithLevelEnabler(defaultLogLevel),
			nglog.AddCaller(),
			nglog.AddStacktrace(defaultLogLevel),
			nglog.AddCallerSkip(defaultLogCallerSkip),
		}...,
	)
	return typ, zl, err
}

func newPanicZapLogger(mod int) (string, nglog.Logger, error) {
	const (
		typ    = "panic"
		subDir = "error/"
	)

	r, err := newWriter(mod, subDir, getDefaultLogFilename())
	if err != nil {
		return "", nil, err
	}

	zl, err := nglog.New(
		defaultLogType,
		[]nglog.Option{
			nglog.WithWriter(r),
			nglog.Fields(defaultZapFields()),
			nglog.WithEncoder(defaultLogEncoder),
			nglog.WithEncoderCfg(defaultZapEncoderCfg()),
			nglog.WithLevelEnabler(defaultLogLevel),
			nglog.AddCaller(),
			nglog.AddStacktrace(defaultLogLevel),
			nglog.AddCallerSkip(defaultLogCallerSkip),
		}...,
	)
	return typ, zl, err
}

func newFatalZapLogger(mod int) (string, nglog.Logger, error) {
	const (
		typ    = "fatal"
		subDir = "error/"
	)

	r, err := newWriter(mod, subDir, getDefaultLogFilename())
	if err != nil {
		return "", nil, err
	}

	zl, err := nglog.New(
		defaultLogType,
		[]nglog.Option{
			nglog.WithWriter(r),
			nglog.Fields(defaultZapFields()),
			nglog.WithEncoder(defaultLogEncoder),
			nglog.WithEncoderCfg(defaultZapEncoderCfg()),
			nglog.WithLevelEnabler(defaultLogLevel),
			nglog.AddCaller(),
			nglog.AddStacktrace(defaultLogLevel),
			nglog.AddCallerSkip(defaultLogCallerSkip),
		}...,
	)
	if err != nil {
		return "", nil, err
	}

	return typ, zl, nil
}

func newMicroZapLogger(mod int) (string, nglog.Logger, error) {
	const (
		typ      = "micro"
		subDir   = "info/"
		filename = "micro.log"
	)

	r, err := newWriter(mod, subDir, filename)
	if err != nil {
		return "", nil, err
	}

	zl, err := nglog.New(
		defaultLogType,
		[]nglog.Option{
			nglog.WithWriter(r),
			nglog.Fields(defaultZapFields()),
			nglog.WithEncoder(defaultLogEncoder),
			nglog.WithEncoderCfg(defaultZapEncoderCfg()),
			nglog.WithLevelEnabler(defaultLogLevel),
			nglog.AddCaller(),
			nglog.AddCallerSkip(defaultLogCallerSkip),
		}...,
	)
	return typ, zl, err
}

func newTracingZapLogger(mod int) (string, nglog.Logger, error) {
	const (
		typ      = "tracing"
		subDir   = "info/"
		filename = "tracing.log"
	)

	r, err := newWriter(mod, subDir, filename)
	if err != nil {
		return "", nil, err
	}

	zl, err := nglog.New(
		defaultLogType,
		[]nglog.Option{
			nglog.WithWriter(r),
			nglog.WithEncoder(nglog.ConsoleEncoder),
			nglog.WithEncoderCfg(nglog.EncoderConfig{MessageKey: "msg"}),
			nglog.WithLevelEnabler(defaultLogLevel),
			nglog.AddCaller(),
			nglog.AddCallerSkip(defaultLogCallerSkip),
		}...,
	)
	return typ, zl, err
}

func defaultZapFields() map[string]interface{} {
	return map[string]interface{}{
		"host":    defaultHostName,
		"project": defaultProjectName,
	}
}

func defaultZapEncoderCfg() nglog.EncoderConfig {
	return nglog.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		EncodeLevel:   nglog.CapitalLevelEncoder,
		TimeKey:       "@timestamp",
		EncodeTime:    nglog.ChinaMilliTimeEncoder,
		CallerKey:     "caller",
		EncodeCaller:  nglog.JavaCallerEncoder,
		StacktraceKey: "detail",
	}
}

func defaultZapLoggerEncoderCfg() nglog.EncoderConfig {
	return nglog.EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		EncodeLevel:   nglog.CapitalLevelEncoder,
		TimeKey:       "@timestamp",
		EncodeTime:    nglog.ChinaMilliTimeEncoder,
		CallerKey:     "caller",
		EncodeCaller:  nglog.FullCallerEncoder,
		StacktraceKey: "detail",
	}
}

func replaceSysLogger() {
	slog.SetOutput(&sysLogger{})
}

type sysLogger struct{}

func (s *sysLogger) Write(p []byte) (n int, err error) {
	WithField("type", "micro").
		WithField("notice", "true").Info(string(p))
	return len(p), nil
}

func copyFields(src map[string]interface{}) map[string]interface{} {
	dst := make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
