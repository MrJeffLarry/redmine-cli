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
GOINSTALL=$(GOCMD) install

EXE= 
ifeq ($(GOOS),windows)
	EXE=.exe
endif

DOC_SRC_TARGET=./cmd/docs/
DOC_BIN_NAME=docs
SRC_TARGET=./cmd/red/
BIN_NAME=red-cli
BIN_FOLDER=build/


# Get commit hash (fallback to 'unknown' if not available)
GIT_COMMIT ?= $(shell { git stash create; git rev-parse HEAD; } | grep -Exm1 '[[:xdigit:]]{40}' || echo "unknown")
# Get version from branch or tag, fallback to 'dev' if not available
VERSION_TMP := $(shell git symbolic-ref -q --short HEAD 2>/dev/null || git describe --tags --exact-match 2>/dev/null || echo "dev")
VERSION ?= $(VERSION_TMP)

BIN_TARGET=$(BIN_FOLDER)$(BIN_NAME)-$(VERSION)$(EXE)
DOC_TARGET=$(BIN_FOLDER)$(DOC_BIN_NAME)-$(VERSION)$(EXE)

export FLAGS += -X "main.version=$(VERSION)"
export FLAGS += -X "main.commit=$(GIT_COMMIT)"
export FLAGS += -X "main.date=$(shell date)"

all: test build

install:
	$(GOINSTALL) ./cmd/red

deps:
	$(GOGET) -v -t -d ./...

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
	@$(GOBUILD) -ldflags='$(FLAGS)' -o $(BIN_TARGET) $(SRC_TARGET)
	$(BIN_TARGET) $(args)

gen-docs:
	@$(GOBUILD) -ldflags='$(FLAGS)' -o $(DOC_TARGET) $(DOC_SRC_TARGET)
	$(DOC_TARGET) $(args)

build-all:
	GOOS=linux GARCH=amd64 $(GOBUILD) -o $(BIN_FOLDER)$(BIN_NAME)_linux_amd64 $(SRC_TARGET)
	GOOS=linux GOARCH=arm GOARM=7 $(GOBUILD) -o $(BIN_FOLDER)$(BIN_NAME)_linux_arm7 $(SRC_TARGET)
	GOOS=windows GOARCH=386 $(GOBUILD) -o $(BIN_FOLDER)$(BIN_NAME).exe $(SRC_TARGET)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BIN_FOLDER)$(BIN_NAME)_mac_amd64 $(SRC_TARGET)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BIN_FOLDER)$(BIN_NAME)_mac_arm64 $(SRC_TARGET)

# Cross compilation
build-ubuntu:
	GOOS=linux GARCH=amd64 $(GOBUILD) -o $(BIN_TARGET) $(SRC_TARGET)
build-arm:
	GOOS=linux GOARCH=arm GOARM=7 $(GOBUILD) -o $(BIN_TARGET) $(SRC_TARGET)
build-win:
	GOOS=windows GOARCH=386 $(GOBUILD) -o $(BIN_TARGET) $(SRC_TARGET)
build-mac:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BIN_TARGET) $(SRC_TARGET)

# dist list
# go tool dist list
