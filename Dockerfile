FROM golang:alpine
WORKDIR /src
ADD . /src/
RUN chmod +x /src/docker-entrypoint.sh

RUN sed -i 's/dl-cdn.alpinelinux.org/repo.huaweicloud.com/g' /etc/apk/repositories && \
    apk update && apk add nodejs nodejs-npm

RUN npm config set registry https://registry.npm.taobao.org && cd /src/webui && npm install node-sass \
    --registry=https://registry.npm.taobao.org --sass_binary_site=https://npm.taobao.org/mirrors/node-sass

RUN cd /src/webui && npm install --registry=https://registry.npm.taobao.org \
    --sass_binary_site=https://npm.taobao.org/mirrors/node-sass && \
    npm run build && rm -rf node_modules

RUN cd /src && GO111MODULE=on GOPROXY=https://goproxy.cn GOOS=linux GOARCH=amd64 go build -tags netgo -ldflags="-s -w" -o xpb .

EXPOSE 21330

ENTRYPOINT ["/src/docker-entrypoint.sh"]
