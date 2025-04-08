override GOFILES = $(wildcard cmd/*.go)
override DOCKERFILE= $(abspath ./Dockerfile)
include ${ROOTPATH}/iPaaS/project-root/devops/hacks/build-go-lib.mk

.DEFAULT_GOAL = all
.PHONY: image all

all: go-linux-bin image
image: go-image-base

authenticator.exe: tools/authenticator/main.go
	GOOS=windows go build -o $@ $^

mocker: tools/mocker/main.go
	go build -o $@ $^
