package session

import (
    "github.com/BurntSushi/xgbutil"
    "github.com/codegangsta/cli"
    "github.com/auroralaboratories/corona-api/modules"
)

type SessionModule struct {
    modules.BaseModule

    X *xgbutil.XUtil
}

func RegisterSubcommands() []cli.Command {
    return []cli.Command{}
}

func (self *SessionModule) Init() (err error) {
    self.X, err = xgbutil.NewConn()
    return
}
