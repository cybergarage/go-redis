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
	"github.com/cybergarage/go-authenticator/auth"
	"github.com/cybergarage/go-tracing/tracer"
)

type Server interface {
	auth.Manager
	Config

	// SetTracer sets a tracer.
	SetTracer(tracer.Tracer)

	// Config returns the server configuration.
	Config() Config

	// SetCommandHandler sets a user handler to handle user commands.
	SetCommandHandler(handler UserCommandHandler)
	// RegisterExexutor sets a command executor.
	RegisterExexutor(cmd string, executor Executor)

	Start() error
	Stop() error
	Restart() error
}
