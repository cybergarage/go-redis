// Copyright (C) 2022 Satoshi Konno All rights reserved.
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
	"errors"
	"io"
	"net"
	"strconv"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-redis/redis/proto"
	"github.com/cybergarage/go-tracing/tracer"
)

// Server is an instance for Redis protocols.
type Server struct {
	ServerConfig
	tracer.Tracer
	Addr                 string
	tcpListener          net.Listener
	authCommandHandler   AuthCommandHandler
	systemCommandHandler SystemCommandHandler
	userCommandHandler   UserCommandHandler
	commandExecutors     Executors
}

// NewServer returns a new server instance.
func NewServer() *Server {
	server := &Server{
		Tracer:               tracer.NullTracer,
		Addr:                 "",
		tcpListener:          nil,
		authCommandHandler:   nil,
		systemCommandHandler: nil,
		userCommandHandler:   nil,
		commandExecutors:     Executors{},
		ServerConfig:         *NewDefaultServerConfig(),
	}
	server.SetPort(DefaultPort)
	server.registerCoreExecutors()
	server.registerSugarExecutors()
	server.systemCommandHandler = server
	return server
}

// SetTracer sets a tracing tracer.
func (server *Server) SetTracer(t tracer.Tracer) {
	server.Tracer = t
}

// SetAuthCommandHandler sets a auth handler to handle auth commands.
func (server *Server) SetAuthCommandHandler(handler AuthCommandHandler) {
	server.authCommandHandler = handler
}

// SetCommandHandler sets a user handler to handle user commands.
func (server *Server) SetCommandHandler(handler UserCommandHandler) {
	server.userCommandHandler = handler
}

// RegisterExexutor sets a command executor.
func (server *Server) RegisterExexutor(cmd string, executor Executor) {
	server.commandExecutors[cmd] = executor
}

// Start starts the server.
func (server *Server) Start() error {
	err := server.open()
	if err != nil {
		return err
	}

	go server.serve()

	addr := net.JoinHostPort(server.Addr, strconv.Itoa(server.ConfigPort()))
	log.Infof("%s/%s (%s) started", PackageName, Version, addr)

	return nil
}

// Stop stops the server.
func (server *Server) Stop() error {
	if err := server.close(); err != nil {
		return err
	}

	addr := net.JoinHostPort(server.Addr, strconv.Itoa(server.ConfigPort()))
	log.Infof("%s/%s (%s) terminated", PackageName, Version, addr)

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
	addr := net.JoinHostPort(server.Addr, strconv.Itoa(server.ConfigPort()))
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
func (server *Server) receive(conn net.Conn) error {
	defer conn.Close()

	isPasswdRequired, _ := server.ConfigRequirePass()

	handlerConn := newConn()
	handlerConn.SetAuthrized(!isPasswdRequired)

	log.Debugf("%s/%s (%s) accepted", PackageName, Version, conn.RemoteAddr().String())

	parser := proto.NewParserWithReader(conn)

	for {
		span := server.Tracer.StartSpan(PackageName)
		handlerConn.SetSpanContext(span)

		handlerConn.StartSpan("parse")
		reqMsg, parserErr := parser.Next()
		handlerConn.FinishSpan()
		if parserErr != nil {
			span.Span().Finish()
			log.Error(parserErr)
			return parserErr
		}
		if reqMsg == nil {
			span.Span().Finish()
			break
		}

		var resMsg *Message
		var reqErr error

		resMsg, reqErr = server.handleMessage(handlerConn, reqMsg)
		if reqErr != nil {
			if !errors.Is(reqErr, ErrQuit) {
				resMsg = NewErrorMessage(reqErr)
			}
		}

		handlerConn.StartSpan("response")
		resErr := server.responseMessage(conn, resMsg)
		handlerConn.FinishSpan()
		if resErr != nil {
			log.Error(resErr)
		}
		if errors.Is(reqErr, ErrQuit) {
			span.Span().Finish()
			conn.Close()
			return nil
		}
		span.Span().Finish()
	}

	return nil
}

// handleMessage handles a client message.
func (server *Server) handleMessage(conn *Conn, msg *proto.Message) (*Message, error) {
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
		return server.handleArrayMessage(conn, msg)
	case proto.ErrorMessage:
		return nil, nil
	}
	return nil, nil
}

// responseMessage returns the response message to the request connection.
func (server *Server) responseMessage(conn io.Writer, msg *Message) error {
	var bytes []byte
	var err error
	if msg != nil {
		bytes, err = msg.RESPBytes()
	} else {
		bytes, err = NewErrorMessage(ErrSystem).Bytes()
	}
	if err != nil {
		return err
	}
	_, err = conn.Write(bytes)
	return err
}

// handleMessage handles a client message.
func (server *Server) handleArrayMessage(conn *Conn, arrayMsg *proto.Array) (*Message, error) {
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
		return server.handleArrayMessage(conn, nestedArrayMsg)
	}

	cmd, err := firstMsg.String()
	if err != nil {
		return nil, err
	}

	return server.executeCommand(conn, cmd, arrayMsg)
}
