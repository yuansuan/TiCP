IMAGE_NAME ?= $(shell basename "$(PWD)")
VERSION ?= $(shell date +%Y%m%d%H%M)
REGISTRY ?= harbor.yuansuan.cn/ticp

IMAGE_LABELS ?= --label git-hash='$(shell git rev-parse -q HEAD)' --label build-time='$(shell date)'
DOCKER_BUILD_ARGS =
IMAGE_REPOSITORY ?= $(REGISTRY)/$(IMAGE_NAME)
IMAGE_TAG ?= $(VERSION)
IMAGE ?= $(IMAGE_REPOSITORY):$(IMAGE_TAG)
DOCKERFILE ?= Dockerfile
DOCKER_CONTEXT ?= $(shell dirname $(DOCKERFILE))

.PHONY: image-base
image-base:
	DOCKER_BUILDKIT=1 docker build $(IMAGE_LABELS) $(DOCKER_BUILD_ARGS) $(DOCKER_BUILD_ARGS_EXTRA) \
		-t $(IMAGE) -f $(DOCKERFILE) $(DOCKER_CONTEXT)
