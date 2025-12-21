BINDIR:=bin
REVISION:=$(shell git rev-parse --short HEAD)
GOPATH:=$(shell go env GOPATH)

GO_LDFLAGS_VERSION:=-X 'main.Revision=${REVISION}'
GO_LDFLAGS:=$(GO_LDFLAGS_VERSION)
GO_BUILD_OPTION:=-ldflags="-s -w $(GO_LDFLAGS)"

.PHONY: build up stop down destroy clean lint ut it test

build: bin/cli

$(BINDIR)/cli:
	go build -o $@ $(GO_BUILD_OPTION) ./cmd/cli/...

clean:
	go clean -cache -testcache
	rm -f $(BINDIR)/*

fmt:
	@ FORMATTED=$$(go fmt ./... 2>&1); \
	if [ -n "$$FORMATTED" ]; then \
		echo "error: go fmt made changes to the following files:"; \
		echo "$$FORMATTED"; \
		exit 1; \
	else \
		echo "all files are formatted."; \
	fi

lint:
	go vet ./...
	staticcheck ./...

test: fmt lint
	go test -race -cover ./...
