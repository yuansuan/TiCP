FROM 0.0.0.0:5555/mobigo:v1.6

RUN date
ADD ../.. /workspace

RUN date
RUN sh /workspace/devops/onpremise/build.sh && sh /workspace/devops/clean.sh
RUN date


FROM centos:7
WORKDIR /workspace/
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
  && echo 'Asia/Shanghai' >/etc/timezone

RUN curl -sLO http://nodejs.org/dist/v14.18.3/node-v14.18.3-linux-x64.tar.gz \
  && tar --strip-components 1 -xzvf node-v* -C /usr/local \
  && rm -f node-v14.18.3-linux-x64.tar.gz

COPY --from=0 /etc/nginx /etc/nginx
RUN rpm -Uvh http://nginx.org/packages/centos/7/noarch/RPMS/nginx-release-centos-7-0.el7.ngx.noarch.rpm
RUN yum install nginx mariadb nmap-ncat -y

RUN date
COPY --from=0 /workspace/devops /workspace/devops
COPY --from=0 /workspace/docker /workspace/docker
COPY --from=0 /workspace/dist /workspace/dist
COPY --from=0 /workspace/src/yuansuan.cn /workspace/src/yuansuan.cn
RUN date
