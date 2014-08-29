package main

type IPlugin interface {
    Init()                            error
    SetConfig(string, interface{})
    GetConfig(string)                  (interface{}, bool)
    GetConfigOr(string, interface{})   (interface{})
}

type BasePlugin struct {
    Config map[string]interface{}
}

func (self *BasePlugin) Init() (err error) {
    panic("Unimplemented plugin initializer")
}

func (self *BasePlugin) SetConfig(name string, value interface{}){
    self.Config[name] = value
}

func (self *BasePlugin) GetConfig(name string) (value interface{}, ok bool){
  value, ok = self.Config[name]
  return
}

func (self *BasePlugin) GetConfigOr(name string, fallback interface{}) (value interface{}){
    value, ok := self.GetConfig(name)

    if ok != true {
      value = fallback
    }

    return
}

