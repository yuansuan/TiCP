# FROM giovtorres/docker-centos7-slurm:latest
FROM registry.intern.yuansuan.cn/euclid/docker-centos7-slurm:latest

ADD config /workspace/config
ADD standard-compute /workspace/standard-compute

ENV YS_MODE=local
ENV YS_LOG_LEVEL=debug

RUN yum install singularity -y

RUN mkdir /root/.ssh

WORKDIR /workspace

CMD [ "./standard-compute" ]