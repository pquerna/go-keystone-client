
ROOT_DIR := $(shell pwd)
GOPATH := ${ROOT_DIR}/deps

all: test build
	@echo "Done"

deps:
	go get -u launchpad.net/gocheck

format:
	go fmt

test: build
	go test

build:
	go build

clean:
	go clean

.PHONY: deps clean build test format
