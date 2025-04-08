all: image run push

image:
	docker build . -t registry.intern.yuansuan.cn/euclid/docker-centos7-pbspro

run:
	docker run --rm  -it --hostname pbspro registry.intern.yuansuan.cn/euclid/docker-centos7-pbspro:latest


push:
	docker push registry.intern.yuansuan.cn/euclid/docker-centos7-pbspro
