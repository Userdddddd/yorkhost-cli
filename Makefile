BINARY  := yorkhost
PKG     := github.com/YorkHost-fr/yorkhost-cli
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null)
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w \
	-X $(PKG)/internal/version.Version=$(VERSION) \
	-X $(PKG)/internal/version.Commit=$(COMMIT) \
	-X $(PKG)/internal/version.Date=$(DATE)

.PHONY: all build install test race lint fmt vet docs docs-check snapshot clean help

all: lint test build ## lint + tests + build

build: ## Binaire dans ./bin/yorkhost
	CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o bin/$(BINARY) ./cmd/yorkhost

install: ## Installe dans $GOPATH/bin
	CGO_ENABLED=0 go install -trimpath -ldflags "$(LDFLAGS)" ./cmd/yorkhost

test: ## Tests unitaires et de bout en bout (faux serveur API)
	go test ./...

race: ## Tests avec détecteur de concurrence
	go test -race -count=1 ./...

fmt: ## Formate le code
	gofmt -w .

vet: ## go vet
	go vet ./...

lint: vet ## gofmt + go vet (+ golangci-lint s'il est installé)
	@test -z "$$(gofmt -l .)" || (echo "Fichiers non formatés :"; gofmt -l .; exit 1)
	@command -v golangci-lint >/dev/null && golangci-lint run || echo "golangci-lint absent : ignoré"

docs: ## Régénère la référence des commandes (docs/commands)
	rm -rf docs/commands
	go run ./cmd/yorkhost gen-docs docs/commands

docs-check: docs ## Échoue si docs/commands n'est pas à jour
	@git diff --exit-code -- docs/commands || (echo "docs/commands n'est pas à jour : lancez make docs"; exit 1)

snapshot: ## Release locale sans publication (dist/)
	goreleaser release --snapshot --clean

clean: ## Supprime bin/ et dist/
	rm -rf bin dist completions manpages

help: ## Cette aide
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'
