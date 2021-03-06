package modules

import (
	"fmt"
	"github.com/husobee/vestigo"
)

type IModule interface {
	Init() error
	LoadRoutes(*vestigo.Router) error
	GetConfigRoot() map[string]interface{}
	SetConfig(string, interface{})
	GetConfig(string) (interface{}, bool)
	GetConfigOr(string, interface{}) interface{}
}

type BaseModule struct {
	IModule
	Config map[string]interface{}
}

func (self *BaseModule) Init() error {
	return fmt.Errorf("Unimplemented plugin initializer")
}

func (self *BaseModule) GetConfigRoot() map[string]interface{} {
	return self.Config
}

func (self *BaseModule) SetConfig(name string, value interface{}) {
	self.Config[name] = value
}

func (self *BaseModule) GetConfig(name string) (interface{}, bool) {
	v, ok := self.Config[name]
	return v, ok
}

func (self *BaseModule) GetConfigOr(name string, fallback interface{}) interface{} {
	if value, ok := self.GetConfig(name); ok {
		return value
	} else {
		return fallback
	}
}
