package main

import (
    "os"
    log "github.com/Sirupsen/logrus"
    "github.com/codegangsta/cli"
    "github.com/auroralaboratories/corona-api/util"
    "github.com/auroralaboratories/corona-api/modules/session"
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
    }

    app.Commands = append(app.Commands, util.RegisterSubcommands()...)
    app.Commands = append(app.Commands, session.RegisterSubcommands()...)

    app.Run(os.Args)
}