# FROM giovtorres/docker-centos7-slurm:latest
FROM registry.intern.yuansuan.cn/euclid/docker-centos7-slurm:latest

#RUN yum install singularity -y

ENV YS_MODE=dev
ENV YS_LOG_LEVEL=debug
RUN useradd -m yuansuan

WORKDIR /workspace/standard-compute
ENV LOG_LVL="info"

ADD config /workspace/standard-compute/config
ADD standard-compute /workspace/standard-compute/standard-compute

CMD [ "./standard-compute" ]
