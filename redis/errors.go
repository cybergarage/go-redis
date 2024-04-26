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
	"fmt"
)

var (
	ErrNotSupported = errors.New("not supported")
	ErrQuit         = errors.New("QUIT")
	ErrSystem       = errors.New("internal system error")
	ErrNotAuthrized = errors.New("not authrized")
	ErrInvalid      = errors.New("invalid")
)

const (
	errorNotSupportedCommand    = "'%s' is %w"
	errorMissingCommandArgument = "%s: missing argument (%s) %w"
	errorUnkownCommandArgument  = "%s: unknown argument (%s)"
	errorInvalidCommandArgument = "%s: %w argument (%s - %s)"
	errorUseOnlyOnce            = "%s may be used only once"
	errorShouldBeGreaterThanInt = "%s should be greater than %d"
)

// NewErrNotSupported returns a new ErrNotSupported.
func NewErrNotSupported(target string) error {
	return fmt.Errorf(errorNotSupportedCommand, target, ErrNotSupported)
}

// NewErrorNotSupportedMessage returns a new ErrNotSupported message.
func NewErrorNotSupportedMessage(cmd string) *Message {
	return NewErrorMessage(NewErrNotSupported(cmd))
}

func newMissingArgumentError(cmd string, arg string, err error) error {
	return fmt.Errorf(errorMissingCommandArgument, cmd, arg, err)
}

func newUnkownArgumentError(cmd string, arg string) error {
	return fmt.Errorf(errorUnkownCommandArgument, cmd, arg)
}

func newInvalidArgumentError(cmd string, arg string, err error) error {
	return fmt.Errorf(errorInvalidCommandArgument, cmd, ErrInvalid, arg, err.Error())
}
