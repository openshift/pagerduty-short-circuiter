unexport GOFLAGS

GOOS?=linux
TESTOPTS ?=
GOARCH?=amd64
GOENV=GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 GOFLAGS=
GOPATH := $(shell go env GOPATH)
HOME=$(shell mktemp -d)
GOLANGCI_LINT_VERSION=v1.51.2

# Ensure go modules are enabled:
export GO111MODULE=on
export GOPROXY=https://proxy.golang.org

# Disable CGO so that we always generate static binaries:
export CGO_ENABLED=0

.PHONY: build
build:
	go build -o kite ./cmd/kite

.PHONY: install
install:
	go build -o ${GOPATH}/bin/kite ./cmd/kite

.PHONY: mod
mod:
	go mod tidy

.PHONY: tools
tools:
	@mkdir -p $(GOPATH)/bin
	@ls $(GOPATH)/bin/ginkgo 1>/dev/null || (echo "Installing ginkgo..." && go install github.com/onsi/ginkgo/ginkgo@v1.16.4)
	@ls $(GOPATH)/bin/mockgen 1>/dev/null || (echo "Installing gomock..." && go install github.com/golang/mock/mockgen@v1.6.0)
	
.PHONY: test
test:
	go test ./... -v $(TESTOPTS)

.PHONY: coverage
coverage:
	chmod 755 hack/codecov.sh
	hack/codecov.sh

# Installed using instructions from: https://golangci-lint.run/usage/install/#linux-and-windows
getlint:
	@mkdir -p $(GOPATH)/bin
	@ls $(GOPATH)/bin/golangci-lint 1>/dev/null || (echo "Installing golangci-lint..." && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin $(GOLANGCI_LINT_VERSION))

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
