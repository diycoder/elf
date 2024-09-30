package elf

import (
	"testing"

	"github.com/diycoder/elf/plugin"
	"github.com/diycoder/elf/plugin/apollo"
	"github.com/diycoder/elf/plugin/nacos"
	"github.com/diycoder/elf/plugin/store"
)

func TestPlugin(t *testing.T) {
	if err := initPlugin(); err != nil {
		t.Error(err)
		return
	}
	t.Log("OK")
	t.Log("--->>", nacos.Get("default", "a").String(""))
}

func initPlugin() error {
	// apollo config init
	if err := apollo.NewOptions(
		apollo.WithAddress("http://apollo.diycoder.com:8080"), // apollo 地址
		apollo.WithNamespace("application"),                   // apollo namespace
		apollo.WithAppID("user-api"),                          // apollo appID
		apollo.WithCluster("dev"),                             // apollo cluster
		apollo.WithBackUp(true),                               // 是否备份
	).Store(); err != nil {
		return err
	}

	// redis、mysql config init
	if err := store.NewOptions(
		store.WithStoreConfig("./env/store.yml"), // mysql、redis 配置文件(读取apollo配置)
	).Store(); err != nil {
		return err
	}

	// nacos config init
	if err := nacos.NewOptions(
		nacos.WithNamespaceId("5fe7e4db-424c-40b9-bb6f-ac893f049d24"),
		nacos.WithAddress("127.0.0.1:8848"),
		nacos.WithWatchConfig("./env/nacos_watch.yaml"),
	).Store(); err != nil {
		return err
	}

	// load plugin（version、log、nacos、apollo、mysql、redis）
	plugins := make([]plugin.Plugin, 0)
	plugins = append(
		DefaultPlugins(),   // 注册elf默认插件
		apollo.NewPlugin(), // 注册apollo插件
		nacos.NewPlugin(),  // 注册nacos插件
		store.NewPlugin(),  // 注册store插件
	)
	if err := InitPlugins(plugins...); err != nil {
		return err
	}
	return nil
}
