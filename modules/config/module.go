package main

import (
	"code.google.com/p/leveldb-go/leveldb"
	"code.google.com/p/leveldb-go/leveldb/db"
)

type ConfigPlugin struct {
	BasePlugin

	ConnectOptions *db.Options
	Connection     *leveldb.DB
}

func (self *ConfigPlugin) Init() (err error) {
	self.ConnectOptions = &db.Options{}
	conn, err := leveldb.Open(self.GetConfigOr("plugins.config.db.path", "/tmp/corona").(string), self.ConnectOptions)

	if err != nil {
		return
	}

	self.Connection = conn
	return
}

// func (self *ConfigPlugin) GetKey(namespace string, key string, fallback interface{}) (err error) {
//     //self.Connection.Get()
// }
