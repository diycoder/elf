package nacos

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/diycoder/elf/kit/runenv"
)

type Options struct {
	Address     string `json:"nacos_address"`
	NamespaceId string `json:"nacos_namespace_id"`
	WatchConfig string `json:"nacos_watch_config"`
	SecretKey   string `json:"nacos_secret_key"`
	AccessKey   string `json:"nacos_access_key"`
}

type Watch struct {
	Group  string `yaml:"group"`
	DataId string `yaml:"dataId"`
}

type Option func(o *Options)

func NewOptions(opts ...Option) *Options {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}
	return &options
}

// WithAddress sets the nacos address.
func WithAddress(addr string) Option {
	return func(o *Options) {
		o.Address = addr
	}
}

// WithWatchConfig sets the nacos watch data config file: group_id、data_id.
func WithWatchConfig(watchConfig string) Option {
	return func(o *Options) {
		o.WatchConfig = watchConfig
	}
}

// WithNamespaceId sets the nacos namespace.
func WithNamespaceId(namespaceId string) Option {
	return func(o *Options) {
		o.NamespaceId = namespaceId
	}
}

// WithAccessKey sets the nacos AccessKey.
func WithAccessKey(accesskey string) Option {
	return func(o *Options) {
		o.AccessKey = accesskey
	}
}

// WithSecretKey sets the nacos SecretKey.
func WithSecretKey(secretkey string) Option {
	return func(o *Options) {
		o.SecretKey = secretkey
	}
}

// Store set apollo config to env
func (o *Options) Store() error {
	typeOf := reflect.TypeOf(o)
	valueOf := reflect.ValueOf(o)

	for i := 0; i < valueOf.Elem().NumField(); i++ {
		field := typeOf.Elem().Field(i)
		key := strings.ToUpper(field.Tag.Get("json"))
		val := valueOf.Elem().Field(i)
		item := val.Interface()
		var value interface{}
		switch item.(type) {
		case string:
			value = val.String()
		case bool:
			value = val.Bool()
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			value = val.Int()
		}
		v := fmt.Sprintf("%v", value)
		if runenv.Exist(key) || v == "" {
			continue
		}
		if err := runenv.SetEnv(key, v); err != nil {
			return err
		}
	}
	return nil
}
