package log

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	nglog "github.com/diycoder/elf/kit/log"
	"github.com/diycoder/elf/kit/log/writer/rotate"
	"github.com/diycoder/elf/kit/runenv"
	"github.com/diycoder/elf/utils/net"
)

var (
	defaultLogDir        = rotate.DefaultLogDir
	defaultLogLevel      = nglog.DebugLevel
	defaultRoateDuration = rotate.DefaultRotateDuration
	defaultFileMode      = rotate.DefaultFileMode
	defaultProjectName   = filepath.Base(os.Args[0])
	defaultLogFileName   = fmt.Sprint(defaultProjectName, ".log")
	defaultHostName      = getHost()
	defaultLogMod        = getLogMod()
)

const (
	defaultLogType       = nglog.ZapLogger
	defaultLogEncoder    = nglog.ConsoleEncoder
	defaultLogCallerSkip = 2
	outTerminal          = 0
	outFile              = 1
	outTerminalAndFile   = 2
	defaultOutTerminal   = outTerminal
)

func SetLogDir(dir string) error {
	if len(dir) == 0 {
		return errors.New("dir can not be empty")
	} else if !strings.HasSuffix(dir, "/") {
		return errors.New("dir must end with /")
	}
	defaultLogDir = dir
	return nil
}

func SetProjectName(name string) error {
	if name == "" {
		return errors.New("project name can not be empty")
	}
	defaultProjectName = name
	return nil
}

func SetLogFilename(name string) error {
	if name == "" {
		return errors.New("log filename can not be empty")
	}
	defaultLogFileName = name
	return nil
}

func SetLogLevel(lv nglog.Level) error {
	if err := defaultLog.SetLogLevel(lv); err != nil {
		return err
	}
	defaultLogLevel = lv
	return nil
}

func getHost() string {
	ip, err := net.GetIP()
	if err != nil {
		ip, _ = os.Hostname()
	}
	return ip
}

func getLogMod() int {
	if e := os.Getenv("LOG_MOD"); e == "" {
		return defaultOutTerminal
	}
	return runenv.GetInt("LOG_MOD")
}

func GetLogDir() string {
	return defaultLogDir
}

func getDefaultLogFilename() string {
	return defaultLogFileName
}
