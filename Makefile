-include .env

GO ?= go
GOTEST ?= $(GO) tests
GOVET ?= $(GO) vet

VERSION ?= $(shell git describe --tags)
BUILD ?= $(shell git rev-parse --short HEAD)
PROJECTNAME ?= $(shell basename "$(PWD)")
GOFILES := $(wildcard *.go)

# Go related variables.
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/out/bin
MAIN := $(GOBASE)/cmd/baagh



# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: help all vendor build run

all: help

## build: Build the binary.
build: $(GOFILES)
	mkdir -p $(GOBIN)
	env GOOS=linux GOARCH=arm64 $(GO) build -mod vendor -o $(GOBIN) $(MAIN)

run: build
	su - pi -c "$(GOBIN)/baagh"


## clean: Remove build related files.
clean: 
	rm -fr ./bin
	rm -fr ./out

## vendor: Copy of all packages needed to support builds and tests in the vendor directory
vendor: 
	$(GO) mod vendor


help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo