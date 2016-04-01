package command

import (
	// "bytes"
	"fmt"
	"net/http"
	// "strings"
	log "github.com/Sirupsen/logrus"
	"github.com/auroralaboratories/corona-api/modules"
	"github.com/auroralaboratories/corona-api/util"
	"github.com/codegangsta/cli"
	"github.com/julienschmidt/httprouter"
	"github.com/shutterstock/go-stockutil/stringutil"
)

type CommandModule struct {
	modules.BaseModule

	Commands map[string]Command
}

func RegisterSubcommands() []cli.Command {
	return []cli.Command{}
}

func (self *CommandModule) LoadRoutes(router *httprouter.Router) error {
	router.GET(`/api/commands/list`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		util.Respond(w, http.StatusOK, self.Commands, nil)
	})

	router.PUT(`/api/commands/run/:name`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		name := params.ByName(`name`)

		if command, ok := self.Commands[name]; ok {
			if results, err := command.Execute(); err == nil {
				util.Respond(w, http.StatusOK, results, nil)
			} else {
				util.Respond(w, http.StatusOK, results, nil)
			}
		} else {
			util.Respond(w, http.StatusNotFound, nil, fmt.Errorf("Cannot find a command called '%s'", name))
		}
	})

	return nil
}

func (self *CommandModule) Init() error {
	self.Commands = make(map[string]Command)

	if commands, err := GenerateCommands(`cmd`, self.GetConfigRoot()); err == nil {
		for k, v := range commands {
			self.Commands[k] = v
		}
	} else {
		return err
	}

	return nil
}

func GenerateCommands(prefix string, config map[string]interface{}) (map[string]Command, error) {
	commands := make(map[string]Command)

	if cmdInterface, ok := config[`commands`]; ok {
		switch cmdInterface.(type) {
		case map[string]interface{}:
			for key, commandConfigI := range cmdInterface.(map[string]interface{}) {
				key = prefix + `:` + key

				log.Infof("CommandModule: initializing command '%s'", key)

				switch commandConfigI.(type) {
				case map[string]interface{}:
					commandConfig := commandConfigI.(map[string]interface{})
					command := Command{
						Key: key,
					}

					if v, ok := commandConfig[`shellwrap`]; ok {
						if s, err := stringutil.ToString(v); err == nil {
							command.ShellWrap = s
						}
					}

					if v, ok := commandConfig[`detach`]; ok {
						if s, err := stringutil.ToString(v); err == nil {
							command.Detach = (s == `true`)
						}
					}

					if v, ok := commandConfig[`command`]; ok {
						if s, err := stringutil.ToString(v); err == nil {
							command.CommandLine = s
						}
					}

					if err := command.Init(); err == nil {
						commands[key] = command
					} else {
						return commands, err
					}
				}
			}
		}
	}

	return commands, nil
}
