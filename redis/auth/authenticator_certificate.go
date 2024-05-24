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

// CertificateAuthenticator represents an authenticator for TLS certificates.
type CertificateAuthenticator struct {
	commonName string
}

// CertificateAuthenticatorOption represents an authenticator option.
type CertificateAuthenticatorOption = func(*CertificateAuthenticator)

// NewCertificateAuthenticator returns a new certificate authenticator.
func NewCertificateAuthenticatorWith(opts ...CertificateAuthenticatorOption) *CertificateAuthenticator {
	authenticator := &CertificateAuthenticator{
		commonName: "",
	}
	for _, opt := range opts {
		opt(authenticator)
	}

	return authenticator
}

// WithCommonName returns an authenticator option to set the common name.
func WithCommonName(name string) func(*CertificateAuthenticator) {
	return func(conn *CertificateAuthenticator) {
		conn.commonName = name
	}
}

// Authenticate authenticates the specified connection.
func (authenticator *CertificateAuthenticator) Authenticate(conn Conn) (bool, error) {
	conState, ok := conn.TLSConnectionState()
	if !ok {
		return false, nil
	}
	for _, cert := range conState.PeerCertificates {
		if 0 < len(authenticator.commonName) {
			if cert.Subject.CommonName == authenticator.commonName {
				return true, nil
			}
		}
	}
	return false, nil
}
