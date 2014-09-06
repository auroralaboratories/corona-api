// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
    "time"
)

type BusPlugin struct {
    BasePlugin

//  Registered connections.
    Connections map[*BusConnection]bool

//  Inbound messages from the connections.
    Broadcast chan *BusMessage

//  Register requests from the connections.
    Register chan *BusConnection

//  Unregister requests from connections.
    Unregister chan *BusConnection
}

const (
    BusMessageTypeConnect    = "connect"
    BusMessageTypeDisconnect = "disconnected"
    BusMessageTypeEvent      = "event"
)

type BusMessage struct {
    ConnectionID        string     `json:"id,omitempty"`
    Timestamp           time.Time  `json:"timestamp,omitempty"`
    Type                string     `json:"type"`
    Message             string     `json:"message,omitempty"`
}

func (self *BusPlugin) Init() (err error) {
    self.Broadcast   = make(chan *BusMessage)
    self.Register    = make(chan *BusConnection)
    self.Unregister  = make(chan *BusConnection)
    self.Connections = make(map[*BusConnection]bool)

    go self.Run()
    return
}

func (self *BusPlugin) Run() {
    logger.Info("Starting Message Bus")

    for {
        select {
    //  on register: flag as active
        case c := <- self.Register:
            self.Connections[c] = true
            c.Write(&BusMessage{
                Type: BusMessageTypeConnect,
            })

            logger.Infof("Message bus client connected: %s", c.ID)

    //  on unregister: cleanup connection
        case c := <- self.Unregister:
            if _, ok := self.Connections[c]; ok {
                c.Write(&BusMessage{
                    Type: BusMessageTypeDisconnect,
                })

                logger.Infof("Message bus client disconnected: %s", c.ID)
                close(c.SendBuffer)
                delete(self.Connections, c)
            }

    //  on broadcast: push message to the send buffer of all active connections
        case m := <- self.Broadcast:
            for c, active := range self.Connections {
                if active {
                    select {
                    case c.SendBuffer <- m:
                    default:
                        close(c.SendBuffer)
                        delete(self.Connections, c)
                    }
                }
            }
        }
    }
}