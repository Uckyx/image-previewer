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

.PHONY: echo
echo:
	echo ${CI_COMMIT_SHA}

.PHONY: test
test:
	go test -v -count=100 -race -timeout=5m ./...

.PHONY: test-coverage
test-coverage:
	go test -p 1 -v -race -coverprofile cover.out ./...
	go tool cover -html=cover.out -o cover.html

.PHONY: lint
lint:
	golangci-lint run --color always --timeout 30m 2>/dev/null

.PHONY: docker-up
docker-up:
	make build && docker-compose up -d

.PHONY: docker-down
docker-down:
	docker-compose down

.PHONY: build
build:
	$(GO_BUILD) -mod vendor -trimpath -o ./bin/server ./cmd/server