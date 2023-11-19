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
	"strings"
	"testing"
)

func AuthCommandTest(t *testing.T, server *Server) {
	t.Helper()

	requirePass := "password"
	server.SetRequirePass(requirePass)
	defer server.RemoveRequirePass()

	auths := []struct {
		passwd   string
		expected bool
	}{
		{"", false},
		{requirePass, true},
		{strings.ToUpper(requirePass), false},
	}
	for _, auth := range auths {
		t.Run(auth.passwd, func(t *testing.T) {
			client := NewClient()
			opts := NewClientOptions()
			opts.Password = auth.passwd
			err := client.OpenWith(LocalHost, &opts)
			if auth.expected {
				if err != nil {
					t.Error(err)
				}
				status := client.Ping()
				if status.Err() != nil {
					t.Error(status.Err())
				}
			} else {
				if err == nil {
					t.Errorf("Expected error : %s", auth.passwd)
				}
				status := client.Ping()
				if status.Err() == nil {
					t.Errorf("Expected error : %s", auth.passwd)
				}
			}
			client.Close()
		})
	}
}
