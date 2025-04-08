override BUILD_CONTEXT ?= $(PWD)
override DOCKERFILE ?= ${ROOTPATH}/iPaaS/project-root/devops/docker/Dockerfile
ifdef ENABLED_DEBUG
override DOCKERFILE = ${ROOTPATH}/iPaaS/project-root/devops/docker/golang.Dockerfile
endif

include ${ROOTPATH}/iPaaS/project-root/devops/hacks/build-lib.mk
GOFILES ?= $(wildcard *.go)

GOLDFLAGS=-compressdwarf=false
EXTRA_LDFLAGS ?=
ifndef ENABLED_DEBUG
GOLDFLAGS+=-w -s
endif

## for protobuf conflict
GOLDFLAGS += -X 'google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=warn'
# inject commitId/branch/buildTime to gin-boot
GOLDFLAGS += -X 'github.com/yuansuan/ticp/common/go-kit/gin-boot/version.commitId=$(shell git rev-parse HEAD)'
GOLDFLAGS += -X 'github.com/yuansuan/ticp/common/go-kit/gin-boot/version.branch=$(shell git branch --show-current)'
GOLDFLAGS += -X 'github.com/yuansuan/ticp/common/go-kit/gin-boot/version.buildTime=$(shell date +"%Y-%m-%d %H:%M:%S")'

GOFLAG=-ldflags="$(GOLDFLAGS) $(EXTRA_LDFLAGS)"

.PHONY: go-linux-bin
go-linux-bin:
	cd $(BUILD_CONTEXT)
	@go mod tidy
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(GOFLAG) -o $(BUILD_CONTEXT)/$(IMAGE_NAME) $(GOFILES)

.PHONY: go-image-base
go-image-base: image-base

