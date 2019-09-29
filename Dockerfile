FROM golang:alpine AS binarybuilder
WORKDIR /qiongbi
ARG GATEWAY=openapi.alipay.com/gateway.do
COPY . .
RUN apk --no-cache --no-progress add --virtual build-deps build-base git linux-pam-dev \
    && go mod tidy -v \
    && go mod vendor \
    # 修改支付宝网关
    && grep openapi.alipay.com -rl ./vendor/github.com/smartwalle/alipay|xargs sed -i "s~openapi.alipay.com/gateway.do~${GATEWAY}~g" \
    && grep ${GATEWAY} -rl ./vendor/github.com/smartwalle/alipay \
    && go build -mod=vendor -o app -ldflags="-s -w" cmd/web/main.go
FROM alpine:latest
RUN apk --no-cache --no-progress add \
    ca-certificates \
    tzdata
WORKDIR /qiongbi
COPY resource /qiongbi/resource
COPY --from=binarybuilder /qiongbi/app ./app

ENV AppID=${AppID} \
    PubKey=${PubKey} \
    PriKey=${PriKey} \
    Domain=${Domain}

VOLUME ["/qiongbi/data"]
EXPOSE 8080
CMD ["/qiongbi/app"]