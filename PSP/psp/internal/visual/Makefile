# Go parameters
ROOTPATH=$(abspath $(CURDIR)/../..)
PROJECTNAME=$(shell basename "$(PWD)")
GOBASE=$(shell pwd)
GOFILES=$(wildcard *.go)

DIR_ROUTER=router
lib=src/gin-boot

PARSER_CMD=parser
PROTO2GO_CMD=p2g
GEN_CMD=gen

GOFLAG="-ldflags=-compressdwarf=false"

OS_WIN=MINGW64_NT MINGW32_NT
OS_VERSION_INFO=$(shell uname -s)
OS_VERSION=$(shell echo "$(OS_VERSION_INFO)"|cut -c 1-10)

ifeq ($(OS_VERSION),$(findstring $(OS_VERSION), $(OS_WIN)))
	FIRST_GO_PATH := $(shell echo "$(ROOTPATH)" |awk -F ';' '{print $$1}')
else
	FIRST_GO_PATH := $(shell echo "$(ROOTPATH)" |awk -F ':' '{print $$1}')
endif

SHPATH=$(FIRST_GO_PATH)/$(lib)/bin/cli.sh

MAKECMD=gen p2g
ifeq ($(firstword $(MAKECMDGOALS)),$(findstring $(firstword $(MAKECMDGOALS)), $(MAKECMD)))
    RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
    $(eval $(RUN_ARGS):;@:)
endif

mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
cur_makefile_path := $(dir $(mkfile_path))



all: help


## build:	make build (generate router and convert *.proto to *.pb.go and vet main.go)
build:
	@echo "  >  Building binary..."
	@echo "  	>  Parser router..."
	@sh $(SHPATH) $(PARSER_CMD) $(cur_makefile_path)/handler
	@echo "  	>  Parser proto..."
	@sh $(SHPATH) $(GEN_CMD) $(PROTO2GO_CMD) ${ROOTPATH}/src/yuansuan.cn/protos
	@echo "  	>  Vet main.go..."
	@go vet $(GOFILES)
	@go build $(GOFLAG) -o  $(PROJECTNAME) $(GOFILES)

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
	@rm -f $(PROJECTNAME)

gofmt:
	@go install -v github.com/incu6us/goimports-reviser/v3@latest
	@find $(CURDIR) -name "*.go" -exec goimports-reviser -company-prefixes yuansuan.cn/psp -company-prefixes yuansuan.cn -rm-unused  {} \;