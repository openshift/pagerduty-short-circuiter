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

.PHONY: build
build:
	go build -o pdcli ./cmd/pdcli

.PHONY: install
install:
	go build -o ${GOPATH}/bin/pdcli ./cmd/pdcli

.PHONY: tools
tools:
	@mkdir -p $(GOPATH)/bin
	@ls $(GOPATH)/bin/ginkgo 1>/dev/null || (echo "Installing ginkgo..." && go get -u github.com/onsi/ginkgo/ginkgo@v1.16.4)
	@ls $(GOPATH)/bin/mockgen 1>/dev/null || (echo "Installing gomock..." && go get -u github.com/golang/mock/mockgen@v1.6.0)
	
.PHONY: test
test:
	ginkgo -v -r tests

# Installed using instructions from: https://golangci-lint.run/usage/install/#linux-and-windows
getlint:
	@mkdir -p $(GOPATH)/bin
	@ls $(GOPATH)/bin/golangci-lint 1>/dev/null || (echo "Installing golangci-lint..." && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.42.0)

.PHONY: lint
lint: getlint
	$(GOPATH)/bin/golangci-lint run

.PHONY: fmt
fmt:
	gofmt -s -l -w cmd pkg tests

.PHONY: clean
clean:
	rm -rf \
		*-darwin-amd64 \
		*-linux-amd64 \
		*-windows-amd64 \
		*.sha256 \
		$(NULL)
