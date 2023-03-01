unexport GOFLAGS

GOOS?=linux
GOARCH?=amd64
GOENV=GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 GOFLAGS=
GOPATH := $(shell go env GOPATH)

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
	go test ./... -covermode=atomic -coverpkg=./... -v

# Installed using instructions from: https://golangci-lint.run/usage/install/#linux-and-windows
getlint:
	@mkdir -p $(GOPATH)/bin
	@ls $(GOPATH)/bin/golangci-lint 1>/dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)

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
