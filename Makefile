# ORCA Makefile

# Metadata
VERSION ?= $(shell cat VERSION)
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Binary names
BINARY_NAME=orca
BINARY_PATH=bin/$(BINARY_NAME)

# Build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.gitCommit=$(GIT_COMMIT) -X main.buildDate=$(BUILD_DATE)"

# Docker parameters
DOCKER_IMAGE=orca
DOCKER_TAG=$(VERSION)

.PHONY: all build clean test coverage lint fmt vet mod-tidy mod-download install help

all: test build

## build: Build the binary
build:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	@mkdir -p bin
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_PATH) ./cmd/orca

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf bin/
	@rm -rf dist/
	@rm -f coverage.txt coverage.html

## test: Run unit tests
test:
	@echo "Running unit tests..."
	$(GOTEST) -v -race -short ./...

## test-all: Run all tests including integration
test-all:
	@echo "Running all tests..."
	$(GOTEST) -v -race ./...

## integration-test: Run integration tests only
integration-test:
	@echo "Running integration tests..."
	$(GOTEST) -v -race -tags=integration ./...

## smoke-test: Run smoke tests against deployed cluster
smoke-test:
	@echo "Running smoke tests..."
	./scripts/smoke-test.sh

## coverage: Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	$(GOCMD) tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated: coverage.html"

## lint: Run linters
lint:
	@echo "Running linters..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run --timeout 5m ./...

## fmt: Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

## vet: Run go vet
vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...

## mod-tidy: Tidy go modules
mod-tidy:
	@echo "Tidying go modules..."
	$(GOMOD) tidy

## mod-download: Download go modules
mod-download:
	@echo "Downloading go modules..."
	$(GOMOD) download

## install: Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	@cp $(BINARY_PATH) $(GOPATH)/bin/$(BINARY_NAME)

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	docker tag $(DOCKER_IMAGE):$(DOCKER_TAG) $(DOCKER_IMAGE):latest

## docker-push: Push Docker image
docker-push:
	@echo "Pushing Docker image..."
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push $(DOCKER_IMAGE):latest

## run: Run the application locally
run: build
	@echo "Running $(BINARY_NAME)..."
	$(BINARY_PATH) --config config.yaml

## help: Show this help message
help:
	@echo "ORCA Makefile Commands:"
	@echo ""
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
