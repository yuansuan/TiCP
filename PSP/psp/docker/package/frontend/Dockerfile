FROM centos:7

## 设置时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# WORKDIR
ENV YS_TOP=/opt/yuansuan
WORKDIR $YS_TOP
RUN mkdir -p ${YS_TOP}/psp
RUN mkdir -p ${YS_TOP}/3rd_party

COPY ./dist/fe $YS_TOP/psp/fe
COPY ./docker/package/frontend/nginx.tar.gz $YS_TOP/3rd_party
COPY ./docker/package/frontend/node.tar.xz $YS_TOP/3rd_party
COPY ./cmd/start_front.sh $YS_TOP/psp/


RUN tar -xvf ${YS_TOP}/3rd_party/nginx.tar.gz -C ${YS_TOP}/3rd_party
RUN tar -xvf ${YS_TOP}/3rd_party/node.tar.xz -C ${YS_TOP}/3rd_party
RUN mkdir -p ${YS_TOP}/psp/certs

ENV PATH="$PATH:/opt/yuansuan/3rd_party/node/bin"
ENV PATH="$PATH:/opt/yuansuan/3rd_party/nginx/sbin"

RUN npm config set registry https://registry.npmmirror.com
RUN npm install -g wscat

EXPOSE 443

ENTRYPOINT ["/opt/yuansuan/psp/start_front.sh"]
