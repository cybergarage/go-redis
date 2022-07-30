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
			{"key", "value", "value"},
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
