FROM node:18.18.0-alpine3.17 AS builder

ENV NODE_OPTIONS=--openssl-legacy-provider

# make npm install more cache use
# modify business code also could use docker cache untill npm install unless modify dependencies
COPY package.json /workspace/package.json
COPY package-lock.json /workspace/package-lock.json
WORKDIR /workspace
RUN npm install 

COPY . /workspace
RUN npm run build:prod

FROM nginx:stable-alpine

COPY --from=builder /workspace/dist /dist
COPY ./nginx/nginx.conf /etc/nginx
COPY ./nginx/fe.conf /etc/nginx/conf.d/default.conf

COPY ./docker-start.sh /

CMD ["/bin/sh", "/docker-start.sh"]