// Package daemon is a micro plugin for daemon gateway service
package version

import (
	"fmt"
	"net/http"
	"os"

	"github.com/diycoder/elf/plugin"
	"github.com/urfave/cli/v2"
)

var (
	Version     = ""
	GitRevision = ""
	GitBranch   = ""
	GoVersion   = ""
	BuildTime   = ""
	OSArch      = ""
)

type version struct{}

func (p *version) Flags() []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:  "version, v",
			Usage: "Set the version info",
		},
		&cli.StringFlag{
			Name: "test.timeout",
		},
		&cli.BoolFlag{
			Name: "test.v",
		},
		&cli.BoolFlag{
			Name: "test.paniconexit0",
		},
		&cli.StringFlag{
			Name: "test.run",
		},
		&cli.StringFlag{
			Name: "test.testlogfile",
		},
	}
}

func (p *version) Commands() []*cli.Command {
	return nil
}

func (p *version) Handler() plugin.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		})
	}
}

func (p *version) Init(ctx *cli.Context) error {
	if ctx.Bool("version") {
		fmt.Println("Version     : \t" + Version)
		fmt.Println("Git   branch: \t" + GitBranch)
		fmt.Println("Git revision: \t" + GitRevision)
		fmt.Println("Go   version: \t" + GoVersion)
		fmt.Println("Build   time: \t" + BuildTime)
		fmt.Println("OS/Arch     : \t" + OSArch)
		os.Exit(0)
	}
	return nil
}

func (p *version) String() string {
	return "version"
}

func NewPlugin() plugin.Plugin {
	return &version{}
}
