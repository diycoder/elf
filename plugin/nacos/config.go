package nacos

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/diycoder/elf/config"
	"github.com/diycoder/elf/plugin"

	"github.com/urfave/cli/v2"
)

type nacos struct {
	opts *Options
}

func (c *nacos) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "nacos_address",
			Usage:   "Set the address of nacos.",
			EnvVars: []string{"NACOS_ADDRESS"},
		},
		&cli.StringFlag{
			Name:    "nacos_watch_config",
			Value:   "./env/nacos_watch.yaml",
			Usage:   "Set the watch config of nacos .",
			EnvVars: []string{"NACOS_WATCH_CONFIG"},
		},
		&cli.StringFlag{
			Name:    "nacos_namespace_id",
			Usage:   "Set the namespace id of nacos .",
			EnvVars: []string{"NACOS_NAMESPACE_ID"},
		},
		&cli.StringFlag{
			Name:    "nacos_secret_key",
			Usage:   "Set the secret key of nacos .",
			EnvVars: []string{"NACOS_SECRET_KEY"},
		},
		&cli.StringFlag{
			Name:    "nacos_access_key",
			Usage:   "Set the access key of nacos .",
			EnvVars: []string{"NACOS_ACCESS_KEY"},
		},
	}
}

// Sub-commands
func (c *nacos) Commands() []*cli.Command {
	return nil
}

// Handle is the middleware handler for HTTP requests. We pass in
// the existing handler so it can be wrapped to create a call chain.
func (c *nacos) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			// serve the request
			h.ServeHTTP(rw, r)
		})
	}
}

func (c *nacos) Init(ctx *cli.Context) error {
	// init nacos param
	if err := c.getConfigure(ctx); err != nil {
		return err
	}
	if err := load(c.opts); err != nil {
		return err
	}

	// init loader source
	for _, l := range config.Loaders {
		loader := l()
		if loader == nil {
			return fmt.Errorf("config plugin error - new loader.%v", "")
		}
		if err := loader.Load(); err != nil {
			return err
		}
	}

	// init watch source
	for _, w := range config.Watchers {
		watch := w()
		if watch == nil {
			return fmt.Errorf("config plugin error - new watcher.%v", "")
		}
		if err := watch.Init(); err != nil {
			return err
		}
	}

	return nil
}

// get config info
func (c *nacos) getConfigure(ctx *cli.Context) error {
	address := ctx.String("nacos_address")
	if address == "" {
		return errors.New("nacos address is empty")
	}

	namespaceId := ctx.String("nacos_namespace_id")
	if namespaceId == "" {
		return errors.New("nacos namespace id is empty")
	}

	watchConfig := ctx.String("nacos_watch_config")
	if watchConfig == "" {
		return errors.New("nacos watch config is empty")
	}

	secretKey := ctx.String("nacos_secret_key")
	accessKey := ctx.String("nacos_access_key")

	c.opts = &Options{
		Address:     address,
		SecretKey:   secretKey,
		AccessKey:   accessKey,
		WatchConfig: watchConfig,
		NamespaceId: namespaceId,
	}

	return nil
}

// Name of the plugin
func (c *nacos) String() string {
	return "nacos_config"
}

func NewPlugin() plugin.Plugin {
	c := &nacos{}
	return c
}
