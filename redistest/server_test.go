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

package redistest

import (
	"testing"
)

func TestServer(t *testing.T) {
	server := NewServer()

	err := server.Start()
	if err != nil {
		t.Error(err)
		return
	}

	client := NewClient()
	err = client.Open(LocalHost)
	if err != nil {
		t.Error(err)
		return
	}

	// ctx := context.Background()

	t.Run("Echo", func(t *testing.T) {
		msgs := []string{
			"Hello World!",
		}
		for _, msg := range msgs {
			t.Run(msg, func(t *testing.T) {
				echo := client.Echo(msg)
				if echo.Err() != nil {
					t.Error(echo.Err())
					return
				}
				if echo.Val() != msg {
					t.Errorf("'%s' != '%s'", echo.Val(), msg)
					return
				}
			})
		}
	})

	t.Run("Set", func(t *testing.T) {
		records := []struct {
			key      string
			val      string
			expected string
		}{
			{"key_set", "value0", "value0"},
			{"key_set", "value1", "value1"},
			{"key_set", "value2", "value2"},
		}

		for _, r := range records {
			t.Run(r.key+":"+r.val, func(t *testing.T) {
				err = client.Set(r.key, r.val, 0).Err()
				if err != nil {
					t.Error(err)
				}

				res, err := client.Get(r.key).Result()
				if err != nil {
					t.Error(err)
				}
				if res != r.expected {
					t.Errorf("%s != %s", res, r.expected)
				}
			})
		}
	})

	t.Run("SetNx", func(t *testing.T) {
		records := []struct {
			key      string
			val      string
			expected bool
		}{
			{"key_setnx", "value0", true},
			{"key_setnx", "value1", false},
			{"key_setnx", "value2", false},
		}

		for _, r := range records {
			t.Run(r.key+":"+r.val, func(t *testing.T) {
				res, err := client.SetNX(r.key, r.val, 0).Result()
				if err != nil {
					t.Error(err)
				}
				if res != r.expected {
					t.Errorf("%t != %t", res, r.expected)
				}
			})
		}
	})

	t.Run("GetSet", func(t *testing.T) {
		records := []struct {
			key      string
			val      string
			expected []byte
		}{
			{"key_getset", "value0", nil},
			{"key_getset", "value1", []byte("value0")},
			{"key_getset", "value2", []byte("value1")},
		}

		for _, r := range records {
			t.Run(r.key+":"+r.val, func(t *testing.T) {
				res, err := client.GetSet(r.key, r.val).Result()
				if r.expected == nil {
					if err == nil {
						t.Errorf("%s != %s", res, string(r.expected))
					}
					return
				} else if err != nil {
					t.Error(err)
				}
				if res != string(r.expected) {
					t.Errorf("%s != %s", res, string(r.expected))
				}
			})
		}
	})

	t.Run("MSet", func(t *testing.T) {
		records := []struct {
			hash string
			keys []string
			vals []string
		}{
			{"key_msetmget", []string{"key1", "key2"}, []string{"Hello", "World"}},
		}

		for _, r := range records {
			t.Run(r.hash, func(t *testing.T) {
				args := []string{}
				for n, key := range r.keys {
					args = append(args, key)
					args = append(args, r.vals[n])
				}
				err := client.MSet(args).Err()
				if err != nil {
					t.Error(err)
				}
			})
		}
	})

	t.Run("YCSB", func(t *testing.T) {
		err = ExecYCSB(t)
		if err != nil {
			t.Error(err)
		}
	})

	// // panic: not implemented
	// err = client.Quit().Err()
	// if err != nil {
	// 	t.Error(err)
	// }

	err = client.Close()
	if err != nil {
		t.Error(err)
	}

	err = server.Stop()
	if err != nil {
		t.Error(err)
		return
	}
}
