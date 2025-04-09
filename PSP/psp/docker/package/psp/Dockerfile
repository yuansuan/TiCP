FROM registry.intern.yuansuan.cn/library/alpine:3.10

## 设置时区
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache tzdata
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" >  /etc/timezone
RUN apk update

# WORKDIR
ENV YS_TOP=/opt/yuansuan
ENV CONFIG_PATH=${YS_TOP}/psp/config
WORKDIR $YS_TOP

RUN mkdir -p ${CONFIG_PATH}
RUN mkdir -p ${YS_TOP}/psp/bin
RUN chmod -R 777 ${YS_TOP}

# copy executable file & config
COPY ./cmd/pspd $YS_TOP/psp/bin
COPY ./cmd/.env $YS_TOP/psp/

EXPOSE 22340
EXPOSE 32432

ENTRYPOINT ["/opt/yuansuan/psp/bin/pspd"]
