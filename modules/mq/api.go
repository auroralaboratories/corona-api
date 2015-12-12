package main

import (
    "github.com/gorilla/websocket"
    "code.google.com/p/go-uuid/uuid"
    "net/http"
)

var websocketUpgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin:     func(r *http.Request) bool {
        return true
    },
}

func (self *CoronaAPI) WebsocketClientConnect(w http.ResponseWriter, r *http.Request){
    bus := self.Plugin("Bus").(*BusPlugin)
    ws, err := websocketUpgrader.Upgrade(w, r, nil)

    if err != nil {
        logger.Debugf("ERR: %s", err)
        //http.Error(w, err.Error(), 500)
        return
    }else{
        conn := &BusConnection{
            ID:         uuid.New(),
            Connection: ws,
            Plugin:     bus,
            SendBuffer: make(chan *BusMessage, 256),
        }

        bus.Register <- conn
        go conn.Writer()
        conn.Reader()
    }

    return
}

