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

// Manader represent an authenticator manager.
type Manader struct {
	auth.Manager
}

// NewManader returns a new authenticator manager.
func NewManader() *Manader {
	manager := &Manader{
		Manager: auth.NewManager(),
	}
	return manager
}
