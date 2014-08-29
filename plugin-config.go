package main

import (
    "github.com/HouzuoGuo/tiedot/db"
)

type ConfigPlugin struct {
    BasePlugin
    Connection *db.DB
}

func (self *ConfigPlugin) Init() (err error) {
    self.Connection, err = db.OpenDB(self.GetConfigOr("config.db.path", "/tmp/sprinkles.db").(string))
    return
}

