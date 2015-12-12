package util

import (
    "fmt"
    "os"
    "time"
    log "github.com/Sirupsen/logrus"
    "github.com/codegangsta/cli"
)

const ApplicationName    = `corona-api`
const ApplicationSummary = `a REST API for building web-based graphical user session applications`
const ApplicationVersion = `0.0.5`

var StartedAt = time.Now()

func ParseLogLevel(level string) {
    log.SetOutput(os.Stderr)
    log.SetFormatter(&log.TextFormatter{
        ForceColors: true,
    })

    switch level {
    case `info`:
        log.SetLevel(log.InfoLevel)
    case `warn`:
        log.SetLevel(log.WarnLevel)
    case `error`:
        log.SetLevel(log.ErrorLevel)
    case `fatal`:
        log.SetLevel(log.FatalLevel)
    case `quiet`:
        log.SetLevel(log.PanicLevel)
    default:
        log.SetLevel(log.DebugLevel)
    }
}

func RegisterSubcommands() []cli.Command {
    return []cli.Command{
        {
            Name:        `version`,
            Usage:       `Output only the version string and exit`,
            Action:      func(c *cli.Context){
                fmt.Println(ApplicationVersion)
            },
        },
    }
}