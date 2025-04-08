LDFLAGS += -X "main.BuildDatetime=$(shell TZ='Asia/Shanghai' date '+%Y-%m-%d %I:%M:%S')"
LDFLAGS += -X "main.GoVersion=$(shell go version)"
LDFLAGS += -X "main.GitBranch=$(shell git rev-parse --abbrev-ref HEAD)"
LDFLAGS += -X "main.GitHash=$(shell git rev-parse HEAD)"
LDFLAGS += -X "main.GitName=$(shell git config --get user.name)"

REMOTE_ADDRESS=harbor.yuansuan.cn
IMAGE_VERSION=$(shell date +%Y%m%d%H%M)
BUILD_BE_IMAGE=$(REMOTE_ADDRESS)/ticp/psp-be:$(IMAGE_VERSION)
BUILD_FE_IMAGE=$(REMOTE_ADDRESS)/ticp/psp-fe:$(IMAGE_VERSION)

.PHONY: all-be all-fe

all-be: go-linux-bin build-psp-be-image
all-fe: yarn-build-desktop build-psp-fe-image

go-linux-bin:
	@rm -f cmd/pspd
	@echo "🚧 正在编译..."
	@start=$$(date +%s); GOOS="linux" GOARCH="amd64" CGO_ENABLED="0" go build -x -v -ldflags '$(LDFLAGS)' -gcflags "all=-N -l" -o cmd/pspd cmd/main.go; end=$$(date +%s); \
	runtime=$$(expr $$end - $$start); echo "✅ 编译完成，用时 $$runtime 秒"

build-psp-be-image:
	@echo "🚧 正在生成镜像..."
	@start=$$(date +%s); docker build -t $(BUILD_BE_IMAGE) -f ./docker/package/psp/Dockerfile .; end=$$(date +%s); runtime=$$(expr $$end - $$start); echo "✅ 生成镜像完成，用时 $$runtime 秒"

yarn-build-desktop:
	@rm -rf ../../node_modules
	@cd web/desktop && yarn install && yarn build

build-psp-fe-image:
	@echo "🚧 正在生成镜像..."
	@cd web/desktop && start=$$(date +%s); docker build -t $(BUILD_FE_IMAGE) .; end=$$(date +%s); runtime=$$(expr $$end - $$start); echo "✅ 生成镜像完成，用时 $$runtime 秒"
