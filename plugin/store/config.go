package store

import (
	"errors"
	"net/http"

	"github.com/diycoder/elf/plugin"
	"github.com/diycoder/elf/plugin/log"
	"github.com/diycoder/elf/plugin/store/mysql"
	"github.com/diycoder/elf/plugin/store/redis"
	"github.com/urfave/cli/v2"
)

type store struct {
	storeCfg string
}

func (s *store) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "store_cfg",
			Usage:   "Set the file path of store config",
			EnvVars: []string{"STORE_CFG"},
		},
	}
}

func (s *store) Commands() []*cli.Command {
	return nil
}

func (c *store) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			// serve the request
			h.ServeHTTP(rw, r)
		})
	}
}

func (s *store) Init(ctx *cli.Context) error {
	s.storeCfg = ctx.String("store_cfg")
	if s.storeCfg == "" {
		return errors.New("store config file path is empty")
	}
	log.Infof("store config path:%v", s.storeCfg)
	if err := s.loadStore(s.storeCfg); err != nil {
		return err
	}
	return nil
}

func (s *store) loadStore(path string) error {
	md, err := getStoreConf(path)
	if err != nil {
		return err
	}
	for _, c := range md.Stores {
		switch c.Type {
		case "redis":
			if err := redis.Load(c.Namespace, c.Key); err != nil {
				return err
			}
		case "mysql":
			if err := mysql.Load(c.Namespace, c.Key); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *store) String() string {
	return "store"
}

func NewPlugin() plugin.Plugin {
	c := &store{}
	return c
}
