#
# Copyright (c) 2021 Red Hat, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# Ensure go modules are enabled:
export GO111MODULE=on
export GOPROXY=https://proxy.golang.org

# Disable CGO so that we always generate static binaries:
export CGO_ENABLED=0

# Constants
GOPATH := $(shell go env GOPATH)

# Allow overriding: `make lint container_runner=docker`.
container_runner:=podman

.PHONY: build
build:
	go build -o pdcli ./cmd/pdcli

.PHONY: install
install:
	go build -o ${GOPATH}/bin/pdcli ./cmd/pdcli

.PHONY: test
test:
	ginkgo -v -r tests

.PHONY: lint
lint:
	$(container_runner) run --rm --security-opt label=disable --volume="$(PWD):/app" --workdir=/app \
		golangci/golangci-lint:v$(shell cat .golangciversion) \
		golangci-lint run

.PHONY: clean
clean:
	rm -rf \
		$$(ls cmd) \
		*-darwin-amd64 \
		*-linux-amd64 \
		*-windows-amd64 \
		*.sha256 \
		$(NULL)

.PHONY: tools
tools: ## Install tools to ${GOPATH}/bin
	go get -u github.com/onsi/ginkgo/ginkgo
