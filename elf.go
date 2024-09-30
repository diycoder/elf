package elf

import (
	"fmt"
	"os"

	"github.com/diycoder/elf/config/cmd"
	"github.com/diycoder/elf/plugin"
	"github.com/diycoder/elf/plugin/log"
	"github.com/diycoder/elf/plugin/version"

	"github.com/urfave/cli/v2"
	"go.uber.org/atomic"
)

var (
	done atomic.Int32
	// ErrorReinitialized
	ErrorReinitialized = fmt.Errorf("plugins init have been called. please merge plugins. ")
)

const (
	Initialized   = 1
	Uninitialized = 0
)

func Init(plugins ...plugin.Plugin) error {
	if !done.CAS(Uninitialized, Initialized) {
		return ErrorReinitialized
	}

	app := cmd.App()
	oldBefore := app.Before
	app.Before = func(context *cli.Context) error {
		for _, p := range plugin.Plugins() {
			if err := p.Init(context); err != nil {
				app.CustomAppHelpTemplate = fmt.Sprintf("plugin %s init error: ", p.String())
				return cli.NewExitError(err, 1)
			}
		}
		app.Before = oldBefore
		return nil
	}

	for _, p := range plugins {
		app.Flags = append(app.Flags, p.Flags()...)
		if err := plugin.Register(p); err != nil {
			return err
		}
	}
	return app.Run(os.Args)
}

// InitPlugins initialize plugins
func InitPlugins(plugins ...plugin.Plugin) error {
	return Init(plugins...)
}

func DefaultPlugins() []plugin.Plugin {
	return []plugin.Plugin{
		log.NewPlugin(),
		version.NewPlugin(),
	}
}
