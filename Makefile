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

.PHONY: all build clean test coverage lint fmt vet mod-tidy mod-download install help localstack-start localstack-stop localstack-logs localstack-status localstack-restart run-local docs docs-serve docs-build docs-deploy pre-commit-install pre-commit-run release-snapshot

all: lint test build

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

## integration-test: Run integration tests only (requires LocalStack)
integration-test:
	@echo "Running integration tests..."
	@if ! docker ps | grep -q orca-localstack; then \
		echo "❌ LocalStack is not running. Start it with 'make localstack-start'"; \
		exit 1; \
	fi
	@./scripts/wait-for-localstack.sh
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

## localstack-start: Start LocalStack for testing
localstack-start:
	@echo "Starting LocalStack..."
	@docker-compose -f docker-compose.localstack.yml up -d
	@echo "Waiting for LocalStack to initialize..."
	@sleep 5
	@echo "✅ LocalStack is starting. Check logs with 'make localstack-logs'"
	@echo "Resources will be initialized automatically."

## localstack-stop: Stop LocalStack
localstack-stop:
	@echo "Stopping LocalStack..."
	@docker-compose -f docker-compose.localstack.yml down

## localstack-logs: Show LocalStack logs
localstack-logs:
	@docker-compose -f docker-compose.localstack.yml logs -f

## localstack-status: Show LocalStack resource IDs
localstack-status:
	@echo "LocalStack Resources:"
	@if [ -f /tmp/localstack-orca-resources.env ]; then \
		cat /tmp/localstack-orca-resources.env; \
	else \
		echo "Resources not yet initialized. Check 'make localstack-logs'"; \
	fi

## localstack-restart: Restart LocalStack
localstack-restart: localstack-stop localstack-start

## run-local: Run ORCA with LocalStack config
run-local: build
	@if ! docker ps | grep -q orca-localstack; then \
		echo "❌ LocalStack is not running. Start it with 'make localstack-start'"; \
		exit 1; \
	fi
	@echo "Running $(BINARY_NAME) with LocalStack configuration..."
	@$(BINARY_PATH) --config config.localstack.yaml

## docs-serve: Serve documentation locally
docs-serve:
	@echo "Serving documentation at http://127.0.0.1:8000"
	@which mkdocs > /dev/null || (echo "❌ mkdocs not installed. Run: pip install mkdocs-material mkdocs-minify-plugin" && exit 1)
	mkdocs serve

## docs-build: Build documentation
docs-build:
	@echo "Building documentation..."
	@which mkdocs > /dev/null || (echo "❌ mkdocs not installed. Run: pip install mkdocs-material mkdocs-minify-plugin" && exit 1)
	mkdocs build --strict

## docs-deploy: Deploy documentation to GitHub Pages
docs-deploy:
	@echo "Deploying documentation to GitHub Pages..."
	@which mkdocs > /dev/null || (echo "❌ mkdocs not installed. Run: pip install mkdocs-material mkdocs-minify-plugin" && exit 1)
	mkdocs gh-deploy --force

## pre-commit-install: Install pre-commit hooks
pre-commit-install:
	@echo "Installing pre-commit hooks..."
	@which pre-commit > /dev/null || (echo "❌ pre-commit not installed. Run: pip install pre-commit" && exit 1)
	pre-commit install
	@echo "✅ Pre-commit hooks installed"

## pre-commit-run: Run pre-commit hooks on all files
pre-commit-run:
	@echo "Running pre-commit hooks..."
	@which pre-commit > /dev/null || (echo "❌ pre-commit not installed. Run: pip install pre-commit" && exit 1)
	pre-commit run --all-files

## release-snapshot: Create a snapshot release with GoReleaser
release-snapshot:
	@echo "Creating snapshot release..."
	@which goreleaser > /dev/null || (echo "❌ goreleaser not installed. Run: go install github.com/goreleaser/goreleaser@latest" && exit 1)
	goreleaser release --snapshot --clean

## release: Create a production release with GoReleaser
release:
	@echo "Creating production release..."
	@which goreleaser > /dev/null || (echo "❌ goreleaser not installed. Run: go install github.com/goreleaser/goreleaser@latest" && exit 1)
	@if [ -z "$(VERSION)" ]; then echo "❌ VERSION not set"; exit 1; fi
	goreleaser release --clean

## deps: Install development dependencies
deps:
	@echo "Installing development dependencies..."
	@echo "Installing Go tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/goreleaser/goreleaser@latest
	@echo "Installing Python tools..."
	pip install --upgrade pip
	pip install mkdocs-material mkdocs-minify-plugin pre-commit
	@echo "✅ Dependencies installed"

## setup: Initial project setup
setup: deps mod-download pre-commit-install
	@echo "✅ Project setup complete"

## help: Show this help message
help:
	@echo "ORCA Makefile Commands:"
	@echo ""
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'
