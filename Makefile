-include .env

GO ?= go
GOTEST ?= $(GO) tests
GOVET ?= $(GO) vet

VERSION ?= $(shell git describe --tags)
BUILD ?= $(shell git rev-parse --short HEAD)
PROJECTNAME ?= $(shell basename "$(PWD)")

# Go related variables.
GOBASE ?= $(shell pwd)
GOBIN ?= $(GOBASE)/out/bin
EXECUTABLE ?= $(GOBIN)/$(PROJECTNAME)
MAIN ?= $(GOBASE)/cmd/baagh



# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: help all vendor build run root

all: help

## build: Build the binary.
build: 
	echo "Compiling"
	mkdir -p $(GOBIN)
	env GOOS=linux GOARCH=arm64 $(GO) build -mod vendor -o $(EXECUTABLE) $(MAIN)

install: 
	echo "Installing..."
	$(GO) install -mod vendor $(MAIN)
	echo "Installed."

## run: runs the application
run: build
	echo "Running"
	$(EXECUTABLE)

## clean: Remove build related files.
clean: 
	rm -fr ./out

## vendor: Copy of all packages needed to support builds and tests in the vendor directory
vendor: 
	$(GO) mod vendor

## root: by running this, you can access 
root:
	sudo python3 $(GOBASE)/scripts/create_gpio_user_permissions.py

help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo