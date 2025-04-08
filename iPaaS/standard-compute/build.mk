override GOFILES = $(wildcard ${ROOTPATH}/iPaaS/standard-compute/cmd/*.go)
override DOCKERFILE = ${ROOTPATH}/iPaaS/standard-compute/docker/Dockerfile
include ${ROOTPATH}/iPaaS/project-root/devops/hacks/build-go-lib.mk

.DEFAULT_GOAL = all
.PHONY: all image

all: go-linux-bin image
image: go-image-base
