# Copyright (C) 2022 The go-redis Authors All rights reserved.
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

GOBIN := $(shell go env GOPATH)/bin
PATH := $(GOBIN):$(PATH)

MODULE_ROOT=github.com/cybergarage/go-redis

PKG_NAME=redis
PKG_VER=$(shell git describe --abbrev=0 --tags)
PKG_COVER=${PKG_NAME}-cover

PKG_ID=${MODULE_ROOT}/${PKG_NAME}
PKG_SRC_DIR=${PKG_NAME}
PKG=${MODULE_ROOT}/${PKG_SRC_DIR}

BIN_DIR=examples
BIN_ID=${MODULE_ROOT}/${BIN_DIR}
BIN_SERVER=go-redisd
BIN_SERVER_DOCKER_TAG=cybergarage/${BIN_SERVER}:${PKG_VER}
BIN_SERVER_ID=${BIN_ID}/${BIN_SERVER}
BINS=\
	${BIN_SERVER_ID}

TEST_PKG_NAME=${PKG_NAME}test
TEST_PKG_ID=${MODULE_ROOT}/${TEST_PKG_NAME}
TEST_PKG_DIR=${TEST_PKG_NAME}
TEST_PKG=${MODULE_ROOT}/${TEST_PKG_DIR}

.PHONY: version format vet lint clean
.IGNORE: lint

all: test

%.md : %.adoc
	asciidoctor -b docbook -a leveloffset=+1 -o - $< | pandoc  --markdown-headings=atx --wrap=preserve -t markdown_strict -f docbook > $@
docs := $(patsubst %.adoc,%.md,$(wildcard *.adoc doc/*.adoc))
doc: $(docs)

version:
	@pushd ${PKG_SRC_DIR} && ./version.gen > version.go && popd
	-git commit v${PKG_SRC_DIR}/version.go -m "Update version"

format: version doc
	gofmt -s -w ${PKG_SRC_DIR} ${BIN_DIR} ${TEST_PKG_DIR}

vet: format
	go vet ${PKG_ID} ${TEST_PKG_ID} ${BINS}

lint: vet
	golangci-lint run ${PKG_SRC_DIR}/... ${BIN_DIR}/... ${TEST_PKG_DIR}/...

test: lint
	go test -v -p 1 -timeout 10m -cover -coverpkg=${PKG}/... -coverprofile=${PKG_COVER}.out ${PKG}/... ${TEST_PKG}/...
	go tool cover -html=${PKG_COVER}.out -o ${PKG_COVER}.html

build: test
	go build -v ${BINS}

install: test
	go install ${BINS}

run: install
	$(GOBIN)/${BIN_SERVER}

image: test
	docker image build -t ${BIN_SERVER_DOCKER_TAG} .

rund: image
	docker container run -it --rm -p 6379:6379 ${BIN_SERVER_DOCKER_TAG}

clean:
	go clean -i ${PKG}
