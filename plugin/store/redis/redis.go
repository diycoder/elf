package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/diycoder/elf/plugin/apollo"
	"github.com/diycoder/elf/plugin/log"
	"github.com/redis/go-redis/extra/redisotel/v9"
	rds "github.com/redis/go-redis/v9"
)

var redisMap sync.Map

func Load(namespace, key string) error {
	if namespace == "" || key == "" {
		return fmt.Errorf("invalid config")
	}
	val := apollo.Get(namespace, key)
	if val.String("") == "" {
		return fmt.Errorf("watch redis 读取json为空")
	}

	var conf rdconfig
	err := json.Unmarshal(val.Bytes(), &conf)
	if err != nil {
		return err
	}

	if err := checkConfig(&conf); err != nil {
		log.Errorf("invalid redis config: %v", val.String(""))
		return err
	}
	client := newRedisPool(&conf)
	ping := client.Ping(context.Background())
	if ping.Val() == "" {
		log.Errorf("redis连接池连接失败 err: %v", ping.Err())
		return ping.Err()
	}
	redisMap.Store(key, client)
	log.Infof("rebuild %v redis connection pool done", key)

	return nil
}

func echoConfig(c *rdconfig) {
	bytes, err := json.Marshal(c)
	if err != nil {
		return
	}
	cfg := make(map[string]interface{})
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		return
	}
	delete(cfg, "password")
	byteStr, err := json.Marshal(&cfg)
	if err != nil {
		return
	}
	log.Infof("apollo redis config %v", string(byteStr))
}

type rdconfig struct {
	Addr         string `json:"addr"`
	Db           int    `json:"db"`
	Password     string `json:"password"`
	PoolSize     int    `json:"pool_size"`
	MinIdleConns int    `json:"min_idle_conns"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
}

func checkConfig(c *rdconfig) error {
	if c == nil {
		return fmt.Errorf("invalid redis addr")
	}
	if c.Addr == "" {
		return fmt.Errorf("invalid redis addr: %v", c.Addr)
	}
	return nil
}

// new redis pool from config center by watch
func newRedisPool(cfg *rdconfig) *rds.Client {
	echoConfig(cfg)
	options := rds.Options{
		Addr:     cfg.Addr,
		DB:       cfg.Db,
		Password: cfg.Password,
	}
	if cfg.PoolSize > 0 {
		options.PoolSize = cfg.PoolSize
	}
	if cfg.MinIdleConns > 0 {
		options.MinIdleConns = cfg.MinIdleConns
	}
	if cfg.ReadTimeout > 0 {
		options.ReadTimeout = time.Duration(cfg.ReadTimeout) * time.Second
	}
	if cfg.WriteTimeout > 0 {
		options.WriteTimeout = time.Duration(cfg.WriteTimeout) * time.Second
	}

	redisCli := rds.NewClient(&options)
	err := redisotel.InstrumentTracing(redisCli)
	if err != nil {
		log.Errorf("redisotel tracing err:%v", err)
		return nil
	}

	return redisCli
}

// GetClient get redis client by key
func GetClient(key string) (*rds.Client, error) {
	value, ok := redisMap.Load(key)
	if ok {
		return value.(*rds.Client), nil
	}

	log.Errorf("get redis client:%v failed", key)
	return nil, fmt.Errorf("get redis client:%v failed", key)
}
