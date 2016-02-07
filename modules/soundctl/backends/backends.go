package backends

import (
    "github.com/auroralaboratories/corona-api/modules/soundctl/types"
)

type BaseBackend struct {
    types.IBackend               `json:"-"`

    Name       string            `json:"name"`
    Outputs    []types.IOutput   `json:"outputs"`
    Properties map[string]string `json:"properties"`
}

func (self *BaseBackend) Initialize() error {
    self.Outputs    = make([]types.IOutput, 0)
    self.Properties = make(map[string]string)
    return nil
}


func (self *BaseBackend) GetName() string {
    return self.Name
}


func (self *BaseBackend) SetName(name string) {
    self.Name = name
}


func (self *BaseBackend) GetProperty(key string, fallback string) string {
    if v, ok := self.Properties[key]; ok {
        return v
    }else{
        return fallback
    }
}


func (self *BaseBackend) SetProperty(key string, value string) {
    self.Properties[key] = value
}

func (self *BaseBackend) GetOutputByName(name string) (types.IOutput, bool) {
    for _, output := range self.Outputs {
        if output.GetName() == name {
            return output, true
        }
    }

    return nil, false
}

func (self *BaseBackend) GetOutputsByProperty(key string, value string) []types.IOutput {
    rv := make([]types.IOutput, 0)

    for _, output := range self.Outputs {
        if propValue := output.GetProperty(key, ``); propValue != `` && propValue == value {
            rv = append(rv, output)
        }
    }

    return rv
}


func (self *BaseBackend) AddOutput(output types.IOutput) error {
    self.Outputs = append(self.Outputs, output)
    return nil
}


func (self *BaseBackend) GetOutputs() []types.IOutput {
    return self.Outputs
}