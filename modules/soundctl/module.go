package soundctl

import (
	"fmt"
	"net/http"
	// "strings"
	"github.com/auroralaboratories/corona-api/modules"
	"github.com/auroralaboratories/corona-api/modules/soundctl/backends/pulseaudio"
	"github.com/auroralaboratories/corona-api/modules/soundctl/types"
	"github.com/auroralaboratories/corona-api/util"
	"github.com/codegangsta/cli"
	"github.com/julienschmidt/httprouter"
	"github.com/shutterstock/go-stockutil/stringutil"
)

type SoundctlModule struct {
	modules.BaseModule

	Backends map[string]types.IBackend
}

func RegisterSubcommands() []cli.Command {
	return []cli.Command{}
}

func (self *SoundctlModule) LoadRoutes(router *httprouter.Router) error {
	router.GET(`/api/soundctl/backends`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		keys := make([]string, 0)

		for key, _ := range self.Backends {
			keys = append(keys, key)
		}

		util.Respond(w, http.StatusOK, keys, nil)
	})

	router.GET(`/api/soundctl/backends/:name`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		name := params.ByName(`name`)

		if backend, ok := self.Backends[name]; ok {
			util.Respond(w, http.StatusOK, backend, nil)
		} else {
			util.Respond(w, http.StatusNotFound, nil, fmt.Errorf("Unable to locate backend '%s'", name))
		}
	})

	router.GET(`/api/soundctl/backends/:name/outputs/:output`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		if output, err := self.getNamedOutput(params.ByName(`name`), params.ByName(`output`)); err == nil {
			util.Respond(w, http.StatusOK, output, nil)
		} else {
			util.Respond(w, http.StatusNotFound, nil, err)
		}
	})

	router.GET(`/api/soundctl/backends/:name/outputs-property/:key/:value`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		name := params.ByName(`name`)
		propName := params.ByName(`key`)
		propValue := params.ByName(`value`)

		if backend, ok := self.Backends[name]; ok {
			if outputs := backend.GetOutputsByProperty(propName, propValue); len(outputs) > 0 {
				util.Respond(w, http.StatusOK, outputs, nil)
			} else {
				util.Respond(w, http.StatusNotFound, nil, fmt.Errorf("Unable to locate any outputs with property %s=%s on backend '%s'", propName, propValue, name))
			}
		} else {
			util.Respond(w, http.StatusNotFound, nil, fmt.Errorf("Unable to locate backend '%s'", name))
		}
	})

	router.PUT(`/api/soundctl/backends/:name/outputs/:output/set-default`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	})

	router.PUT(`/api/soundctl/backends/:name/outputs/:output/mute`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		if output, err := self.getNamedOutput(params.ByName(`name`), params.ByName(`output`)); err == nil {
			if err := output.Mute(); err == nil {
				util.Respond(w, http.StatusNoContent, nil, nil)
			} else {
				util.Respond(w, http.StatusInternalServerError, nil, err)
			}
		} else {
			util.Respond(w, http.StatusNotFound, nil, err)
		}
	})

	router.PUT(`/api/soundctl/backends/:name/outputs/:output/unmute`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		if output, err := self.getNamedOutput(params.ByName(`name`), params.ByName(`output`)); err == nil {
			if err := output.Unmute(); err == nil {
				util.Respond(w, http.StatusNoContent, nil, nil)
			} else {
				util.Respond(w, http.StatusInternalServerError, nil, err)
			}
		} else {
			util.Respond(w, http.StatusNotFound, nil, err)
		}
	})

	router.PUT(`/api/soundctl/backends/:name/outputs/:output/toggle`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		if output, err := self.getNamedOutput(params.ByName(`name`), params.ByName(`output`)); err == nil {
			if err := output.ToggleMute(); err == nil {
				util.Respond(w, http.StatusNoContent, nil, nil)
			} else {
				util.Respond(w, http.StatusInternalServerError, nil, err)
			}
		} else {
			util.Respond(w, http.StatusNotFound, nil, err)
		}
	})

	router.PUT(`/api/soundctl/backends/:name/outputs/:output/volume/:factor`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		if output, err := self.getNamedOutput(params.ByName(`name`), params.ByName(`output`)); err == nil {
			if v, err := stringutil.ConvertToFloat(params.ByName(`factor`)); err == nil {
				if err := output.SetVolume(v); err == nil {
					util.Respond(w, http.StatusNoContent, nil, nil)
				} else {
					util.Respond(w, http.StatusInternalServerError, nil, err)
				}
			} else {
				util.Respond(w, http.StatusBadRequest, nil, err)
			}
		} else {
			util.Respond(w, http.StatusNotFound, nil, err)
		}
	})

	router.PUT(`/api/soundctl/backends/:name/outputs/:output/volume-up/:factor`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		if output, err := self.getNamedOutput(params.ByName(`name`), params.ByName(`output`)); err == nil {
			if v, err := stringutil.ConvertToFloat(params.ByName(`factor`)); err == nil {
				if err := output.IncreaseVolume(v); err == nil {
					util.Respond(w, http.StatusNoContent, nil, nil)
				} else {
					util.Respond(w, http.StatusInternalServerError, nil, err)
				}
			} else {
				util.Respond(w, http.StatusBadRequest, nil, err)
			}
		} else {
			util.Respond(w, http.StatusNotFound, nil, err)
		}
	})

	router.PUT(`/api/soundctl/backends/:name/outputs/:output/volume-down/:factor`, func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		if output, err := self.getNamedOutput(params.ByName(`name`), params.ByName(`output`)); err == nil {
			if v, err := stringutil.ConvertToFloat(params.ByName(`factor`)); err == nil {
				if err := output.DecreaseVolume(v); err == nil {
					util.Respond(w, http.StatusNoContent, nil, nil)
				} else {
					util.Respond(w, http.StatusInternalServerError, nil, err)
				}
			} else {
				util.Respond(w, http.StatusBadRequest, nil, err)
			}
		} else {
			util.Respond(w, http.StatusNotFound, nil, err)
		}
	})

	return nil
}

func (self *SoundctlModule) getNamedOutput(backendName string, outputName string) (types.IOutput, error) {
	if backend, ok := self.Backends[backendName]; ok {
		outputs := backend.GetOutputs()

		if outputName == `current` {
			if output, err := backend.GetCurrentOutput(); err == nil {
				return output, nil
			} else {
				return output, err
			}
		} else if i, err := stringutil.ConvertToInteger(outputName); err == nil && int(i) < len(outputs) {
			return outputs[i], nil
		} else {
			if output, ok := backend.GetOutputByName(outputName); ok {
				return output, nil
			}
		}

		return nil, fmt.Errorf("Unable to locate output '%s' on backend '%s'", outputName, backendName)
	} else {
		return nil, fmt.Errorf("Unable to locate backend '%s'", backendName)
	}
}

func (self *SoundctlModule) PopulateBackends() error {
	self.Backends = make(map[string]types.IBackend)

	self.Backends[`default`] = pulseaudio.New()

	for name, backend := range self.Backends {
		if err := backend.Initialize(); err == nil {
			if err := backend.Refresh(); err != nil {
				return fmt.Errorf("Failed to refresh soundctl backend '%s': %+v", name, err)
			}
		} else {
			return fmt.Errorf("Failed to initialize soundctl backend '%s': %+v", name, err)
		}
	}

	return nil
}

func (self *SoundctlModule) Init() (err error) {
	if err := self.PopulateBackends(); err != nil {
		return err
	}

	return nil
}
