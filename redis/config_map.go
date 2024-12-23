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
	"strconv"
	"strings"
)

const (
	ConfigSep = " "
)

// configMap represents a server configuration.
type configMap struct {
	params map[string]string
}

// newConfig returns a new configuration.
func newConfig() *configMap {
	return &configMap{
		params: map[string]string{},
	}
}

// SetConfig sets a specified parameter.
func (cfg *configMap) SetConfig(key string, params string) {
	cfg.params[key] = params
}

// AppendConfig appends a specified parameter.
func (cfg *configMap) AppendConfig(key string, params string) {
	currParams, ok := cfg.params[key]
	if !ok {
		cfg.params[key] = params
		return
	}
	cfg.params[key] = strings.Join([]string{currParams, params}, ConfigSep)
}

// ConfigString return the specified parameter.
func (cfg *configMap) ConfigString(key string) (string, bool) {
	params, ok := cfg.params[key]
	return params, ok
}

// ConfigInteger returns the specified parameter as an integer.
func (cfg *configMap) ConfigInteger(key string) (int, bool) {
	params, ok := cfg.params[key]
	if !ok {
		return 0, false
	}
	v, err := strconv.Atoi(params)
	if err != nil {
		return 0, false
	}
	return v, true
}

// RemoveConfig removes the specified parameter.
func (cfg *configMap) RemoveConfig(key string) {
	delete(cfg.params, key)
}
