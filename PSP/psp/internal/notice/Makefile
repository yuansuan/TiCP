# Go parameters

ROOTPATH=$(abspath $(CURDIR)/../..)
PROJECT_NAME=$(shell basename "$(PWD)")
GO_BASE=$(shell pwd)
GO_FILES=$(wildcard *.go)

DIR_ROUTER=router

PARSER_CMD=parser
PROTO2GO_CMD=p2g
GEN_CMD=gen

OS_WIN=MINGW64_NT MINGW32_NT
OS_VERSION_INFO=$(shell uname -s)
OS_VERSION=$(shell echo "$(OS_VERSION_INFO)"|cut -c 1-10)

ifeq ($(OS_VERSION),$(findstring $(OS_VERSION), $(OS_WIN)))
	FIRST_GO_PATH := $(shell echo "$(ROOTPATH)" |awk -F ';' '{print $$1}')
else
	FIRST_GO_PATH := $(shell echo "$(ROOTPATH)" |awk -F ':' '{print $$1}')
endif

SHPATH=$(FIRST_GO_PATH)/devops/build/bin/cli.sh

MAKE_CMD=gen p2g
ifeq ($(firstword $(MAKECMDGOALS)),$(findstring $(firstword $(MAKECMDGOALS)), $(MAKE_CMD)))
    RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
    $(eval $(RUN_ARGS):;@:)
endif

mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
cur_makefile_path := $(dir $(mkfile_path))



all: help


## build:	make build (generate router and convert *.proto to *.pb.go and vet main.go)
build:
	@echo "  >  Building binary..."
	@#echo "  	>  Parser router..."
	@#sh $(SHPATH) $(PARSER_CMD) $(cur_makefile_path)/handler
	@echo "  	>  Parser proto..."
	@sh $(SHPATH) $(GEN_CMD) $(PROTO2GO_CMD) ${ROOTPATH}/internal/common/proto
	@#echo "  	>  Vet main.go..."
	@#go vet $(GOFILES)
	@#go build -race -o $(PROJECTNAME) $(GOFILES)

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
	@find $(CURDIR) -name "*.go" -exec goimports-reviser -company-prefixes yuansuan.cn/psp -company-prefixes yuansuan.cn -rm-unused  {} \;