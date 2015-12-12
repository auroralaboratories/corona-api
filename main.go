package main

import (
  "flag"
  "github.com/ghetzel/go-logger"
)

const (
    DEFAULT_STATIC_ASSETS_PATH string = "/usr/share/corona/assets"
)

type CLIOptions struct {
    logLevel         *string
    logFilename      *string
    staticAssetsRoot *string
}

var logger        gologger.Logger
var api           CoronaAPI
var options       CLIOptions

func PanicIfErr(err error) {
    if err != nil {
        panic(err)
    }
}


func init_cli_arguments(){
    options.logLevel         = flag.String("level",       "debug",                    "Level of logging verbosity")
    options.logFilename      = flag.String("logfile",     "-",                        "The file to log output to, or dash (-) for standard output")
    options.staticAssetsRoot = flag.String("static-root", DEFAULT_STATIC_ASSETS_PATH, "Path where the API will serve static assets from")

    flag.Parse()
}

func main() {
    logger.Debug("Initializing Corona")
    init_cli_arguments()
    logger.Init(*options.logFilename, *options.logLevel)

    api.Options = CoronaOptions{
        StaticRoot: *options.staticAssetsRoot,
    }

    api.Init()


    logger.Infof("Starting Corona on %s:%d", api.Interface, api.Port)
    err := api.Serve()

    if err != nil {
        logger.Errorf("Error launching Corona: %s", err)
    }

}
