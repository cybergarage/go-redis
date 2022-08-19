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
	"errors"
	"strings"
	"time"

	"github.com/cybergarage/go-redis/redis/proto"
)

func nextHashArgument(cmd string, args Arguments) (string, error) {
	hash, err := args.NextString()
	if err != nil {
		return "", newMissingArgumentError(cmd, "hash", err)
	}
	return hash, nil
}

func nextKeyArgument(cmd string, args Arguments) (string, error) {
	key, err := args.NextString()
	if err != nil {
		return "", newMissingArgumentError(cmd, "key", err)
	}
	return key, nil
}

func nextKeysArguments(cmd string, args Arguments) ([]string, error) {
	var key string
	var err error
	keys := []string{}
	key, err = args.NextString()
	for err == nil {
		keys = append(keys, key)
		key, err = args.NextString()
	}
	if !errors.Is(err, proto.ErrEOM) {
		return nil, err
	}
	return keys, nil
}

func nextIntegerArgument(cmd string, key string, args Arguments) (int, error) {
	val, err := args.NextInteger()
	if err != nil {
		return 0, newMissingArgumentError(cmd, key, err)
	}
	return val, nil
}

func nextSetArguments(cmd string, args Arguments) (string, string, error) {
	key, err := args.NextString()
	if err != nil {
		return "", "", newMissingArgumentError(cmd, "key", err)
	}
	val, err := args.NextString()
	if err != nil {
		return "", "", newMissingArgumentError(cmd, "value", err)
	}
	return key, val, err
}

func nextMSetArguments(cmd string, args Arguments) (map[string]string, error) {
	var key, val string
	var err error
	dir := map[string]string{}
	key, err = args.NextString()
	for err == nil {
		val, err = args.NextString()
		if err != nil {
			newMissingArgumentError(cmd, key, err)
		}
		dir[key] = val
		key, err = args.NextString()
	}
	if !errors.Is(err, proto.ErrEOM) {
		return nil, err
	}
	return dir, nil
}

func nextExpireArgument(cmd string, ttl time.Time, args Arguments) (ExpireOption, error) {
	opt := ExpireOption{
		Time: ttl,
		NX:   false,
		XX:   false,
		GT:   false,
		LT:   false,
	}
	var err error
	arg, err := args.NextString()
	if err == nil {
		switch strings.ToUpper(arg) {
		case "NX":
			opt.NX = true
		case "XX":
			opt.XX = true
		case "GT":
			opt.GT = true
		case "LT":
			opt.LT = true
		default:
			return opt, newUnkownArgumentError(cmd, arg)
		}
	}
	if !errors.Is(err, proto.ErrEOM) {
		return opt, err
	}
	return opt, nil
}
