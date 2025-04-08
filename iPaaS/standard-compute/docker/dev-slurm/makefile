all: image push

image: 
	docker build . -t registry.intern.yuansuan.cn/euclid/standard-compute-dev-slurm

push:
	docker push registry.intern.yuansuan.cn/euclid/standard-compute-dev-slurm
