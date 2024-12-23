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

	"github.com/cybergarage/go-redis/redis"
	"github.com/cybergarage/go-redis/redis/auth"
)

const (
	serverKey  = "./certs/key.pem"
	serverCert = "./certs/cert.pem"
	rootCert   = "./certs/root_cert.pem"
)

func TestServer(t *testing.T) {
	server := NewServer()

	err := server.Start()
	if err != nil {
		t.Error(err)
		return
	}

	// AuthCommandTest

	AuthCommandTest(t, server)

	// CommandTest

	client := NewClient()
	err = client.Open(LocalHost)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Command", func(t *testing.T) {
		CommandTest(t, client)
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

	// redis-benchmark

	params := []string{"-t get,set -n 10000"}
	t.Run("redis-benchmark", func(t *testing.T) {
		for _, param := range params {
			t.Run(param, func(t *testing.T) {
				ExecRedisBenchmark(t, param)
			})
		}
	})

	// YCSB

	workloads := []string{"workloada", "workloadb"}
	t.Run("YCSB", func(t *testing.T) {
		for _, workload := range workloads {
			t.Run(workload, func(t *testing.T) {
				ExecYCSBWorkload(t, workload)
			})
		}
	})

	err = server.Stop()
	if err != nil {
		t.Error(err)
		return
	}
}

func TestTLSServer(t *testing.T) {
	var err error

	server := NewServer()
	server.SetPort(0)
	server.SetTLSPort(redis.DefaultPort)
	server.SetServerKeyFile(serverKey)
	server.SetServerCertFile(serverCert)
	server.SetRootCertFiles(rootCert)

	ca, err := auth.NewCertificateAuthenticator(auth.WithCertificateAuthenticatorCommonNameRegexp("localhost"))
	if err != nil {
		t.Error(err)
		return
	}
	server.SetCertificateAuthenticator(ca)

	err = server.Start()
	if err != nil {
		t.Error(err)
		return
	}

	// CommandTest

	client := NewClient()
	clientOpts := NewClientOptions()

	tlsConfig, err := server.TLSConfig()
	if err != nil {
		t.Error(err)
		return
	}
	clientOpts.TLSConfig = tlsConfig

	err = client.OpenWith(LocalHost, redis.DefaultPort, &clientOpts)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Command", func(t *testing.T) {
		CommandTest(t, client)
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
