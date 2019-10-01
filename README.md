# 穷逼 ![构建状态](https://github.com/naiba/qiongbi/workflows/.github/workflows/docker.yml/badge.svg)

:neckbeard: 穷逼捐赠系统，使用支付宝即时到账接口，前端来自另一个[穷逼](https://qiong.bi)。

## 食用指南

1. 首先构建镜像

    ```shell
    docker build -t docker.pkg.github.com/naiba/qiongbi/web .
    ```

2. 使用 `docker-compose.yml` 方式运行，方便

    ```yaml
    version: '3.3'

    services:
    db:
        image: docker.pkg.github.com/naiba/qiongbi/web:latest
        volumes:
        - ./data/:/qiongbi/data
        restart: always
        ports:
        - "172.17.0.1:8001:8080"
        environment:
        - AppID=qiongbi # 支付宝应用AppID
        - PubKey=xxxxxx # 支付宝公钥
        - PriKey=xxxxxx # 商户私钥
        - Domain=example.com # 绑定的域名
    ```

这边注意不要使用本仓库中构建的 image，本仓库中构建的 image 是我的支付宝网关专用的。
