# include env from local
LOCAL_ENV_FILE=local.env
ifneq ("$(wildcard $(LOCAL_ENV_FILE))","")
	include $(LOCAL_ENV_FILE)
	export $(shell sed 's/=.*//' $(LOCAL_ENV_FILE))
endif

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOTOOL=$(GOCMD) tool
GOGET=$(GOCMD) get

EXE= 
ifeq ($(GOOS),windows)
	EXE=.exe
endif

SRC_TARGET=./cmd/red/
BIN_NAME=red
BIN_FOLDER=build/

BIN_TARGET=$(BIN_FOLDER)$(BIN_NAME)$(EXE)

GIT_COMMIT ?= $(shell { git stash create; git rev-parse HEAD; } | grep -Exm1 '[[:xdigit:]]{40}')
VERSION ?= $(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)

export FLAGS += -X "main.Version=$(VERSION)"
export FLAGS += -X "main.GitCommit=$(GIT_COMMIT)"
export FLAGS += -X "main.BuildTime=$(shell date)"

all: test build

deps:
	go get -v -t -d ./...

build: 
	$(GOBUILD) -ldflags='$(FLAGS)' -o $(BIN_TARGET) $(SRC_TARGET)

test:
	$(GOTEST) -v ./... -cover

cover:
	$(GOTEST) -v -coverprofile=./build/c.out ./...
	$(GOTOOL) cover -html=./build/c.out -o ./build/coverage.html

coverAll:
	$(GOTEST) -v -tags=integration -coverprofile=./build/c.out ./...
	$(GOTOOL) cover -html=./build/c.out -o ./build/coverage.html

clean:
	$(GOCLEAN)
	rm -r build

run:
	@$(GOBUILD) -o $(BIN_TARGET) $(SRC_TARGET)
	$(BIN_TARGET) $(args)

# Cross compilation
build-ubuntu:
	GOOS=linux GARCH=amd64 $(GOBUILD) -o $(BIN_TARGET) $(SRC_TARGET)
build-arm:
	GOOS=linux GOARCH=arm GOARM=7 $(GOBUILD) -o $(BIN_TARGET) $(SRC_TARGET)
build-win:
	GOOS=windows GOARCH=386 $(GOBUILD) -o $(BIN_TARGET) $(SRC_TARGET)
build-mac:
	GOOS=darwin GOARCH=386 $(GOBUILD) -o $(BIN_TARGET) $(SRC_TARGET)