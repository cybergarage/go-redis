/// Copyright (C) 2019 The go-redis Authors. All rights reserved.
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

// Credential represents a credential.
type Credential = auth.Credential

// CredentialOptionFn represents an option function for a credential.
type CredentialOptionFn = auth.CredentialOptionFn

// NewCredential returns a new credential with options.
func NewCredential(opts ...CredentialOptionFn) Credential {
	return auth.NewCredential(opts...)
}

// WithCredentialGroup returns an option to set the group.
func WithCredentialGroup(group string) CredentialOptionFn {
	return auth.WithCredentialGroup(group)
}

// WithCredentialUsername returns an option to set the username.
func WithCredentialUsername(username string) CredentialOptionFn {
	return auth.WithCredentialUsername(username)
}

// WithCredentialPassword returns an option to set the password.
func WithCredentialPassword(password string) CredentialOptionFn {
	return auth.WithCredentialPassword(password)
}
