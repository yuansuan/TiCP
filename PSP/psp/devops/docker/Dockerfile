FROM centos:7

RUN yum -y update && yum clean all

ENV WORKSPACE workspace

RUN mkdir -p /$WORKSPACE && chmod -R 777 /$WORKSPACE && \
    yum install -y centos-release-scl unzip curl wget gcc mariadb nmap-ncat && \
    yum clean all

ENV GOPATH=/$WORKSPACE

RUN mkdir -p /usr/local

# Build protobuf against configured revision
#
ENV PROTOBUF_REVISION 3.6.1

RUN curl -sLO https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOBUF_REVISION}/protoc-${PROTOBUF_REVISION}-linux-x86_64.zip \
  && unzip protoc-${PROTOBUF_REVISION}-linux-x86_64.zip -d /usr/local \
  && chmod +x /usr/local/bin/protoc \
  && chmod -R 755 /usr/local/include/ \
  && rm -f protoc-${PROTOBUF_REVISION}-linux-x86_64.zip

RUN wget -qO go1.19.linux-amd64.tar.gz https://dl.google.com/go/go1.19.linux-amd64.tar.gz \
    && tar xf go1.19.linux-amd64.tar.gz \
    && mv go /usr/local/ \
    && ln -sf /usr/local/go/bin/* /usr/local/bin/ \
    && rm -f go1.19.linux-amd64.tar.gz

RUN yum install make pam-devel -y

#build protoc-gen-go for gRPC
COPY centos/protoc-gen-go /usr/local/bin/protoc-gen-go
ENV PATH=$GOPATH/bin:/usr/local/bin:$PATH
ENV GOPROXY=https://goproxy.cn,direct
ENV GONOPROXY=*.yuansuan.cn
ENV GONOSUMDB=*.yuansuan.cn
ENV GOPRIVATE=*.yuansuan.cn
ENV CGO_ENABLED=1
ENV GO111MODULE=on
ENV ROOTPATH=/$WORKSPACE
WORKDIR /$WORKSPACE

# node.js stuff

RUN curl -sLO http://nodejs.org/dist/v10.15.3/node-v10.15.3-linux-x64.tar.gz \
  && tar --strip-components 1 -xzvf node-v* -C /usr/local \
  && npm config set registry https://registry.npm.taobao.org \
  && npm install -g typescript \
  && rm -f node-v10.15.3-linux-x64.tar.gz

RUN rpm -Uvh http://nginx.org/packages/centos/7/noarch/RPMS/nginx-release-centos-7-0.el7.ngx.noarch.rpm
RUN yum install openssh-clients nginx -y

COPY nginx.conf yscloud.key yscloud.crt client.crt /etc/nginx/
COPY frontend.conf  /etc/nginx/conf.d/
RUN yum -y install pam-devel git  ruby-devel gcc make rpm-build rubygems gcc-c++ make 
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
  && echo 'Asia/Shanghai' >/etc/timezone

COPY centos/protoc-go-inject-tag /usr/local/bin/protoc-go-inject-tag
COPY centos/protoc-gen-grpc-gateway /usr/local/bin/protoc-gen-grpc-gateway