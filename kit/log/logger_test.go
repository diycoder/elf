package log

import (
	"os"
	"testing"
)

func TestLogger(t *testing.T) {
	logger, err := New(0,
		WithWriter(os.Stdout),
		WithEncoder(ConsoleEncoder),
		WithEncoderCfg(defaultZapEncoderCfg()),
		WithLevelEnabler(-1),
		AddCaller(),
	)
	if err != nil {
		t.Error(err)
		return
	}
	logger.Info("haha")
}

func defaultZapEncoderCfg() EncoderConfig {
	return EncoderConfig{
		MessageKey:    "msg",
		LevelKey:      "level",
		EncodeLevel:   CapitalLevelEncoder,
		TimeKey:       "@timestamp",
		EncodeTime:    ChinaMilliTimeEncoder,
		CallerKey:     "caller",
		EncodeCaller:  JavaCallerEncoder,
		StacktraceKey: "detail",
	}
}
