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


ENV ?= env
GOARCH ?= GOARCH
GOOS ?= GOOS
ALLENV ?= $(ENV) $(GOOS)=linux $(GOARCH)=arm64

APPPATH ?= ~/.baagh/
INSTALLPATH ?= ~/go/bin/

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: help all vendor build run root

all: help

## install: Installs the binary
install:
	echo "Installing..."
	$(GO) build -o $(INSTALLPATH)$(PROJECTNAME) $(MAIN)
	echo "Installed."

## rinstall: Installs the binary in race mode
rinstall:
	echo "Installing in race mode"
	$(GO) build -race -o $(INSTALLPATH)$(PROJECTNAME).race $(MAIN)

## build: Build the binary.
build:
	echo "Compiling"
	mkdir -p $(GOBIN)
	$(GO) build -o $(EXECUTABLE) $(MAIN)

## run: Runs the application
run: build
	echo "Running"
	sudo $(EXECUTABLE)

## rbuild: Builds the binary in race mode
rbuild:
	echo "Compiling in race mode"
	mkdir -p $(GOBIN)
	env GOOS=linux GOARCH=arm64 $(GO) build -race -v -o $(EXECUTABLE).race $(MAIN)

# rrun: Runs the binary in race mode
rrun: rbuild
	echo "Running in race mode"
	sudo $(EXECUTABLE).race

## vendor: Run Go Vendor
vendor: 
	$(GO) mod vendor

## clean: Remove build related files.
clean: 
	rm -fr ./out




help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo