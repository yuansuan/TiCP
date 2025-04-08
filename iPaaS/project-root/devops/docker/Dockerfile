FROM harbor.yuansuan.cn/ticp/alpine:3

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache tzdata
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" >  /etc/timezone
RUN apk update
RUN apk add --no-cache curl shadow

RUN useradd -m yuansuan

ARG MODULENAME

WORKDIR /workspace/
ENV LOG_LVL="info"

ADD $MODULENAME /workspace/$MODULENAME/

ENV COMMAND=/workspace/$MODULENAME/$MODULENAME
WORKDIR /workspace/$MODULENAME

ENTRYPOINT ["/bin/sh", "-c", "$COMMAND"]