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

MODULE_ROOT=github.com/cybergarage/go-redis

PKG_NAME=redis
PKG_VER=$(shell git describe --abbrev=0 --tags)

PKG_ID=${MODULE_ROOT}/${PKG_NAME}
PKG_SRC_DIR=${PKG_NAME}
PKG_SRCS=\
        ${PKG_SRC_DIR} \
        ${PKG_SRC_DIR}/regexp \
        ${PKG_SRC_DIR}/proto
PKGS=\
	${PKG_ID} \
	${PKG_ID}/regexp \
	${PKG_ID}/proto

BIN_DIR=examples
BIN_ID=${MODULE_ROOT}/${BIN_DIR}
BIN_SERVER=go-redisd
BIN_SERVER_DOCKER_TAG=cybergarage/${BIN_SERVER}:${PKG_VER}
BIN_SERVER_ID=${BIN_ID}/${BIN_SERVER}
BIN_SRCS=\
	${BIN_DIR}/${BIN_SERVER} \
	${BIN_DIR}/${BIN_SERVER}/server
BINS=\
	${BIN_SERVER_ID}

TEST_PKG_NAME=${PKG_NAME}test
TEST_PKG_ID=${MODULE_ROOT}/${TEST_PKG_NAME}
TEST_PKG_DIR=${TEST_PKG_NAME}
TEST_PKG_SRCS=\
	${TEST_PKG_DIR}
TEST_PKGS=\
	${TEST_PKG_ID}

.PHONY: version format vet lint clean

all: test

%.md : %.adoc
	asciidoctor -b docbook -a leveloffset=+1 -o - $< | pandoc  --markdown-headings=atx --wrap=preserve -t markdown_strict -f docbook > $@
docs := $(patsubst %.adoc,%.md,$(wildcard *.adoc doc/*.adoc))
doc: $(docs)

version:
	@pushd ${PKG_SRC_DIR} && ./version.gen > version.go && popd

format: version doc
	gofmt -s -w ${PKG_SRC_DIR} ${BIN_DIR} ${TEST_PKG_DIR}

vet: format
	go vet ${PKG_ID} ${TEST_PKG_ID} ${BINS}

lint: vet
	golangci-lint run ${PKG_SRCS} ${BIN_SRCS} ${TEST_PKG_SRCS}

test: lint
	go test -v -cover -timeout 60s ${PKGS} ${TEST_PKGS}

build: test
	go build -v ${BINS}

install: test
	go install ${BINS}

run: build
	./${BIN_SERVER}

image: test
	docker image build -t ${BIN_SERVER_DOCKER_TAG} .

rund: image
	docker container run -it --rm -p 6379:6379 ${BIN_SERVER_DOCKER_TAG}

clean:
	go clean -i ${PKGS}
