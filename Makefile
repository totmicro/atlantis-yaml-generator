PROJECTNAME=$(shell basename "$(PWD)")
CURRENT=$(shell echo $(PWD))

.PHONY: help

help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## install: Install missing dependencies.
install:
	go mod download

## Run tests.
tests: main.go
	go test -v ./...

## Check coverage.
coverage: main.go
	go test  ./... -cover

## build: Builds the project.
build: main.go
	go build -o build/atlantis-yaml-generator

## build-all: Build amd64 and arm64 binaries for linux and darwin.
build-all: main.go
	for arch in amd64 arm64 ; do \
		for os in linux darwin; do \
			CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -o "build/atlantis-yaml-generator_"$$os"_$$arch" $(LDFLAGS) ; \
		done; \
	done;
	/bin/chmod +x build/*


