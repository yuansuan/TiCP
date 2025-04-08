FROM registry.intern.yuansuan.cn/library/alpine:3.10

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache tzdata
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" >  /etc/timezone

ARG BUILDPATH
ARG MODULENAME

WORKDIR /workspace/
ENV LOG_LVL="info"

COPY internal/iamserver/$MODULENAME $BUILDPATH/
COPY internal/iamserver/config $BUILDPATH/config
RUN mkdir -p $BUILDPATH/log

ENV COMMAND $BUILDPATH/$MODULENAME
WORKDIR $BUILDPATH

ENTRYPOINT $COMMAND
