ifeq ($(OS),Windows_NT)
    GOOS := windows
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        GOOS := linux
    endif
    ifeq ($(UNAME_S),Darwin)
        GOOS := darwin
    endif
endif

CI_COMMIT_SHA ?= local
CGO_ENABLED = 0
GOARCH = amd64
LDFLAGS = -ldflags "-X main.shaCommit=${CI_COMMIT_SHA}"
GO = $(shell which go)
GO_BUILD = GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) $(GO) build $(LDFLAGS)

.PHONY: build-server
build-server:
	$(GO_BUILD) -mod vendor -trimpath -o ./bin/server ./cmd/server

.PHONY: run-server
run-server:
	$(GO) run ./cmd/server

.PHONY: test-single
test-single:
	$(GO) test -p 2 -v -count=1 -timeout=1m ./pkg/... ./internal/...

.PHONY: test-race
test-race:
	$(GO) test -p 2 -v -count=100 -race -timeout=8m ./pkg/... ./internal/...

.PHONY: test-coverage
test-coverage:
	$(GO) test -p 2 -v -race -coverprofile cover.out ./pkg/... ./internal/...
	$(GO) tool cover -html=cover.out -o cover.html

.PHONY: test-integration
test-integration:
	$(GO) clean -testcache
	$(GO) test -short ./integration_test/... -v

.PHONY: lint
lint:
	golangci-lint run --color always --timeout 30m 2>/dev/null

.PHONY: start
start:
	docker-compose up -d --remove-orphans

.PHONY: stop
stop:
	docker-compose down --remove-orphans
