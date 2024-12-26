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
	*serverConfig
	*auth.Manader
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
	credStore            map[string]auth.Credential
}

// NewServer returns a new server instance.
func NewServer() Server {
	server := &server{
		serverConfig:         newDefaultServerConfig(),
		Manader:              auth.NewManader(),
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
		credStore:            make(map[string]auth.Credential),
	}

	server.SetPort(DefaultPort)
	server.registerCoreExecutors()
	server.registerSugarExecutors()
	server.systemCommandHandler = server
	server.authCommandHandler = server
	server.SetCredentialStore(server)

	return server
}

// SetTracer sets a tracing tracer.
func (server *server) SetTracer(t tracer.Tracer) {
	server.Tracer = t
}

// Config returns the server configuration.
func (server *server) Config() Config {
	return server.serverConfig
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
		cred := auth.NewCredential(
			auth.WithCredentialPassword(password),
		)
		server.SetCredential(cred)
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
		addr := net.JoinHostPort(server.Addr, strconv.Itoa(server.Port()))
		log.Infof("%s/%s (%s) terminated", PackageName, Version, addr)
	}

	if server.IsTLSPortEnabled() {
		addr := net.JoinHostPort(server.Addr, strconv.Itoa(server.TLSPort()))
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
		addr := net.JoinHostPort(server.Addr, strconv.Itoa(server.Port()))
		server.portListener, err = net.Listen("tcp", addr)
		if err != nil {
			return err
		}
		log.Infof("%s/%s (%s) started", PackageName, Version, addr)
	}

	if server.IsTLSPortEnabled() {
		tlsConfig, err := server.TLSConfig()
		if err != nil {
			return err
		}
		server.tlsConfig = tlsConfig
		addr := net.JoinHostPort(server.Addr, strconv.Itoa(server.TLSPort()))
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
		ok, err := server.Manager.VerifyCertificate(tlsConn)
		if !ok {
			log.Error(err)
			return err
		}

		go server.receive(tlsConn, tlsConn)
	}

	return nil
}

// receive handles a client connection.
func (server *server) receive(conn net.Conn, tlsConn *tls.Conn) error {
	_, isPasswdRequired := server.ConfigRequirePass()

	handlerConn := newConnWith(conn, tlsConn)
	defer func() {
		handlerConn.Close()
	}()

	handlerConn.SetAuthrized(!isPasswdRequired)
	if tlsConn != nil {
		ok, err := server.VerifyCertificate(tlsConn)
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
