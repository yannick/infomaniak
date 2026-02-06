BINARY     := informaniak
MODULE     := github.com/yannick/informaniak
VERSION    ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS    := -trimpath -ldflags "-s -w -X main.version=$(VERSION)"
GOFILES    := $(shell find . -name '*.go' -not -path './vendor/*')

.PHONY: all build clean deps lint test vet fmt help

all: deps build ## Build after fetching deps

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

deps: ## Download all dependencies
	go mod download
	go mod tidy

build: ## Build binary for the current platform
	go build $(LDFLAGS) -o bin/$(BINARY) .

install: build ## Install binary to GOPATH/bin
	go install $(LDFLAGS) .

test: ## Run tests with race detector
	go test -race -count=1 ./...

vet: ## Run go vet
	go vet ./...

fmt: ## Run gofmt check
	@test -z "$$(gofmt -l $(GOFILES))" || (echo "Files need formatting:"; gofmt -l $(GOFILES); exit 1)

lint: vet fmt ## Run all linters
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, skipping"; \
	fi

clean: ## Remove build artifacts
	rm -rf bin/
