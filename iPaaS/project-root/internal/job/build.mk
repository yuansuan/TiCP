override GOFILES = $(wildcard ${ROOTPATH}/iPaaS/project-root/cmd/job/*.go)
include ${ROOTPATH}/iPaaS/project-root/devops/hacks/build-go-lib.mk
override EXTRA_LDFLAGS := -X 'yuansuan.cn/project-root/internal/job/handler_rpc.NSPrefix=$(NS)' # 放在build-go-lib.mk后面，NS值在build-go-lib.mk里定义

.DEFAULT_GOAL = all
.PHONY: all image

all: go-linux-bin image
image: go-image-base
