BINARY := pokedexcli
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: help build run tui test cover vet fmt fmt-check lint clean install snapshot tidy demo demo-tui

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

run: build ## Build and run the REPL
	./$(BINARY)

tui: build ## Build and run the TUI
	./$(BINARY) tui

test: ## Run tests with the race detector
	go test -race -cover ./...

cover: ## Run tests and open an HTML coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

vet: ## Run go vet
	go vet ./...

fmt: ## Format the code
	gofmt -w .

fmt-check: ## Fail if any file is not gofmt-ed
	@test -z "$$(gofmt -l .)" || (echo "Unformatted files:"; gofmt -l .; exit 1)

lint: ## Run golangci-lint (must be installed)
	golangci-lint run

tidy: ## Tidy go.mod/go.sum
	go mod tidy

install: ## Install the binary into GOBIN
	go install -ldflags "$(LDFLAGS)" .

snapshot: ## Build a local release snapshot with GoReleaser
	goreleaser release --snapshot --clean

demo: ## Record the REPL demo GIF (needs asciinema, agg, tmux)
	bash scripts/record-demo.sh

demo-tui: ## Record the TUI demo GIF (needs asciinema, agg, tmux)
	bash scripts/record-tui-demo.sh

clean: ## Remove build artifacts
	rm -f $(BINARY) coverage.out
	rm -rf dist/
