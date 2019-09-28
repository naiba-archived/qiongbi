FROM golang:alpine AS binarybuilder
WORKDIR /qiongbi
COPY . .
RUN cd cmd/web \
    && go build -o app -ldflags="-s -w"
FROM alpine:latest
RUN apk --no-cache --no-progress add \
    ca-certificates \
    tzdata
WORKDIR /qiongbi
COPY resource /qiongbi/resource
COPY --from=binarybuilder /qiongbi/cmd/web/app ./app

ENV AppID=$AppID \
    PubKey=$PubKey \
    PriKey=$PriKey \
    Domain=$Domain

VOLUME ["/qiongbi/data"]
EXPOSE 8080
CMD ["/qiongbi/app"]