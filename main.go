package main

import (
    "fmt"
    "os"
    log "github.com/Sirupsen/logrus"
    "github.com/codegangsta/cli"
    "github.com/auroralaboratories/corona-api/util"
    "github.com/auroralaboratories/corona-api/modules/session"
    "github.com/auroralaboratories/corona-api/modules/command"
)

func main(){
    app                      := cli.NewApp()
    app.Name                  = util.ApplicationName
    app.Usage                 = util.ApplicationSummary
    app.Version               = util.ApplicationVersion
    app.EnableBashCompletion  = false
    app.Before                = func(c *cli.Context) error {
        if c.Bool(`quiet`) {
            util.ParseLogLevel(`quiet`)
        }else{
            util.ParseLogLevel(c.String(`log-level`))
        }

        log.Infof("%s v%s started at %s", util.ApplicationName, util.ApplicationVersion, util.StartedAt)

        api := NewApi()

        api.Address     = c.String(`address`)
        api.Port        = c.Int(`port`)

        if l := c.StringSlice(`modules`); l != nil && len(l) > 0 {
            api.ModulesList = l
        }

        if err := api.Init(); err == nil {
            if err := api.Serve(); err != nil {
                return fmt.Errorf("Failed to start API: %v", err)
            }
        }else{
            return fmt.Errorf("Failed to initialize API: %v", err)
        }

        return nil
    }

    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name:   `log-level, L`,
            Usage:  `Level of log output verbosity`,
            Value:  `info`,
            EnvVar: `LOGLEVEL`,
        },
        cli.BoolFlag{
            Name:   `quiet, q`,
            Usage:  `Don't print any log output to standard error`,
        },
        cli.StringFlag{
            Name:   `address, a`,
            Usage:  `The address the API server should listen on`,
            Value:  `127.0.0.1`,
        },
        cli.IntFlag{
            Name:   `port, p`,
            Usage:  `The port the API server should listen on`,
            Value:  25672,
        },
        cli.StringSliceFlag{
            Name:   `modules, m`,
            Usage:  `The set of modules to load on startup (default: all)`,
        },
    }

    app.Commands = append(app.Commands, util.RegisterSubcommands()...)
    app.Commands = append(app.Commands, session.RegisterSubcommands()...)
    app.Commands = append(app.Commands, command.RegisterSubcommands()...)

    app.Run(os.Args)
}