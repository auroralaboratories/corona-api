package main

import (
	"github.com/ghetzel/go-logger"
)

var logger gologger.Logger
var sprinkles_api SprinklesAPI

func PanicIfErr(err error) {
  if err != nil {
    panic(err)
  }
}


func init() {
	logger.Init("-", "debug")
}

func main() {
	logger.Debug("Initializing Sprinkles")
	sprinkles_api.Init()

	logger.Infof("Starting Sprinkles on %s:%d", sprinkles_api.Interface, sprinkles_api.Port)
	err := sprinkles_api.Serve()

	if err != nil {
		logger.Errorf("Error launching Sprinkles: %s", err)
	}

}
