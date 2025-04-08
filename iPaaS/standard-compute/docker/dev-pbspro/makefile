all: image run push


image:
	docker build . -t registry.intern.yuansuan.cn/euclid/standard-compute-dev-pbspro


run:
	docker run --rm --hostname pbspro registry.intern.yuansuan.cn/euclid/standard-compute-dev-pbspro:latest

push:
	docker push registry.intern.yuansuan.cn/euclid/standard-compute-dev-pbspro
