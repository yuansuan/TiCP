FROM centos:7

ENV WORKSPACE workspace

RUN mkdir -p /$WORKSPACE \
    && chmod -R 777 /$WORKSPACE \
    && yum install -y make git wget gcc \
    && yum -y update \
    && yum clean all

RUN mkdir -p /usr/local \
    && wget -qO go1.19.linux-amd64.tar.gz https://dl.google.com/go/go1.19.linux-amd64.tar.gz \
    && tar xf go1.19.linux-amd64.tar.gz \
    && mv go /usr/local/ \
    && ln -sf /usr/local/go/bin/* /usr/local/bin/ \
    && rm -f go1.19.linux-amd64.tar.gz

ENV PATH=/$WORKSPACE/bin:/usr/local/bin:$PATH
ENV GOPROXY=https://goproxy.cn,direct
ENV GONOPROXY=*.yuansuan.cn
ENV GONOSUMDB=*.yuansuan.cn
ENV GOPRIVATE=*.yuansuan.cn
ENV CGO_ENABLED=1
ENV GO111MODULE=on
ENV ROOTPATH=/$WORKSPACE

WORKDIR /$WORKSPACE
