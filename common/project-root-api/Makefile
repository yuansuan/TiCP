override DOCKER_REGISTRY ?= harbor.yuansuan.cn/ticp
GO := go
ROOT_PACKAGE := github.com/yuansuan/ticp/common/project-root-api
ifeq ($(origin ROOT_DIR),undefined)
ROOT_DIR := $(shell pwd)
endif

OUTPUT_DIR := $(ROOT_DIR)/_output
PROTOC_INC_PATH=$(dir $(shell which protoc 2>/dev/null))/../include
API_DEPS=$(shell find proto -type f -name "*.proto")
API_DEPSRCS=$(API_DEPS:.proto=.pb.go)

# Linux command settings                    
FIND := find . ! -path './third_party/*' ! -path './vendor/*'        
XARGS := xargs --no-run-if-empty        

all: gen


## gen: Generate protobuf files.
.PHONY: gen
gen: gen.proto

.PHONY: gen.plugin.verify
gen.plugin.verify:
ifeq (,$(shell which protoc-gen-go 2>/dev/null))
	@echo "===========> Installing protoc-gen-go"
	@GO111MODULE=on $(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@latest
endif
ifeq (,$(shell which protoc-gen-go-grpc 2>/dev/null))
	@echo "===========> Installing protoc-gen-go-grpc"
	@GO111MODULE=on $(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
endif

$(API_DEPSRCS): gen.plugin.verify $(API_DEPS)
	@echo "===========> Generate protobuf files for $@"
	@mkdir -p $(OUTPUT_DIR)
	@protoc -I $(PROTOC_INC_PATH) -I. \
	 --go_out=$(OUTPUT_DIR) --go-grpc_out=$(OUTPUT_DIR) $(@:.pb.go=.proto)
	-@cp $(OUTPUT_DIR)/$(ROOT_PACKAGE)/$@ $@ || cp $(OUTPUT_DIR)/$@ $@
	-@cp $(OUTPUT_DIR)/$(ROOT_PACKAGE)/$(@:.pb.go=_grpc.pb.go) $(@:.pb.go=_grpc.pb.go) || cp $(OUTPUT_DIR)/$(@:.pb.go=_grpc.pb.go) $(@:.pb.go=_grpc.pb.go) || true
	@rm -rf $(OUTPUT_DIR)


.PHONY: gen.proto
 gen.proto: $(API_DEPSRCS)

.PHONY: gen.clean
gen.clean:
	@echo "===========> Cleaning generated protobuf files for $(API_DEPS)"
	@rm -f $(API_DEPS:.proto=.pb.go) $(API_DEPS:.proto=_grpc.pb.go)

.PHONY: protos
protos:
	@docker pull $(DOCKER_REGISTRY)/euclid-go-buildbuf:v1.0
	docker run -i $(DOCKER_TTY) --rm --entrypoint="/bin/bash" \
		-v $(PWD):/workspace -v $(HOME)/.cache/buf:/root/.cache/buf \
		$(DOCKER_REGISTRY)/euclid-go-buildbuf:v1.0 -c 'buf generate'