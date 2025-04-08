override DOCKER_REGISTRY ?= harbor.yuansuan.cn/ticp
override DOCKER_TTY ?= -t

.PHONY: protos
protos:
	@docker pull $(DOCKER_REGISTRY)/euclid-go-buildbuf:v1.0
ifndef DOCKER_HOST
	docker run -i $(DOCKER_TTY) --rm --entrypoint="/bin/bash" \
		-v $(PWD):/workspace -v $(HOME)/.cache/buf:/root/.cache/buf \
		$(DOCKER_REGISTRY)/euclid-go-buildbuf:v1.0 -c 'buf generate'
else
	@$(eval VOLUME_NAME=$(shell bash -c 'echo $$RANDOM'))
	@docker volume create $(VOLUME_NAME)
	@docker container create --name ${VOLUME_NAME} -v ${VOLUME_NAME}:/data busybox
	@docker cp ./protos ${VOLUME_NAME}:/data/protos
	@docker cp ./buf.gen.yaml ${VOLUME_NAME}:/data
	@docker cp ./buf.lock ${VOLUME_NAME}:/data
	@docker cp ./buf.yaml ${VOLUME_NAME}:/data
	docker run -i $(DOCKER_TTY) -v $(VOLUME_NAME):/data -v /UserHome/hanna/.cache/buf:/root/.cache/buf --rm --entrypoint="/bin/bash" \
		$(DOCKER_REGISTRY)/euclid-go-buildbuf:v1.0 -ec 'cd /data && buf generate'
	@docker cp ${VOLUME_NAME}:/data/protos .
	@docker rm ${VOLUME_NAME}
	@docker volume rm ${VOLUME_NAME}
endif
