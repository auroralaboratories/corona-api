package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/auroralaboratories/corona-api/modules"
	"github.com/auroralaboratories/corona-api/modules/command"
	"github.com/auroralaboratories/corona-api/modules/session"
	"github.com/auroralaboratories/corona-api/modules/soundctl"
	"github.com/auroralaboratories/corona-api/util"
	"github.com/codegangsta/negroni"
	"github.com/ghetzel/diecast/diecast"
	"github.com/ghodss/yaml"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"

	log "github.com/Sirupsen/logrus"
)

const (
	DEFAULT_CONFIG_PATH  = `corona.yml`
	DEFAULT_MODULES_LIST = `command,session,soundctl`
)

type Configuration struct {
	Modules map[string]map[string]interface{} `json:"modules,omitempty"`
}

type API struct {
	Address     string
	Port        int
	UiDirectory string
	ModulesList []string
	Modules     []modules.IModule
	ConfigPath  string
	Config      Configuration

	router *httprouter.Router
	server *negroni.Negroni
	cors   *cors.Cors
}

func NewApi() *API {
	return &API{
		Address:     `localhost`,
		ConfigPath:  DEFAULT_CONFIG_PATH,
		Modules:     make([]modules.IModule, 0),
		ModulesList: strings.Split(DEFAULT_MODULES_LIST, `,`),
		Port:        25672,
		UiDirectory: `embedded`,
	}
}

func (self *API) LoadConfig() error {
	log.Infof("Loading configuration at %s", self.ConfigPath)

	if data, err := ioutil.ReadFile(self.ConfigPath); err == nil {
		var config Configuration

		if err := yaml.Unmarshal(data, &config); err == nil {
			self.Config = config
		} else {
			return err
		}
	}

	return nil
}

func (self *API) LoadModules() error {
	log.Debugf("Loading modules: %s", strings.Join(self.ModulesList, `, `))

	for _, name := range self.ModulesList {
		var module modules.IModule
		var moduleConfig map[string]interface{}

		if v, ok := self.Config.Modules[name]; ok {
			moduleConfig = v
		} else {
			moduleConfig = make(map[string]interface{})
		}

		switch name {
		case `session`:
			module = &session.SessionModule{
				BaseModule: modules.BaseModule{
					Config: moduleConfig,
				},
			}
		case `command`:
			module = &command.CommandModule{
				BaseModule: modules.BaseModule{
					Config: moduleConfig,
				},
			}
		case `soundctl`:
			module = &soundctl.SoundctlModule{
				BaseModule: modules.BaseModule{
					Config: moduleConfig,
				},
			}
		default:
			log.Fatalf("Unrecognized module name '%s'", name)
		}

		self.Modules = append(self.Modules, module)
	}

	return nil
}

func (self *API) InitModules() error {
	log.Debugf("Initializing all modules")

	for _, module := range self.Modules {
		if err := module.Init(); err != nil {
			return err
		}
	}

	return nil
}

func (self *API) Init() error {
	loadFuncs := [](func() error){
		self.LoadConfig,
		self.LoadModules,
		self.InitModules,
	}

	for _, loadFunc := range loadFuncs {
		if err := loadFunc(); err != nil {
			return err
		}
	}

	return nil
}

func (self *API) Serve() error {
	self.router = httprouter.New()

	if err := self.loadRoutes(); err != nil {
		return err
	}

	go func() {
		dc := diecast.NewServer()
		dc.Address = self.Address
		dc.Port = self.Port + 1
		dc.LogLevel = `debug`

		if err := dc.Initialize(); err == nil {
			dc.Serve()
		}
	}()

	self.cors = cors.New(cors.Options{
		AllowedOrigins: []string{`*`},
		AllowedHeaders: []string{`*`},
	})

	self.server = negroni.New()
	self.server.Use(negroni.NewRecovery())
	self.server.Use(self.cors)
	self.server.UseHandler(self.router)

	self.server.Run(fmt.Sprintf("%s:%d", self.Address, self.Port))

	return nil
}

func (self *API) loadRoutes() error {
	self.router.GET(`/api/status`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		util.Respond(w, http.StatusOK, map[string]interface{}{
			`started_at`: util.StartedAt,
		}, nil)
	})

	for _, module := range self.Modules {
		log.Debugf("Loading routes for %T", module)

		if err := module.LoadRoutes(self.router); err != nil {
			return err
		}
	}

	return nil
}
