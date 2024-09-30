package log

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/diycoder/elf/plugin"

	"github.com/urfave/cli/v2"
)

type log struct {
	md map[string]string
}

const (
	DebugKey    = "debug"
	PathKey     = "path"
	FileNameKey = "log"
)

// Global Flags
func (l *log) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "log_setting",
			Value:   "log=gw.log",
			Usage:   "Set logger file and pretty or not. \"pretty=false\", \"log=gw.log\", \"path=/data/log\"",
			EnvVars: []string{"LOG_SETTING"},
		},
		&cli.StringFlag{
			Name:    "log_duration",
			Value:   "1h",
			Usage:   "Sets the time between rotation, e.g. \"24m\", \"24h\".",
			EnvVars: []string{"LOG_DURATION"},
		},
		&cli.StringFlag{
			Name:    "log_mod",
			Value:   "0",
			Usage:   "Sets the log mod, e.g. \"0\", \"1\", \"2\" .",
			EnvVars: []string{"LOG_MOD"},
		},
	}
}

// Sub-commands
func (l *log) Commands() []*cli.Command {
	return nil
}

// Handle is the middleware handler for HTTP requests. We pass in
// the existing handler so it can be wrapped to create a call chain.
func (l *log) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			// serve the request
			h.ServeHTTP(rw, r)
		})
	}
}

// Init called when command line args are parsed.
// The initialized cli.Context is passed in.
func (l *log) Init(ctx *cli.Context) error {
	conf := ctx.String("log_setting")
	if len(conf) == 0 {
		return nil
	}

	// write log to file while has -d flag or run in k8s
	l.md[DebugKey] = "true"

	// 特殊处理
	if ctx.Bool("daemon") {
		l.md[DebugKey] = "false"
	}

	if ctx.String("registry") == "kubernetes" {
		l.md[DebugKey] = "false"
	}

	logMod := ctx.String("log_mod")
	if mod, err := strconv.Atoi(logMod); err == nil && mod >= 0 {
		defaultLogMod = mod
	}

	// 日志切割间隔
	rotateDuration := ctx.String("log_duration")
	if d, err := time.ParseDuration(rotateDuration); err == nil && d > 0 {
		defaultRoateDuration = d
	} else if d == 0 {
		return errors.New("[log] The log_duration cannot be empty")
	} else {
		return err
	}

	return Init(l.md)
}

// Name of the plugin
func (l *log) String() string {
	return "log_setting"
}

func NewPlugin() plugin.Plugin {
	return &log{
		md: make(map[string]string),
	}
}
