# Copyright (C) 2022 Satoshi Konno All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

SHELL := bash

#PREFIX?=$(shell pwd)
#GOPATH:=$(shell pwd)
#export GOPATH

PACKAGE_NAME=redis

MODULE_ROOT=github.com/cybergarage/go-redis
SOURCE_DIR=${PACKAGE_NAME}
PACKAGE_ROOT=${MODULE_ROOT}/${PACKAGE_NAME}

SOURCES=\
        ${SOURCE_DIR} \
        ${SOURCE_DIR}/protocol

PACKAGE_ID=${PACKAGE_ROOT}
PACKAGES=\
	${PACKAGE_ID} \
	${PACKAGE_ID}/protocol

BINARY_DIR=examples
BINARY_ROOT=${MODULE_ROOT}/${BINARY_DIR}

BINARIES=\
	${BINARY_ROOT}/go-redisd

.PHONY: version format vet lint clean

all: test

version:
	@pushd ${SOURCE_DIR} && ./version.gen > version.go && popd

format: version
	gofmt -w ${SOURCE_DIR} ${BINARY_DIR}

vet: format
	go vet ${PACKAGE_ROOT}

lint: vet
	golangci-lint run ${SOURCES}

build:
	go build -v ${PACKAGES}

test:
	go test -v -cover -timeout 60s ${PACKAGES}

install:
	go install ${BINARIES}

clean:
	go clean -i ${PACKAGES}
