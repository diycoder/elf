package store

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/diycoder/elf/kit/runenv"
	"github.com/diycoder/elf/plugin/log"
	"gopkg.in/yaml.v3"
)

type metadata struct {
	// map类型
	Stores []*storeConf `yml:"stores"`
}

type storeConf struct {
	Type      string `yml:"type"`
	Key       string `yml:"key"`
	Namespace string `yml:"namespace"`
}

// 读取Yaml配置文件
func getStoreConf(inputPath string) (*metadata, error) {
	// 获取绝对地址
	absFilename, _ := filepath.Abs(inputPath)
	yamlFile, err := ioutil.ReadFile(absFilename)
	if err != nil {
		log.Errorf("read store config err:%v", err.Error())
		return nil, err
	}
	metadata := new(metadata)
	err = yaml.Unmarshal(yamlFile, metadata)
	if err != nil {
		log.Errorf("unmarshal store config err:%v", err.Error())
		return nil, err
	}
	return metadata, nil
}

type Options struct {
	StoreCfg string `json:"store_cfg"`
}

func NewOptions(opts ...Option) *Options {
	options := Options{}
	for _, o := range opts {
		o(&options)
	}
	return &options
}

type Option func(o *Options)

func WithStoreConfig(storeCfg string) Option {
	return func(o *Options) {
		o.StoreCfg = storeCfg
	}
}

func (o *Options) Store() error {
	typeOf := reflect.TypeOf(o)
	valueOf := reflect.ValueOf(o)

	for i := 0; i < valueOf.Elem().NumField(); i++ {
		field := typeOf.Elem().Field(i)
		tag := field.Tag.Get("json")
		key := strings.ToUpper(tag)
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
		default:
			value = ""
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
