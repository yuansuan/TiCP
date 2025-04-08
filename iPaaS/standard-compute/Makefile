override DOCKER_REGISTRY ?= registry.intern.yuansuan.cn/euclid
override DOCKER_TTY ?= -t
override DEVOPS_VERSION := $(shell cat tools/devops/VERSION)
override BUILD_TIME := $(shell date '+%Y-%m-%dT%H:%M:%S')
override GIT_REVISION := $(shell git rev-parse --short HEAD)

all: linux-bin

linux-bin:
	GOOS=linux go build -ldflags "-X github.com/yuansuan/ticp/iPaaS/standard-compute/version.gitCommit=$(shell git rev-parse HEAD)" \
 		-o standard-compute github.com/yuansuan/ticp/iPaaS/standard-compute/cmd

docs-image:
	docker run --rm -it --user $(shell id -u):$(shell id -g) -e GOPATH=$(shell go env GOPATH):/go -e HOME=$(HOME) -v $(HOME):$(HOME) -w $(PWD) quay.io/goswagger/swagger generate spec -o swagger.json
	docker build . -f docs/docker/Dockerfile -t registry.intern.yuansuan.cn/euclid/standard-compute-docs

singularity-devops:
	docker build -f docker/singularity-devops/Dockerfile . -t $(DOCKER_REGISTRY)/standard-compute-devops

devops:
ifdef BUILD_ENV
	go build -o devops -tags embedded \
		-ldflags "-X 'github.com/yuansuan/ticp/iPaaS/standard-compute/tools/devops/config.env=$(BUILD_ENV)' \
			-X 'github.com/yuansuan/ticp/iPaaS/standard-compute/tools/devops.BuildTime=$(BUILD_TIME)' \
			-X 'github.com/yuansuan/ticp/iPaaS/standard-compute/tools/devops.Version=$(DEVOPS_VERSION)' \
			-X 'github.com/yuansuan/ticp/iPaaS/standard-compute/tools/devops.Revision=$(GIT_REVISION)'" \
		github.com/yuansuan/ticp/iPaaS/standard-compute/tools/devops/cmd
else
	go build -o devops github.com/yuansuan/ticp/iPaaS/standard-compute/tools/devops/cmd
endif


package:
	docker run -it --rm \
		-v ~/.ssh:/root/.ssh \
		-v $(PWD)/rpm/SPECS:/root/rpmbuild/SPECS \
		-v $(PWD):/root/rpmbuild/BUILD/standard-compute \
		-v $(PWD)/rpm/RPMS:/root/rpmbuild/RPMS \
		registry.intern.yuansuan.cn/euclid/standard-compute-rpm-maker rpmbuild -bb /root/rpmbuild/SPECS/build.spec

clean:
	rm -f devops
	rm -f standard-compute
