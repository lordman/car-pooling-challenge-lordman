# Makefile for car-pooling-challenge
# vim: set ft=make ts=8 noet
# Copyright Cabify.com
# Licence MIT

# Basic go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Binary names
BINARY_NAME=car-pooling-challenge
BINARY_UNIX=$(BINARY_NAME)_unix

# Variables
UNAME		:= $(shell uname -s)

.EXPORT_ALL_VARIABLES:

# this is godly
# https://news.ycombinator.com/item?id=11939200
.PHONY: help
help:	### this screen. Keep it first target to be default
ifeq ($(UNAME), Linux)
	@grep -P '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
else
	@# this is not tested, but prepared in advance for you, Mac drivers
	@awk -F ':.*###' '$$0 ~ FS {printf "%15s%s\n", $$1 ":", $$2}' \
		$(MAKEFILE_LIST) | grep -v '@awk' | sort
endif

# Targets
#
.PHONY: all
all: test build	### Execute the tests and build the app binary

.PHONY: build
build:	### Build the app binary
	$(GOBUILD) -o $(BINARY_NAME) -v

.PHONY: test
test:	### Execute the tests (not implemented yet)
	$(GOTEST) -v ./...

.PHONY: clean
clean:	### Remove files
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

.PHONY: run
run:	### Build and run the app
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

.PHONY: deps
deps:	### Download and install packages and dependencies
	$(GOGET) gopkg.in/go-playground/validator.v9
	$(GOGET) github.com/gorilla/mux
	
.PHONY: debug
debug:	### Debug Makefile itself
	@echo $(UNAME)

.PHONY: dockerize
dockerize: build	### Build Docker image
	@docker build -t car-pooling-challenge:latest .
