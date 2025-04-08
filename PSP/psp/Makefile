
PROJECT_ROOT := $(CURDIR)

LDFLAGS += -X "main.BuildDatetime=$(shell TZ='Asia/Shanghai' date '+%Y-%m-%d %I:%M:%S')"
LDFLAGS += -X "main.GoVersion=$(shell go version)"
LDFLAGS += -X "main.GitBranch=$(shell git rev-parse --abbrev-ref HEAD)"
LDFLAGS += -X "main.GitHash=$(shell git rev-parse HEAD)"
LDFLAGS += -X "main.GitName=$(shell git config --get user.name)"

image-push:
	cd $(PROJECT_ROOT)/cmd && make docker-image-build-push

proto:
	@cd cmd && make proto

psp: doc
	@cd cmd && make docker-build-psp

deploy:
	@bash ./deploy.sh

local:
	@rm -f cmd/pspd
	@echo "🚧 正在编译..."
	@start=$$(date +%s); GOOS="linux" GOARCH="amd64" CGO_ENABLED="0" go build -x -v -ldflags '$(LDFLAGS)' -gcflags "all=-N -l" -o cmd/pspd cmd/main.go; end=$$(date +%s); \
	runtime=$$(expr $$end - $$start); echo "✅ 编译完成，用时 $$runtime 秒"

local-agent:
	@rm -f agent/psp_agent
	@echo "🚧 正在编译..."
	@start=$$(date +%s); GOOS="linux" GOARCH="amd64" CGO_ENABLED="0" go build -x -v -ldflags '$(LDFLAGS)' -gcflags "all=-N -l" -o agent/psp_agent agent/main.go; end=$$(date +%s); \
	runtime=$$(expr $$end - $$start); echo "✅ 编译完成，用时 $$runtime 秒"

agent:
	@cd cmd/agent && make docker-build-agent

doc:
	@go install github.com/swaggo/swag/cmd/swag@latest
	@cd cmd && swag fmt && cd ../pkg && cd ../internal && swag fmt
	@cd $(PROJECT_ROOT)/cmd && swag init -d .,.. --exclude devops,web -p snakecase -ot go

check-cycle:
	@go install github.com/elza2/go-cyclic@latest
	@go-cyclic run --dir .


IMAGE_TAG=v1.6

build-psp-be-image:
	docker build -t psp-be:$(IMAGE_TAG) -f ./docker/package/psp/Dockerfile .

export-psp-be-image: build-psp-be-image
	docker save -o ./docker/package/psp-be-image.tar psp-be:$(IMAGE_TAG)

rm-psp-be-image:
	-@docker stop psp-be
	-@docker rm psp-be
	-@docker rmi psp-be:$(IMAGE_TAG)


build-psp-fe-image:
	docker build --network=host --no-cache  -t psp-fe:$(IMAGE_TAG) -f ./docker/package/frontend/Dockerfile .

export-psp-fe-image: build-psp-fe-image
	docker save -o ./docker/package/psp-fe-image.tar psp-fe:$(IMAGE_TAG)

rm-psp-fe-image:
	-@docker stop psp-fe
	-@docker rm psp-fe
	-@docker rmi psp-fe:$(IMAGE_TAG)

gofmt:
	@go install -v github.com/incu6us/goimports-reviser/v3@latest
	@time find $(CURDIR)/internal -not -path "$(CURDIR)/internal/common/proto/*" -name "*.go" -exec goimports-reviser -company-prefixes yuansuan.cn/psp -company-prefixes yuansuan.cn -rm-unused  {} \;
