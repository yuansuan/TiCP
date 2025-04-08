GOLDFLAGS += -X 'github.com/yuansuan/ticp/common/go-kit/gin-boot/version.commitId=$(shell git rev-parse HEAD)'
GOLDFLAGS += -X 'github.com/yuansuan/ticp/common/go-kit/gin-boot/version.branch=$(shell git branch --show-current)'
GOLDFLAGS += -X 'github.com/yuansuan/ticp/common/go-kit/gin-boot/version.buildTime=$(shell date +"%Y-%m-%d %H:%M:%S")'

GOFLAG=-ldflags="$(GOLDFLAGS)"
linux-bin:
	GOOS=linux go build $(GOFLAG) -o storage github.com/yuansuan/ticp/iPaaS/project-root/cmd/storage

override GOFILES = $(wildcard ${ROOTPATH}/iPaaS/project-root/cmd/storage/*.go)
include ${ROOTPATH}/iPaaS/project-root/devops/hacks/build-go-lib.mk

.DEFAULT_GOAL = all
.PHONY: all image

all: go-linux-bin image
image: go-image-base
