FROM nginx:stable
ARG GIT_SHA1=""
ENV VERSION=$GIT_SHA1
ENV YS_TOP=/opt/yuansuan
RUN mkdir $YS_TOP
WORKDIR $YS_TOP
COPY ./dist $YS_TOP

RUN rm -f /etc/nginx/conf.d/default.conf
COPY ./frontend.conf /etc/nginx/conf.d/default.conf
