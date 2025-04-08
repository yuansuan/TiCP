#!/bin/bash

mkdir -p /usr/local/nginx/run
mkdir -p /usr/local/nginx/client_temp
mkdir -p /usr/local/nginx/proxy_temp
mkdir -p /var/log/nginx
touch /var/log/nginx/access.log
mkdir -p /opt/yuansuan/psp/work/nginx
touch /opt/yuansuan/psp/work/nginx/nginx_error.log
mkdir -p /opt/yuansuan/psp/log
touch /opt/yuansuan/psp/log/nginx_error.log

cd /opt/yuansuan/psp
nginx -c /opt/yuansuan/psp/config/nginx/nginx.conf -g "daemon off;"
