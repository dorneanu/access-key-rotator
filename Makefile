# Configuration
GO111MODULES=on
BUILD_FLAGS="-ldflags=-s -w"
export GOFLAGS?=-mod=mod

# Basic go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
GOLINT=golangci-lint
GOPATH=$(GOCMD) env GOPATH

# Application configureation
MAIN_BINARY=access-key-rotator
MAIN_LAMBDA=access-key-rotator.lambda
LAMBDA_NAME=AccessKeyRotator

SRC_FILES=cmd/access-key-rotator/main.go
SRC_FILES_LAMBDA=cmd/lambda/access-key-rotator.go

DIST_DIR=./build

# Adjust PATH
export PATH := $(GOPATH):$(PATH)

zip: build-lambda # Make zip for AWS lambda deployment
	zip ${DIST_DIR}/${LAMBDA_NAME}.zip ${DIST_DIR}/${MAIN_LAMBDA}

build: clean lint fmt test build-cli build-lambda

build-cli: # Build CLI binary
	@echo "> Building CLI binary"
	go build -o ${DIST_DIR}/${MAIN_BINARY} ${SRC_FILES}

build-lambda: # Build Lambda
	@echo "> Building lambda binary"
	CGO_ENABLED=0 go build ${BUILD_FLAGS} -o ${DIST_DIR}/${MAIN_LAMBDA} ${SRC_FILES_LAMBDA}

# Install dependencies and addtional tools
install: install-deps install-tools

install-deps: # Install all dependencies 
	@echo "> Download go.mod dependencies"
	@${GOCMD} mod download

install-tools: # Install additional tools for dev and CI/CD
	@echo "> Install dev tools"
	# Install golangci-lint (this is the recommended way)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.40.1

fmt: # Format source code
	@echo "> Formatting source files"
	$(GOFMT) ./...

lint: # Lint source files
	@echo "> Linting source files"
	${GOLINT} run ./...

test: # Run unit tests
	@echo "> Running unit tests"
	$(GOTEST) ./...

clean: # Clean project
	@echo "> Cleaning project"
	$(GOCLEAN)

all: install build