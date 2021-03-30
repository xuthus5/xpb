FROM golang:alpine
WORKDIR /src
ADD . /src/
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --update --no-cache nodejs nodejs-npm

RUN npm config set registry https://registry.npm.taobao.org && npm cache clean --force && cd webui && npm install && npm run build

RUN GO111MODULE=on GOPROXY=https://goproxy.cn GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -o xpb .

EXPOSE 21330

ENTRYPOINT ["./docker-entrypoint.sh"]
