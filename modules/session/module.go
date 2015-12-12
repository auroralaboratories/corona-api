package main

import (
    "github.com/BurntSushi/xgbutil"
)

type SessionPlugin struct {
    BasePlugin
    X          *xgbutil.XUtil
}


func (self *SessionPlugin) Init() (err error) {
    self.X, err = xgbutil.NewConn()
    return
}

