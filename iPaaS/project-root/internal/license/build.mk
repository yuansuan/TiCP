override GOFILES = $(wildcard ${ROOTPATH}/iPaaS/project-root/cmd/license/*.go)
include ${ROOTPATH}/iPaaS/project-root/devops/hacks/build-go-lib.mk

.DEFAULT_GOAL = all
.PHONY: all image

all: go-linux-bin image
image: go-image-base
