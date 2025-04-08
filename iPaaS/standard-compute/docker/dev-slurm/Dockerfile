#FROM giovtorres/docker-centos7-slurm:latest
FROM registry.intern.yuansuan.cn/euclid/docker-centos7-slurm:latest

RUN wget https://go.dev/dl/go1.18.linux-amd64.tar.gz && rm -rf /usr/local/go && tar -C /usr/local -xzf go1.18.linux-amd64.tar.gz


ENV GOPROXY=https://goproxy.cn
ENV YS_MODE=dev
ENV YS_LOG_LEVEL=debug
ENV GOPATH=/root/gopath2
ENV PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN yum install singularity -y

ENV GOPATH=/root/go

WORKDIR /workspace