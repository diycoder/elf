package apollo

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/diycoder/elf/config"
	"github.com/diycoder/elf/plugin"

	"github.com/urfave/cli/v2"
)

type apollo struct {
	opts *Options
}

func (c *apollo) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "apollo_ip",
			Usage:   "Set the address of apollo.",
			EnvVars: []string{"APOLLO_IP"},
		},
		&cli.StringFlag{
			Name:    "apollo_appid",
			Usage:   "Set the appid of apollo config.",
			EnvVars: []string{"APOLLO_APPID"},
		},
		&cli.StringFlag{
			Name:    "apollo_namespaces",
			Usage:   "Set the namespace of apollo config.",
			EnvVars: []string{"APOLLO_NAMESPACES"},
		},
		&cli.StringFlag{
			Name:    "apollo_cluster",
			Usage:   "Set the cluster of apollo config.",
			EnvVars: []string{"APOLLO_CLUSTER"},
		},
		&cli.StringFlag{
			Name:    "apollo_backup_path",
			Value:   "./cache/apollo",
			Usage:   "Set the backup path of apollo config.",
			EnvVars: []string{"APOLLO_BACKUP_PATH"},
		},
		&cli.StringFlag{
			Name:    "apollo_backup",
			Value:   "true",
			Usage:   "Set the is backup of apollo config.",
			EnvVars: []string{"APOLLO_BACKUP"},
		},
	}
}

// Sub-commands
func (c *apollo) Commands() []*cli.Command {
	return nil
}

// Handle is the middleware handler for HTTP requests. We pass in
// the existing handler so it can be wrapped to create a call chain.
func (c *apollo) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			// serve the request
			h.ServeHTTP(rw, r)
		})
	}
}

func (c *apollo) Init(ctx *cli.Context) error {
	// init apollo param
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
func (c *apollo) getConfigure(ctx *cli.Context) error {
	address := ctx.String("apollo_ip")
	if address == "" {
		return errors.New("apollo address is empty")
	}

	appID := ctx.String("apollo_appid")
	if appID == "" {
		return errors.New("apollo app id is empty")
	}

	namespace := ctx.String("apollo_namespaces")
	if namespace == "" {
		return errors.New("apollo namespace is empty")
	}

	cluster := ctx.String("apollo_cluster")
	backupPath := ctx.String("apollo_backup_path")
	backup := ctx.Bool("apollo_backup")

	c.opts = &Options{
		Address:    address,
		AppID:      appID,
		Namespace:  namespace,
		Cluster:    cluster,
		Backup:     backup,
		BackupPath: backupPath,
	}

	return nil
}

// Name of the plugin
func (c *apollo) String() string {
	return "apollo_config"
}

func NewPlugin() plugin.Plugin {
	c := &apollo{}
	return c
}
