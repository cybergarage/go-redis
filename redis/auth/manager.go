// Copyright (C) 2019 The go-redis Authors. All rights reserved.
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

package auth

import (
	"github.com/cybergarage/go-authenticator/auth"
)

// AuthManager represent an authenticator manager.
type AuthManager struct {
	auth.Manager
	authenticators []Authenticator
}

// NewAuthManager returns a new authenticator manager.
func NewAuthManager() *AuthManager {
	manager := &AuthManager{
		Manager:        auth.NewManager(),
		authenticators: make([]Authenticator, 0),
	}
	return manager
}

// AddAuthenticator adds a new authenticator.
func (mgr *AuthManager) AddAuthenticator(authenticator Authenticator) {
	mgr.authenticators = append(mgr.authenticators, authenticator)
}

// ClearAuthenticators clears all authenticators.
func (mgr *AuthManager) ClearAuthenticators() {
	mgr.authenticators = make([]Authenticator, 0)
}

// Authenticate authenticates the connection with the startup message.
func (mgr *AuthManager) Authenticate(conn Conn) (bool, error) {
	if len(mgr.authenticators) == 0 {
		return true, nil
	}
	for _, authenticator := range mgr.authenticators {
		ok, err := authenticator.Authenticate(conn)
		if !ok {
			return false, err
		}
	}
	return true, nil
}

// HasClearTextPasswordAuthenticator returns true if the manager has the clear text password authenticator.
func (mgr *AuthManager) HasClearTextPasswordAuthenticator(username string, password string) bool {
	for _, authenticator := range mgr.authenticators {
		clearTextAuthenticator, ok := authenticator.(*ClearTextPasswordAuthenticator)
		if !ok {
			continue
		}
		if clearTextAuthenticator.username == username && clearTextAuthenticator.password == password {
			return true
		}
	}
	return false
}
