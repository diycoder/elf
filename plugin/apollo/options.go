package apollo

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/diycoder/elf/kit/runenv"
)

type Options struct {
	Address    string `json:"apollo_ip"`
	AppID      string `json:"apollo_appid"`
	Namespace  string `json:"apollo_namespaces"`
	Cluster    string `json:"apollo_cluster"`
	Backup     bool   `json:"apollo_backup"`
	BackupPath string `json:"apollo_backup_path"`
}

type Option func(o *Options)

func NewOptions(opts ...Option) *Options {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}
	return &options
}

func WithNamespace(namespace string) Option {
	return func(o *Options) {
		o.Namespace = namespace
	}
}

func WithAddress(addr string) Option {
	return func(o *Options) {
		o.Address = addr
	}
}

func WithBackupPath(backupPath string) Option {
	return func(o *Options) {
		o.BackupPath = backupPath
	}
}

func WithAppID(appID string) Option {
	return func(o *Options) {
		o.AppID = appID
	}
}

func WithCluster(cluster string) Option {
	return func(o *Options) {
		o.Cluster = cluster
	}
}

func WithBackUp(backup bool) Option {
	return func(o *Options) {
		o.Backup = backup
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
