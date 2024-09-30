package mysql

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/diycoder/elf/plugin/apollo"
	"github.com/diycoder/elf/plugin/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

var dbMap sync.Map

type dbconfig struct {
	User            string `json:"user"`
	Password        string `json:"password"`
	Protol          string `json:"protol"`
	Host            string `json:"host"`
	Port            string `json:"port"`
	DBName          string `json:"db_name"`
	MaxOpenConn     int    `json:"max_open_conn"`
	MaxIdleConn     int    `json:"max_idle_conn"`
	ConnMaxLifetime int    `json:"conn_max_life_time"`
	ConnMaxIdleTime int    `json:"conn_max_idle_time"`
	Extra           string `json:"extra"`
}

func Load(namespace, key string) error {
	val := apollo.Get(namespace, key)
	if val.String("") == "" {
		return fmt.Errorf("watch mysql 读取json为空")
	}

	var conf dbconfig
	err := json.Unmarshal(val.Bytes(), &conf)
	if err != nil {
		return err
	}
	db, err := newRedisPool(&conf)
	if err != nil {
		log.Errorf("mysql 连接连接失败 error: %v", err)
		return err
	}
	if err := db.Ping(); err != nil {
		log.Errorf("mysql ping 失败 error: %v", err)
		return err
	}
	dbMap.Store(key, db)
	log.Infof("rebuild %v mysql connection pool done", key)

	return nil
}

// new mysql pool from config center by watch
func newRedisPool(c *dbconfig) (*sqlx.DB, error) {
	echoConfig(c)
	db, err := addDBTrace(c)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func addDBTrace(c *dbconfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%s)/%s", c.User, c.Password, c.Protol, c.Host, c.Port, c.DBName)
	if c.Extra != "" {
		dsn += "?" + c.Extra
	}

	db, err := otelsqlx.Open("mysql", dsn,
		otelsql.WithAttributes(semconv.DBSystemMySQL),
		otelsql.WithDBName(c.DBName),
	)
	if err != nil {
		log.Errorf("mysql open:%s, err:%v", dsn, err)
		return nil, err
	}
	db.SetMaxOpenConns(c.MaxOpenConn)
	db.SetMaxIdleConns(c.MaxIdleConn)
	db.SetConnMaxIdleTime(time.Duration(c.ConnMaxLifetime) * time.Second)
	db.SetConnMaxLifetime(time.Duration(c.ConnMaxIdleTime) * time.Second)

	return db, nil
}

func echoConfig(c *dbconfig) {
	bytes, err := json.Marshal(c)
	if err != nil {
		return
	}
	cfg := make(map[string]interface{})
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		return
	}
	delete(cfg, "user")
	delete(cfg, "password")
	byteStr, err := json.Marshal(&cfg)
	if err != nil {
		return
	}
	log.Infof("apollo mysql config %v", string(byteStr))
}

// GetDB get mysql client by key
func GetDB(key string) (*sqlx.DB, error) {
	value, ok := dbMap.Load(key)
	if ok {
		return value.(*sqlx.DB), nil
	}

	log.Errorf("get mysql client:%v failed ", key)
	return nil, fmt.Errorf("get mysql client:%v failed ", key)
}
