FROM centos:7

RUN yum -y install rpm-build rpmdevtools make git wget
RUN rpmdev-setuptree

RUN wget https://go.dev/dl/go1.18.1.linux-amd64.tar.gz && tar -C /usr/local -xzf go1.18.1.linux-amd64.tar.gz && rm go1.18.1.linux-amd64.tar.gz

ENV PATH $PATH:/usr/local/go/bin
ENV GOPROXY https://goproxy.cn
ENV GO111MODULE on
ENV CGO_ENABLED 0
