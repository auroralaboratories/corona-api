package main

import (
	"fmt"
	"github.com/auroralaboratories/corona-api/modules/command"
	"github.com/auroralaboratories/corona-api/modules/session"
	"github.com/auroralaboratories/corona-api/util"
	"github.com/ghetzel/cli"
	"github.com/op/go-logging"
	"os"
)

var log = logging.MustGetLogger(`main`)

func main() {
	var api *API

	app := cli.NewApp()
	app.Name = util.ApplicationName
	app.Usage = util.ApplicationSummary
	app.Version = util.ApplicationVersion
	app.EnableBashCompletion = false

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   `log-level, L`,
			Usage:  `Level of log output verbosity`,
			Value:  `info`,
			EnvVar: `LOGLEVEL`,
		},
		cli.StringFlag{
			Name:  `address, a`,
			Usage: `The address the API server should listen on`,
			Value: `127.0.0.1`,
		},
		cli.IntFlag{
			Name:  `port, p`,
			Usage: `The port the API server should listen on`,
			Value: 25672,
		},
		cli.StringSliceFlag{
			Name:  `modules, m`,
			Usage: `The set of modules to load on startup (default: all)`,
		},
	}

	app.Commands = append(app.Commands, util.RegisterSubcommands()...)
	app.Commands = append(app.Commands, session.RegisterSubcommands()...)
	app.Commands = append(app.Commands, command.RegisterSubcommands()...)

	app.Before = func(c *cli.Context) error {
		logging.SetFormatter(logging.MustStringFormatter(`%{color}%{level:.4s}%{color:reset}[%{id:04d}] %{message}`))

		if level, err := logging.LogLevel(c.String(`log-level`)); err == nil {
			logging.SetLevel(level, ``)
		} else {
			return err
		}

		log.Infof("Starting %s %s", c.App.Name, c.App.Version)

		api = NewApi()

		api.Address = c.String(`address`)
		api.Port = c.Int(`port`)

		if l := c.StringSlice(`modules`); l != nil && len(l) > 0 {
			api.ModulesList = l
		}

		if err := api.Init(); err == nil {
			if err := api.Serve(); err != nil {
				return fmt.Errorf("Failed to start API: %v", err)
			}
		} else {
			return fmt.Errorf("Failed to initialize API: %v", err)
		}

		return nil
	}

	app.Run(os.Args)
}
