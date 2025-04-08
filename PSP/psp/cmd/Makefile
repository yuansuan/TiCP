# Go parameters
ROOTPATH=$(abspath $(CURDIR)/..)

PROJECTNAME=$(shell basename "$(PWD)")
GOBASE=$(shell pwd)
GOFILES=$(wildcard *.go)

PSPNAME=pspd
BOLTDB=rfm_file_info_db

DIR_ROUTER=router

PARSER_CMD=parser
PROTO2GO_CMD=p2g
GEN_CMD=gen

REMOTE_ADDRESS=registry.intern.yuansuan.cn
BUILD_IMAGE=$(REMOTE_ADDRESS)/psp-build:latest

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
GIT_CONFIG_CMD=git config --global --add url."git@code.intern.yuansuan.cn:".insteadOf "https://code.intern.yuansuan.cn/"

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

LDFLAGS += -X "main.BuildDatetime=$(shell TZ='Asia/Shanghai' date '+%Y-%m-%d %I:%M:%S')"
LDFLAGS += -X "main.GoVersion=$(shell go version)"
LDFLAGS += -X "main.GitBranch=$(shell git rev-parse --abbrev-ref HEAD)"
LDFLAGS += -X "main.GitHash=$(shell git rev-parse HEAD)"

all: help

proto:
	@echo "  	>  Parser proto..."
	@export ROOTPATH=${ROOTPATH} && export PATH="${ROOTPATH}/devops/build/bin/protoc-gen:${PATH}" \
	 && $(SHPATH) $(GEN_CMD) $(PROTO2GO_CMD) ${ROOTPATH}/internal/common/proto

docker-build:
	@echo "  >  Building binary..."
	@echo "  	>  Vet main.go..."
	@go vet $(GOFILES)
	@echo "  	>  go build..."
	@go build -v -x -ldflags '$(LDFLAGS)' -gcflags "all=-N -l" -race -o $(PSPNAME) $(GOFILES)

docker-build-psp: proto
	@docker run --rm --network host -v ${ROOTPATH}:${ROOTPATH} -v ~/.cache/ys/psp/go-build:/root/.cache/go-build -v ~/.cache/ys/psp/mod:/root/go/pkg/mod -v ~/.ssh:/root/.ssh  \
	$(BUILD_IMAGE) bash -c "cd $(shell pwd) && ${GIT_CONFIG_CMD} && time make docker-build"

docker-image-build-push:
	@docker $(DOCKER_BUILD_COMMAND) -t $(BUILD_IMAGE) -f ${ROOTPATH}/docker/image/Dockerfile ${ROOTPATH}/docker/image
	@docker push $(BUILD_IMAGE)

## build:	make build (generate router and convert *.proto to *.pb.go and vet main.go)
build: proto
	@echo "  	>  Vet main.go..."
	@go vet $(GOFILES)
	@echo "  	>  Go build..."
	@go build -race -ldflags '$(LDFLAGS)' -o $(PSPNAME) $(GOFILES)

.PHONY : gen
## gen:	make gen [j2g/y2g/p2g] [file or dir(include .json or.yml or .proto)]
gen:
	@sh $(SHPATH) $(GEN_CMD) $(RUN_ARGS)


.PHONY : help
help: Makefile
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@echo -e "\033[31m\033[1m Usage: \033[0m"
	@echo
	@echo "   make command [arguments]:"
	@echo
	@echo " The commands are: "
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

clean:
	@rm -f $(PROJECTNAME) $(PSPNAME) $(BOLTDB)

gofmt:
	@go install -v github.com/incu6us/goimports-reviser/v3@latest
	@find $(CURDIR) -name "*.go" -exec goimports-reviser -company-prefixes yuansuan.cn/psp -company-prefixes yuansuan.cn -rm-unused  {} \;
