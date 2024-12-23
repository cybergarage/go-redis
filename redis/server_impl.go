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
	"crypto/tls"
	"errors"
	"io"
	"net"
	"strconv"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-redis/redis/auth"
	"github.com/cybergarage/go-redis/redis/proto"
	"github.com/cybergarage/go-tracing/tracer"
)

type server struct {
	*ServerConfig
	*auth.AuthManager
	*ConnManager
	tracer.Tracer
	Addr                 string
	portListener         net.Listener
	tlsPortListener      net.Listener
	tlsConfig            *tls.Config
	authCommandHandler   AuthCommandHandler
	systemCommandHandler SystemCommandHandler
	userCommandHandler   UserCommandHandler
	commandExecutors     Executors
}

// NewServer returns a new server instance.
func NewServer() Server {
	server := &server{
		ServerConfig:         NewDefaultServerConfig(),
		AuthManager:          auth.NewAuthManager(),
		ConnManager:          NewConnManager(),
		Tracer:               tracer.NullTracer,
		Addr:                 "",
		portListener:         nil,
		tlsPortListener:      nil,
		tlsConfig:            nil,
		authCommandHandler:   nil,
		systemCommandHandler: nil,
		userCommandHandler:   nil,
		commandExecutors:     Executors{},
	}
	server.SetPort(DefaultPort)
	server.registerCoreExecutors()
	server.registerSugarExecutors()
	server.systemCommandHandler = server
	server.authCommandHandler = server
	return server
}

// SetTracer sets a tracing tracer.
func (server *server) SetTracer(t tracer.Tracer) {
	server.Tracer = t
}

// Config returns the server configuration.
func (server *server) Config() Config {
	return server.ServerConfig
}

// SetAuthCommandHandler sets a auth handler to handle auth commands.
func (server *server) SetAuthCommandHandler(handler AuthCommandHandler) {
	server.authCommandHandler = handler
}

// SetCommandHandler sets a user handler to handle user commands.
func (server *server) SetCommandHandler(handler UserCommandHandler) {
	server.userCommandHandler = handler
}

// RegisterExexutor sets a command executor.
func (server *server) RegisterExexutor(cmd string, executor Executor) {
	server.commandExecutors[cmd] = executor
}

// Start starts the server.
func (server *server) Start() error {
	password, requirePass := server.ConfigRequirePass()
	if requirePass {
		if !server.HasClearTextPasswordAuthenticator("", password) {
			server.AddAuthenticator(auth.NewClearTextPasswordAuthenticatorWith("", password))
		}
	}

	err := server.ConnManager.Start()
	if err != nil {
		return err
	}

	err = server.open()
	if err != nil {
		return err
	}

	if server.IsPortEnabled() {
		go server.serve()
	}

	if server.IsTLSPortEnabled() {
		go server.tlsServe()
	}

	return nil
}

// Stop stops the server.
func (server *server) Stop() error {
	if err := server.ConnManager.Stop(); err != nil {
		return err
	}

	if err := server.close(); err != nil {
		return err
	}

	if server.IsPortEnabled() {
		addr := net.JoinHostPort(server.Addr, strconv.Itoa(server.ConfigPort()))
		log.Infof("%s/%s (%s) terminated", PackageName, Version, addr)
	}

	if server.IsTLSPortEnabled() {
		addr := net.JoinHostPort(server.Addr, strconv.Itoa(server.ConfigTLSPort()))
		log.Infof("%s/%s (%s) terminated", PackageName, Version, addr)
	}

	return nil
}

// Restart restarts the server.
func (server *server) Restart() error {
	if err := server.Stop(); err != nil {
		return err
	}
	return server.Start()
}

// open opens a listen socket.
func (server *server) open() error {
	var err error

	if server.IsPortEnabled() {
		addr := net.JoinHostPort(server.Addr, strconv.Itoa(server.ConfigPort()))
		server.portListener, err = net.Listen("tcp", addr)
		if err != nil {
			return err
		}
		log.Infof("%s/%s (%s) started", PackageName, Version, addr)
	}

	if server.IsTLSPortEnabled() {
		tlsConfig, ok := server.ConfigTLSConfig()
		if ok {
			server.tlsConfig = tlsConfig
		} else {
			tlsConfig, err := NewTLSConfigFrom(server.ServerConfig)
			if err != nil {
				return err
			}
			server.tlsConfig = tlsConfig
		}
		addr := net.JoinHostPort(server.Addr, strconv.Itoa(server.ConfigTLSPort()))
		server.tlsPortListener, err = net.Listen("tcp", addr)
		if err != nil {
			return err
		}
		log.Infof("%s/%s (%s) started", PackageName, Version, addr)
	}

	return nil
}

// close closes a listening socket.
func (server *server) close() error {
	if server.portListener != nil {
		err := server.portListener.Close()
		if err != nil {
			return err
		}
		server.portListener = nil
	}

	if server.tlsPortListener != nil {
		err := server.tlsPortListener.Close()
		if err != nil {
			return err
		}
		server.tlsPortListener = nil
	}

	return nil
}

// serve handles client connections.
func (server *server) serve() error {
	defer server.close()

	l := server.portListener
	for {
		if l == nil {
			break
		}
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		go server.receive(conn, nil)
	}

	return nil
}

// tlsServe handles client connections with TLS.
func (server *server) tlsServe() error {
	defer server.close()
	l := server.tlsPortListener
	for {
		if l == nil {
			break
		}
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		tlsConn := tls.Server(conn, server.tlsConfig)
		if err := tlsConn.Handshake(); err != nil {
			return err
		}
		tlsState := tlsConn.ConnectionState()

		go server.receive(tlsConn, &tlsState)
	}

	return nil
}

// receive handles a client connection.
func (server *server) receive(conn net.Conn, tlsState *tls.ConnectionState) error {
	_, isPasswdRequired := server.ConfigRequirePass()

	handlerConn := newConnWith(conn, tlsState)
	defer func() {
		handlerConn.Close()
	}()

	handlerConn.SetAuthrized(!isPasswdRequired)
	if tlsState != nil {
		ok, err := server.Authenticate(handlerConn)
		if !ok {
			err = errors.New("invalid client certificates")
		}
		if err != nil {
			log.Error(err)
			return errors.Join(err, handlerConn.Close())
		}
	}

	server.AddConn(handlerConn)
	defer func() {
		server.RemoveConn(handlerConn)
	}()

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
			return nil
		}
		span.Span().Finish()
	}

	return nil
}

// handleMessage handles a client message.
func (server *server) handleMessage(conn *Conn, msg *proto.Message) (*Message, error) {
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
func (server *server) responseMessage(conn io.Writer, msg *Message) error {
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
func (server *server) handleArrayMessage(conn *Conn, arrayMsg *proto.Array) (*Message, error) {
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
