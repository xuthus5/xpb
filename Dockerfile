FROM node:14.16.0-alpine AS webui-builder
WORKDIR /builder
COPY ./webui/ /builder/
RUN npm config set registry https://registry.npm.taobao.org
RUN npm install node-sass --registry=https://registry.npm.taobao.org --sass_binary_site=https://npm.taobao.org/mirrors/node-sass
RUN npm install --registry=https://registry.npm.taobao.org --sass_binary_site=https://npm.taobao.org/mirrors/node-sass
RUN npm run build

FROM golang:1.16.2-alpine AS serv-builder
WORKDIR /src
COPY ./cmd /src/cmd
COPY ./common /src/common
COPY ./config /src/config
COPY ./logger /src/logger
COPY ./server /src/server
COPY ./main.go ./go.mod ./go.sum /src/
ENV CGO_ENABLED=0
RUN GO111MODULE=on GOPROXY=https://goproxy.cn GOOS=linux GOARCH=amd64 go build -tags netgo -ldflags="-s -w" -o xpb .

FROM alpine:latest
WORKDIR /app
COPY ./docker-entrypoint.sh .
RUN sed -i "s/dl-cdn.alpinelinux.org/repo.huaweicloud.com/g" /etc/apk/repositories && \
    apk --no-cache add ca-certificates && update-ca-certificates
RUN mkdir webui
COPY --from=webui-builder /builder/dist ./webui/dist
COPY --from=serv-builder /src/xpb .
RUN chmod +x ./docker-entrypoint.sh
EXPOSE 21330
ENTRYPOINT ["./docker-entrypoint.sh"]
