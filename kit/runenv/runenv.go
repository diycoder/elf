package runenv

import (
	"errors"
	"os"
	"strings"

	"github.com/diycoder/elf/utils/convert"
)

type REnv = string

const (
	DefaultREnv      = Dev
	Dev         REnv = "dev"
	Test        REnv = "test"
	Uat         REnv = "uat"
	Gray        REnv = "gray"
	Prod        REnv = "product"
	OutTerminal      = true
)

var runEnvKey = "RUN_ENV"

// Is reports whether the server is running in its env configuration
func Is(env REnv) bool {
	return strings.HasSuffix(GetRunEnv(), strings.ToLower(env))
}

func Not(env REnv) bool {
	return !Is(env)
}

// IsDev reports whether the server is running in its development configuration
func IsDev() bool {
	return Is(Dev)
}

// IsTest reports whether the server is running in its testing configuration
func IsTest() bool {
	return Is(Test)
}

// IsUat reports whether the server is running in its testing configuration
func IsUat() bool {
	return Is(Uat)
}

// IsGray reports whether the server is running in its gray configuration
func IsGray() bool {
	return Is(Gray)
}

// IsProd reports whether the server is running in its production configuration
func IsProd() bool {
	return Is(Prod)
}

// GetRunEnv Gets the current runtime environment
func GetRunEnv() (e REnv) {
	if e = os.Getenv(runEnvKey); e == "" {
		// Returns a specified default value (Dev) if an empty or invalid value is detected.
		e = DefaultREnv
	}
	return strings.ToLower(e)
}

// GetRunEnvKey Gets the key of the runtime environment
func GetRunEnvKey() string {
	return runEnvKey
}

// SetRunEnvKey Sets the key of the runtime environment
func SetRunEnvKey(key string) error {
	if key == "" {
		return errors.New("[runEnv] RunEnvKey cannot be empty")
	}
	runEnvKey = key
	return nil
}

// SetDevEnv Sets the env of the dev environment
func SetDevEnv(key, value string) error {
	if key == "" {
		return errors.New("[runEnv] RunEnvKey cannot be empty")
	}
	if !IsDev() {
		return errors.New("[runEnv] is not dev")
	}
	return os.Setenv(key, value)
}

// SetEnv Set the env of the evniroment
func SetEnv(key, value string) error {
	return os.Setenv(key, value)
}

func GetInt(key string) int {
	return convert.ToInt(os.Getenv(key))
}

func GetInt64(key string) int64 {
	return convert.ToInt64(os.Getenv(key))
}

func GetBool(key string) bool {
	return convert.ToBool(os.Getenv(key))
}

func GetFloat64(key string) float64 {
	return convert.ToFloat64(os.Getenv(key))
}

func Get(key string) string {
	return os.Getenv(key)
}

func Exist(key string) bool {
	val, exist := os.LookupEnv(key)
	return exist && len(val) > 0
}
