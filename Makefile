# ----- Makefile for loghub -----
# Run `make build` or simply `make` (default target).

.PHONY: all build test lint fmt tidy

# Default target
all: build

## Build the CLI -----------------------------------------------------
build:
	@echo "→ Building loghub..."
	go build -o loghub.exe .    

## Run unit tests ----------------------------------------------------
test:
	@echo "→ Running unit tests..."
	go test ./... -v

## Run linters (golangci-lint) --------------------------------------
lint:
	@echo "→ Linting..."
	golangci-lint run

## goimports + gofmt -------------------------------------------------
fmt:
	@echo "→ Formatting code..."
	goimports -w $(shell go env GOPATH)/pkg/mod
	go fmt ./...

## Clean dependencies, ensure go.sum tidy ---------------------------
tidy:
	@echo "→ Tidying modules..."
	go mod tidy
