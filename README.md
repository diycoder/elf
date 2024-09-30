### 项目插件


### 插件：

- 日志插件: `log`  
- 版本插件: `version`
- Nacos插件: `nacos`  
- Apollo插件: `apollo`  
- Store插件: `store`  

#### 示例

```go
package main

import (
	"context"
	"time"

	"github.com/diycoder/elf"
	"github.com/diycoder/elf/plugin"
	"github.com/diycoder/elf/plugin/apollo"
	"github.com/diycoder/elf/plugin/log"
	"github.com/diycoder/elf/plugin/nacos"
	"github.com/diycoder/elf/plugin/store"
	"github.com/diycoder/elf/plugin/store/mysql"
	"github.com/diycoder/elf/plugin/store/redis"
)

func main() {
	if err := Init(); err != nil {
		log.Error(err)
		return
	}
	// key为mysql配置模板json对应的key
	db, err := mysql.GetDB("mudutv")
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("db ping:%v", db.Ping())
	// key为redis配置模板json对应的key
	rds, err := redis.GetClient("default")
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("redis ping:%v", rds.Ping(context.Background()))

	for {
		time.Sleep(time.Second)
		log.Info(apollo.Get("service", "vote_max_count").Int(0)) // apollo 获取对应配置
	}
}

func Init() error {
	// nacos config init
	if err := nacos.NewOptions(
		nacos.WithAddress("http://10.10.8.16:8848"),                   // nacos 地址
		nacos.WithNamespaceID("184b882a-bd12-4498-9377-ae6dda4f4a98"), // nacos namespaceID
		nacos.WithSeviceName("mdl-nacos-test"),                        // nacos 注册服务名称
		nacos.WithSevicePort(8080),                                    // nacos 注册服务监听端口号
		nacos.WithCluster("DEFAULT"),                                  // nacos 集群名称(默认DEFAULT)
		nacos.WithGroup("mdl"),                                        // nacos 组名称(默认DEFAULT_GROUP)
	).Store(); err != nil {
		return err
	}

	// apollo config init
	if err := apollo.NewOptions(
		apollo.WithAddress("http://119.3.102.128:8080"), // apollo 地址
		apollo.WithNamespace("application,service"),     // apollo namespace
		apollo.WithAppID("mudutv-company-srv"),          // apollo appID
		apollo.WithCluster("dev"),                       // apollo cluster
		apollo.WithBackUp(true),                         // 是否备份
	).Store(); err != nil {
		return err
	}

	// redis、mysql config init
	if err := store.NewOptions(
		store.WithStoreConfig("./env/dev/store.yml"), // mysql、redis 配置文件(读取apollo配置)
	).Store(); err != nil {
		return err
	}

	// load plugin（version、log、nacos、apollo、mysql、redis）
	plugins := make([]plugin.Plugin, 0)
	plugins = append(elf.DefaultPlugins(), nacos.NewPlugin(), apollo.NewPlugin(), store.NewPlugin())
	if err := elf.InitPlugins(plugins...); err != nil {
		return err
	}
	return nil
}

```