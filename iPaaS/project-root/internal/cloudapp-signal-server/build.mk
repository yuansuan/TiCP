override GOFILES = $(wildcard ${ROOTPATH}/iPaaS/project-root/cmd/cloudapp-signal-server/*.go)
include ${ROOTPATH}/iPaaS/project-root/devops/hacks/build-go-lib.mk

.DEFAULT_GOAL = all
.PHONY: image all

all: go-linux-bin image
image: go-image-base
