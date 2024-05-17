# Copyright (C) 2019 The go-redis Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http:#www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

.PHONY: certgen clean

all: certgen

# PostgreSQL: Documentation: 16: 19.9. Secure TCP/IP Connections with SSL
# https://www.postgresql.org/docs/current/ssl-tcp.html#SSL-CERTIFICATE-CREATION

certgen: clean
	openssl req -new -x509 -days 36500 -nodes -text -out cert.pem -keyout key.pem -extensions san -config san.cfg -subj "/CN=localhost"
	chmod og-rwx key.pem
	cp cert.pem root_cert.pem

clean:
	rm -f cert.pem key.pem root_cert.pem