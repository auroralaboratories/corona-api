package util

import (
    "fmt"
    "os"
    "net/http"
    "time"
    log "github.com/Sirupsen/logrus"
    "github.com/codegangsta/cli"
    "gopkg.in/unrolled/render.v1"
)

const ApplicationName    = `corona-api`
const ApplicationSummary = `a REST API for building web-based graphical user session applications`
const ApplicationVersion = `0.0.5`

var StartedAt = time.Now()
var ApiRenderer = render.New()

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

func Respond(w http.ResponseWriter, code int, payload interface{}, err error) {
    response := make(map[string]interface{})
    response[`responded_at`] = time.Now().Format(time.RFC3339)
    response[`payload`]      = payload

    if code >= http.StatusBadRequest {
        response[`success`] = false

        if err != nil {
            response[`error`] = err.Error()
        }
    }else{
        response[`success`] = true
    }

    ApiRenderer.JSON(w, code, response)
}