# Go parameters
ROOTPATH=$(abspath $(CURDIR)/..)

PROJECTNAME=$(shell basename "$(PWD)")
GOBASE=$(shell pwd)
GOFILES=$(wildcard *.go)

PSP_AGENT_NAME=psp_agent
BOLTDB=rfm_file_info_db

DIR_ROUTER=router

PARSER_CMD=parser
PROTO2GO_CMD=p2g
GEN_CMD=gen

REMOTE_ADDRESS=registry.intern.yuansuan.cn
PSP_APP=psp-app

BUILD_IMAGE=$(REMOTE_ADDRESS)/$(PSP_APP)-build:latest

OS_WIN=MINGW64_NT MINGW32_NT
OS_VERSION_INFO=$(shell uname -s)
OS_VERSION=$(shell echo "$(OS_VERSION_INFO)"|cut -c 1-10)

TAG=$(shell git describe --abbrev=0 --tags || echo "none")
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
CommitID=$(shell git log --pretty=format:"%h" -1)
DATE=$(shell date)

ifeq ($(OS_VERSION),$(findstring $(OS_VERSION), $(OS_WIN)))
	FIRST_GO_PATH := $(shell echo "$(ROOTPATH)" |awk -F ';' '{print $$1}')
else
	FIRST_GO_PATH := $(shell echo "$(ROOTPATH)" |awk -F ':' '{print $$1}')
endif

SHPATH=$(FIRST_GO_PATH)/devops/build/bin/cli.sh

MAKECMD=gen p2g
ifeq ($(firstword $(MAKECMDGOALS)),$(findstring $(firstword $(MAKECMDGOALS)), $(MAKECMD)))
    RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
    $(eval $(RUN_ARGS):;@:)
endif

mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
cur_makefile_path := $(dir $(mkfile_path))

DOCKER_BUILD_COMMAND = build
ifneq ($(shell uname -m | grep -E 'arm*'), "")
	DOCKER_BUILD_COMMAND = buildx build --platform=linux/amd64
endif

LDFLAGS += -X "main.BuildDatetime=$(shell date -u '+%Y-%m-%d %I:%M:%S')"
LDFLAGS += -X "main.GoVersion=$(shell go version)"
LDFLAGS += -X "main.GitBranch=$(shell git rev-parse --abbrev-ref HEAD)"
LDFLAGS += -X "main.GitHash=$(shell git rev-parse HEAD)"

all: help

docker-build:
	@echo "  >  Building binary..."
	@echo "  	>  Vet main.go..."
	@go vet $(GOFILES)
	@echo "  	>  go build..."
	@go build -v -x -ldflags '$(LDFLAGS)' -gcflags "all=-N -l" -race -o $(PSP_AGENT_NAME) $(GOFILES)

docker-build-agent:
	@docker run --rm -v ${ROOTPATH}:${ROOTPATH} -v ~/.cache/ys/psp/go-build:/root/.cache/go-build -v ~/.cache/ys/psp/mod:/workspace/pkg/mod -v ~/.ssh:/root/.ssh  \
	$(BUILD_IMAGE) bash -c "cd $(shell pwd) && time make docker-build"

docker-image-build-push:
	@docker $(DOCKER_BUILD_COMMAND) -t $(BUILD_IMAGE) -f ${ROOTPATH}/docker/Dockerfile ${ROOTPATH}/docker
	@docker push $(BUILD_IMAGE)

gofmt:
	@go install -v github.com/incu6us/goimports-reviser/v3@latest
	@find $(CURDIR) -name "*.go" -exec goimports-reviser -company-prefixes yuansuan.cn/psp -company-prefixes yuansuan.cn -rm-unused  {} \;