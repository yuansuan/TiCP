include ${ROOTPATH}/iPaaS/project-root/devops/hacks/build-lib.mk

.DEFAULT_GOAL = all
.PHONY: image all

all: image
image: image-base
