package main

import (
    "github.com/HouzuoGuo/tiedot/db"
)

type ConfigPlugin struct {
    BasePlugin
    Connection *db.DB
}

func (self *ConfigPlugin) Init() (err error) {
    conn, err := db.OpenDB(self.GetConfigOr("plugins.config.db.path", "/tmp/sprinkles").(string))

    if err != nil {
      return
    }

    self.Connection = conn
    return
}

