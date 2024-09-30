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
	db, err := mysql.GetDB("auth")
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
}

func Init() error {

	// apollo config init
	if err := apollo.NewOptions(
		apollo.WithAddress("http://127.0.0.1:8080"), // apollo 地址
		apollo.WithNamespace("application,service"),     // apollo namespace
		apollo.WithAppID("auth-srv"),          // apollo appID
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
	plugins = append(elf.DefaultPlugins(), apollo.NewPlugin(), store.NewPlugin())
	if err := elf.InitPlugins(plugins...); err != nil {
		return err
	}
	return nil
}

```