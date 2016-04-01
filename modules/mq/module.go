// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
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
	BusMessageTypeControl    = "control"
	BusMessageTypeSuccess    = "succeeded"
	BusMessageTypeFailure    = "failed"
	BusMessageTypeDisconnect = "disconnected"
	BusMessageTypeEvent      = "event"
)

const (
	BusMessageVerbTag        = "tag"
	BusMessageVerbActivate   = "activate"
	BusMessageVerbDeactivate = "deactivate"
)

type BusMessage struct {
	ConnectionID string    `json:"id,omitempty"`
	Tags         []string  `json:"tags,omitempty"`
	Global       bool      `json:"global,omitempty"`
	Timestamp    time.Time `json:"timestamp,omitempty"`
	Type         string    `json:"type"`
	Verb         string    `json:"verb,omitempty"`
	Arguments    []string  `json:"arguments,omitempty"`
	Message      string    `json:"message,omitempty"`
}

func (self *BusPlugin) Init() (err error) {
	self.Broadcast = make(chan *BusMessage)
	self.Register = make(chan *BusConnection)
	self.Unregister = make(chan *BusConnection)
	self.Connections = make(map[*BusConnection]bool)

	go self.Run()
	return
}

func (self *BusPlugin) GetConnection(id string) (connection *BusConnection, err error) {
	if id == "" {
		err = errors.New("Cannot get connection for empty ID")
		return
	}

	for conn, _ := range self.Connections {
		if conn.ID == id {
			connection = conn
			return
		}
	}

	err = errors.New("Unable to find connection '" + id + "'")
	return
}

func (self *BusPlugin) HandleControlVerbTag(connection *BusConnection, message *BusMessage, response *BusMessage) {
	connection.FilterTags = make([]string, len(message.Arguments))
	copy(connection.FilterTags, message.Arguments)
	response.Message = fmt.Sprintf("Connection tags set to %s", connection.FilterTags)
}

func (self *BusPlugin) HandleControlVerbActivate(connection *BusConnection, message *BusMessage, response *BusMessage) {
	connection.Active = true
	response.Message = fmt.Sprintf("Connection has been activated and will receive all messages destined for it")
}

func (self *BusPlugin) HandleControlVerbDeactivate(connection *BusConnection, message *BusMessage, response *BusMessage) {
	connection.Active = false
	response.Message = fmt.Sprintf("Connection has been deactivated and will only process control messages")
}

func (self *BusPlugin) HandleControlErrorUnknownVerb(connection *BusConnection, message *BusMessage, response *BusMessage) {
	response.Type = BusMessageTypeFailure
	response.Message = "Unknown verb '" + message.Verb + "'"
	logger.Warnf("Message bus received control message with unknown verb '%s'", message.Verb)
}

func (self *BusPlugin) Run() {
	logger.Info("Starting Message Bus")

	for {
		select {
		//  on register: flag as active
		case c := <-self.Register:
			c.Active = true
			self.Connections[c] = true

			//  client should save this connection ID if it wants to make future control requests
			c.Write(&BusMessage{
				ConnectionID: c.ID,
				Type:         BusMessageTypeConnect,
			})

			logger.Infof("Message bus client connected: %s", c.ID)

			//  on unregister: cleanup connection
		case c := <-self.Unregister:
			if _, ok := self.Connections[c]; ok {
				c.Write(&BusMessage{
					Type: BusMessageTypeDisconnect,
				})

				logger.Infof("Message bus client disconnected: %s", c.ID)
				close(c.SendBuffer)
				delete(self.Connections, c)
			}

			//  on broadcast: push message to the send buffer of all active connections
		case m := <-self.Broadcast:
			switch m.Type {

			//  CONTROL MESSAGES
			//  if the message type is set to 'control', then it is interpreted as a command
			//  designed to modify the current state of the specified connection ID
			//
			//  The following additional fields are required for control messages:
			//    "id"        => determines which connection to modify (clients are only ever shown their own ID as a response at connect time)
			//    "verb"      => a command that determines a course of action to take.  some commands will use "arguments" for command-specific input
			//
			//  Optional:
			//    "arguments" => an array of strings that have context-specific meaning to each verb
			//

			case BusMessageTypeControl:
				//  ID is required to
				if conn, err := self.GetConnection(m.ConnectionID); err == nil {
					//  allocate response object
					response := &BusMessage{
						Type: BusMessageTypeSuccess,
					}

					//  route verb to correct function
					switch m.Verb {
					case BusMessageVerbTag:
						self.HandleControlVerbTag(conn, m, response)
					case BusMessageVerbActivate:
						self.HandleControlVerbActivate(conn, m, response)
					case BusMessageVerbDeactivate:
						self.HandleControlVerbDeactivate(conn, m, response)

						//  anything not handled is an error
					default:
						self.HandleControlErrorUnknownVerb(conn, m, response)
					}

					//  send the response
					conn.Write(response)
				} else {
					logger.Warnf("Message bus could not get connection '%s': %s", m.ConnectionID, err)
				}
			default:
				//  send inbound message to all connections
				for c, _ := range self.Connections {
					skip_send := false

					//  ...but on only send to active connections
					if c.Active {
						//  TAG FILTERS
						//  are a mechanism for allowing clients to select only messages that they are interested in
						//
						//  - Clients register a set of tags that they would like to subscribe to
						//  - Other clients broadcast events with certain tags attached to the messages
						//  - Only clients that have all of the messages tags will receive the events
						//
						//
						//
						//  the Global flag will override this and let the message through regardless of message tags
						//  allow all messages if we don't have any filters setup
						if !m.Global && len(c.FilterTags) > 0 {
							//  return before sending if this connection's FilterTags do not include all of this messages tags
							for _, tag := range m.Tags {
								//  there's a tag in the message that isn't on our connection, don't pass it along
								if !contains(c.FilterTags, tag) {
									logger.Debugf("Connection %s does not include tag %s", c.ID, tag)
									skip_send = true
									break
								}
							}
						}

						//  only send if we're supposed to
						if !skip_send {
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
	}
}
