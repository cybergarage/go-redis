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

package certs

import (
	"testing"

	"github.com/cybergarage/go-redis/redis"
)

const (
	certFile   = "./cert.pem"
	keyFile    = "./key.pem"
	caCertFile = "./root_cert.pem"
)

func TestCerts(t *testing.T) {
	conf := redis.NewDefaultServerConfig()
	conf.SetTLSCertFile(certFile)
	conf.SetTLSKeyFile(keyFile)
	conf.SetTLSCaCertFile(caCertFile)

	_, err := redis.NewTLSConfigFrom(conf)
	if err != nil {
		t.Error(err)
	}
}
