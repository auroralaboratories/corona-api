package main

import (
    "github.com/BurntSushi/xgbutil"
)

type SessionPlugin struct {
    X    *xgbutil.XUtil
    Name string
}

func (self *SessionPlugin) Init() (err error) {
    self.X, err = xgbutil.NewConn()
    return
}

