include ${ROOTPATH}/iPaaS/project-root/devops/hacks/build-go-lib.mk

.DEFAULT_GOAL = all
.PHONY: all image push

all: go-linux-bin image push
image: go-image-base
push: push-base
