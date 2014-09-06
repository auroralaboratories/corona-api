package main

import (
    "time"
    "encoding/json"
    "github.com/gorilla/websocket"
)

const (
    // Time allowed to write a message to the peer.
    writeWait = 10 * time.Second

    // Time allowed to read the next pong message from the peer.
    pongWait = 60 * time.Second

    // Send pings to peer with this period. Must be less than pongWait.
    pingPeriod = (pongWait * 9) / 10

    // Maximum message size allowed from peer.
    maxMessageSize = 512
)

// connection is an middleman between the websocket connection and the hub.
type BusConnection struct {
//  a unique id used to identify this connection
    ID         string

//  The websocket connection.
    Connection *websocket.Conn

//  the parent plugin
    Plugin     *BusPlugin

//  Buffered channel of outbound messages.
    SendBuffer chan *BusMessage
}

// readPump pumps messages from the websocket connection to the hub.
func (self *BusConnection) Reader() {
    defer func() {
        self.Plugin.Unregister <- self
        self.Connection.Close()
    }()

    self.Connection.SetReadLimit(maxMessageSize)
    self.Connection.SetReadDeadline(time.Now().Add(pongWait))
    self.Connection.SetPongHandler(func(string) error {
        self.Connection.SetReadDeadline(time.Now().Add(pongWait)); return nil
    })

    for {
        message := &BusMessage{}
        err := self.Connection.ReadJSON(message)
        if err != nil {
            break
        }

        // logger.Debugf("Received message: %s", string(message))
        self.Plugin.Broadcast <- message
    }
}

func (self *BusConnection) PrepareMessage(message *BusMessage) (*BusMessage) {
//  inject connection id
    message.ConnectionID = self.ID

//  inject timestamp
    message.Timestamp = time.Now()

    return message
}

// write writes a message with the given message type and payload.
func (self *BusConnection) Write(payload *BusMessage) error {
    self.Connection.SetWriteDeadline(time.Now().Add(writeWait))
    return self.Connection.WriteJSON(self.PrepareMessage(payload))
}

func (self *BusConnection) Control(mt int, payload *BusMessage) (err error) {
    out, err := json.Marshal(self.PrepareMessage(payload))

    if err != nil {
        logger.Errorf("Unable to send control message: %s", err)
        return
    }

    return self.Connection.WriteControl(mt, out, time.Now().Add(writeWait))
}

// writePump pumps messages from the hub to the websocket connection.
func (self *BusConnection) Writer() {
    ticker := time.NewTicker(pingPeriod)

    defer func() {
        ticker.Stop()
        self.Connection.Close()
    }()

    for {
        select {
        case message, ok := <-self.SendBuffer:
            if !ok {
                self.Control(websocket.CloseMessage, &BusMessage{
                    ConnectionID: self.ID,
                    Type:         "error",
                })
                return
            }

            if err := self.Write(message); err != nil {
                return
            }
        case <-ticker.C:
            if err := self.Control(websocket.PingMessage, &BusMessage{
                    ConnectionID: self.ID,
                    Type:         "ping",
                }); err != nil {
                return
            }
        }
    }
}