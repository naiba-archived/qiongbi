FROM golang:alpine AS binarybuilder
WORKDIR /qiongbi
ARG GATEWAY=openapi.alipay.com/gateway.do
COPY . .
RUN apk --no-cache --no-progress add git\
    && go mod vendor \
    && go mod tidy \
    && sed -i "s/openapi.alipay.com\/gateway.do/${GATEWAY}/g" grep openapi.alipay.com -rl ./vendor/github.com/smartwalle/alipay \
    && cd cmd/web \
    && go build -o app -ldflags="-s -w"
FROM alpine:latest
RUN apk --no-cache --no-progress add \
    ca-certificates \
    tzdata
WORKDIR /qiongbi
COPY resource /qiongbi/resource
COPY --from=binarybuilder /qiongbi/cmd/web/app ./app

ENV AppID=${AppID} \
    PubKey=${PubKey} \
    PriKey=${PriKey} \
    Domain=${Domain}

VOLUME ["/qiongbi/data"]
EXPOSE 8080
CMD ["/qiongbi/app"]