package soundctl

import (
	"fmt"
	"net/http"
	// "strings"
	"github.com/auroralaboratories/corona-api/modules"
	"github.com/auroralaboratories/corona-api/modules/soundctl/backends/pulseaudio"
	"github.com/auroralaboratories/corona-api/modules/soundctl/types"
	"github.com/auroralaboratories/corona-api/util"
	"github.com/ghetzel/cli"
	"github.com/husobee/vestigo"
	"github.com/shutterstock/go-stockutil/stringutil"
)

type SoundctlModule struct {
	modules.BaseModule
	Backends map[string]types.IBackend
}

func RegisterSubcommands() []cli.Command {
	return []cli.Command{}
}

func (self *SoundctlModule) LoadRoutes(router *vestigo.Router) error {
	router.Get(`/api/soundctl/backends`, func(w http.ResponseWriter, req *http.Request) {
		keys := make([]string, 0)

		for key := range self.Backends {
			keys = append(keys, key)
		}

		util.Respond(w, http.StatusOK, keys, nil)
	})

	router.Get(`/api/soundctl/backends/:backend`, func(w http.ResponseWriter, req *http.Request) {
		backendName := vestigo.Param(req, `backend`)

		if backend, ok := self.Backends[backendName]; ok {
			util.Respond(w, http.StatusOK, backend, nil)
		} else {
			util.Respond(w, http.StatusNotFound, nil, fmt.Errorf("Unable to locate backend '%s'", backendName))
		}
	})

	router.Get(`/api/soundctl/backends/:backend/outputs/:output`, func(w http.ResponseWriter, req *http.Request) {
		backendName := vestigo.Param(req, `backend`)
		outputName := vestigo.Param(req, `output`)

		if output, err := self.getNamedOutput(backendName, outputName); err == nil {
			util.Respond(w, http.StatusOK, output, nil)
		} else {
			util.Respond(w, http.StatusNotFound, nil, err)
		}
	})

	router.Get(`/api/soundctl/backends/:backend/outputs-property/:key/:value`, func(w http.ResponseWriter, req *http.Request) {
		backendName := vestigo.Param(req, `backend`)
		propName := vestigo.Param(req, `key`)
		propValue := vestigo.Param(req, `value`)

		if backend, ok := self.Backends[backendName]; ok {
			if outputs := backend.GetOutputsByProperty(propName, propValue); len(outputs) > 0 {
				util.Respond(w, http.StatusOK, outputs, nil)
			} else {
				util.Respond(w, http.StatusNotFound, nil, fmt.Errorf("Unable to locate any outputs with property %s=%s on backend '%s'", propName, propValue, backendName))
			}
		} else {
			util.Respond(w, http.StatusNotFound, nil, fmt.Errorf("Unable to locate backend '%s'", backendName))
		}
	})

	router.Put(`/api/soundctl/backends/:backend/outputs/:output/set-default`, func(w http.ResponseWriter, req *http.Request) {
	})

	router.Put(`/api/soundctl/backends/:backend/outputs/:output/mute`, func(w http.ResponseWriter, req *http.Request) {
		backendName := vestigo.Param(req, `backend`)
		outputName := vestigo.Param(req, `output`)

		if output, err := self.getNamedOutput(backendName, outputName); err == nil {
			if err := output.Mute(); err == nil {
				defer self.refreshBackend(backendName)
				util.Respond(w, http.StatusNoContent, nil, nil)
			} else {
				util.Respond(w, http.StatusInternalServerError, nil, err)
			}
		} else {
			util.Respond(w, http.StatusNotFound, nil, err)
		}
	})

	router.Put(`/api/soundctl/backends/:backend/outputs/:output/unmute`, func(w http.ResponseWriter, req *http.Request) {
		backendName := vestigo.Param(req, `backend`)
		outputName := vestigo.Param(req, `output`)

		if output, err := self.getNamedOutput(backendName, outputName); err == nil {
			if err := output.Unmute(); err == nil {
				defer self.refreshBackend(backendName)
				util.Respond(w, http.StatusNoContent, nil, nil)
			} else {
				util.Respond(w, http.StatusInternalServerError, nil, err)
			}
		} else {
			util.Respond(w, http.StatusNotFound, nil, err)
		}
	})

	router.Put(`/api/soundctl/backends/:backend/outputs/:output/toggle`, func(w http.ResponseWriter, req *http.Request) {
		backendName := vestigo.Param(req, `backend`)
		outputName := vestigo.Param(req, `output`)

		if output, err := self.getNamedOutput(backendName, outputName); err == nil {
			if err := output.ToggleMute(); err == nil {
				defer self.refreshBackend(backendName)
				util.Respond(w, http.StatusNoContent, nil, nil)
			} else {
				util.Respond(w, http.StatusInternalServerError, nil, err)
			}
		} else {
			util.Respond(w, http.StatusNotFound, nil, err)
		}
	})

	router.Put(`/api/soundctl/backends/:backend/outputs/:output/volume/:factor`, func(w http.ResponseWriter, req *http.Request) {
		backendName := vestigo.Param(req, `backend`)
		outputName := vestigo.Param(req, `output`)
		factor := vestigo.Param(req, `factor`)

		if output, err := self.getNamedOutput(backendName, outputName); err == nil {
			if v, err := stringutil.ConvertToFloat(factor); err == nil {
				if err := output.SetVolume(v); err == nil {
					defer self.refreshBackend(backendName)
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

	router.Put(`/api/soundctl/backends/:backend/outputs/:output/volume-up/:factor`, func(w http.ResponseWriter, req *http.Request) {
		backendName := vestigo.Param(req, `backend`)
		outputName := vestigo.Param(req, `output`)
		factor := vestigo.Param(req, `factor`)

		if output, err := self.getNamedOutput(backendName, outputName); err == nil {
			if v, err := stringutil.ConvertToFloat(factor); err == nil {
				if err := output.IncreaseVolume(v); err == nil {
					defer self.refreshBackend(backendName)
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

	router.Put(`/api/soundctl/backends/:backend/outputs/:output/volume-down/:factor`, func(w http.ResponseWriter, req *http.Request) {
		backendName := vestigo.Param(req, `backend`)
		outputName := vestigo.Param(req, `output`)
		factor := vestigo.Param(req, `factor`)

		if output, err := self.getNamedOutput(backendName, outputName); err == nil {
			if v, err := stringutil.ConvertToFloat(factor); err == nil {
				if err := output.DecreaseVolume(v); err == nil {
					defer self.refreshBackend(backendName)
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

func (self *SoundctlModule) refreshBackend(backendName string) error {
	if backend, ok := self.Backends[backendName]; ok {
		return backend.Refresh()
	} else {
		return fmt.Errorf("Unable to locate backend '%s'", backendName)
	}
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

	// TODO: make this configurable
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
