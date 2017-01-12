package backends

import (
	"github.com/auroralaboratories/corona-api/modules/soundctl/types"
)

type BaseOutput struct {
	types.IOutput `json:"-"`
	Name          string                 `json:"name"`
	Backend       types.IBackend         `json:"-"`
	Properties    map[string]interface{} `json:"properties"`
}

func (self *BaseOutput) Initialize(backend types.IBackend) error {
	self.Backend = backend
	self.Properties = make(map[string]interface{})
	return nil
}

func (self *BaseOutput) GetName() string {
	return self.Name
}

func (self *BaseOutput) SetName(name string) {
	self.Name = name
}

func (self *BaseOutput) GetProperty(key string, fallback interface{}) interface{} {
	if v, ok := self.Properties[key]; ok {
		return v
	} else {
		return fallback
	}
}

func (self *BaseOutput) SetProperty(key string, value interface{}) {
	self.Properties[key] = value
}
