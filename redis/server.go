// Copyright (C) 2022 The go-redis Authors All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package redis

import (
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-redis/redis/proto"
)

// Server is an instance for Redisprotocols.
type Server struct {
	Addr             string
	Port             int
	tcpListener      net.Listener
	systemCmdHandler SystemCommandHandler
	CommandHandler   CommandHandler
}

// NewServer returns a new server instance.
func NewServer() *Server {
	server := &Server{
		Addr:             "",
		Port:             DefaultPort,
		tcpListener:      nil,
		systemCmdHandler: nil,
		CommandHandler:   nil,
	}
	server.systemCmdHandler = server
	return server
}

// Start starts the server.
func (server *Server) Start() error {
	err := server.Stop()
	if err != nil {
		return err
	}

	err = server.open()
	if err != nil {
		return err
	}

	go server.serve()

	return nil
}

// Stop stops the server.
func (server *Server) Stop() error {
	if err := server.close(); err != nil {
		return err
	}
	return nil
}

// Restart restarts the server.
func (server *Server) Restart() error {
	if err := server.Stop(); err != nil {
		return err
	}
	return server.Start()
}

// open opens a listen socket.
func (server *Server) open() error {
	var err error
	addr := net.JoinHostPort(server.Addr, strconv.Itoa(server.Port))
	server.tcpListener, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return nil
}

// close closes a listening socket.
func (server *Server) close() error {
	if server.tcpListener != nil {
		err := server.tcpListener.Close()
		if err != nil {
			return err
		}
	}

	server.tcpListener = nil

	return nil
}

// serve handles client connections.
func (server *Server) serve() error {
	defer server.close()

	l := server.tcpListener
	for {
		if l == nil {
			break
		}
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		go server.receive(conn)
	}

	return nil
}

// receive handles a client connection.
func (server *Server) receive(conn io.ReadWriteCloser) error {
	defer conn.Close()

	parser := proto.NewParserWithReader(conn)
	reqMsg, err := parser.Next()
	for reqMsg != nil {
		if err != nil {
			log.Error(err.Error())
			return err
		}
		var resMsg *Message
		resMsg, err = server.handleMessage(reqMsg)
		if err == nil {
			err = server.responseMessage(conn, resMsg)
		}
		if err != nil {
			log.Error(err.Error())
		}
		reqMsg, err = parser.Next()
	}

	return nil
}

// handleMessage handles a client message.
func (server *Server) handleMessage(msg *proto.Message) (*Message, error) {
	switch msg.Type {
	case proto.StringMessage:
		return nil, nil
	case proto.IntegerMessage:
		return nil, nil
	case proto.BulkMessage:
		return nil, nil
	case proto.ArrayMessage:
		msg, err := msg.Array()
		if err != nil {
			return nil, err
		}
		return server.handleArrayMessage(msg)
	case proto.ErrorMessage:
		return nil, nil
	}
	return nil, nil
}

// responseMessage returns the response message to the request connection.
func (server *Server) responseMessage(conn io.Writer, msg *Message) error {
	bytes, err := msg.RESPBytes()
	if err != nil {
		return err
	}
	_, err = conn.Write(bytes)
	return err
}

// handleMessage handles a client message.
func (server *Server) handleArrayMessage(arrayMsg *proto.Array) (*Message, error) {
	firstMsg, err := arrayMsg.Next()
	if err != nil {
		return nil, err
	}

	// Nested array ?
	if firstMsg.IsArray() {
		nestedArrayMsg, err := firstMsg.Array()
		if err != nil {
			return nil, err
		}
		return server.handleArrayMessage(nestedArrayMsg)
	}

	cmd, err := firstMsg.String()
	if err != nil {
		return nil, err
	}

	args, err := arrayMsg.NextMessages()
	if err != nil {
		return nil, err
	}

	var resMsg *Message

	if strings.ToUpper(cmd) == "PING" {
		if len(args) < 1 {
			resMsg, err = server.systemCmdHandler.Ping("")
		} else {
			var pingMsg string
			pingMsg, err = args[0].String()
			if err != nil {
				return nil, err
			}
			resMsg, err = server.systemCmdHandler.Ping(pingMsg)
		}
	}

	return resMsg, err
}
