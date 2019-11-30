# 穷逼 ![构建状态](https://github.com/naiba/qiongbi/workflows/Build%20Private%20Image/badge.svg)

:neckbeard: 穷逼捐赠系统，使用支付宝即时到账接口，前端来自另一个[穷逼](https://qiong.bi)。

## 食用指南

1. 拉取镜像

    ```shell
    docker pull docker.pkg.github.com/naiba/dockerfiles/qiongbi
    ```

2. 使用 `docker-compose.yml` 方式运行，方便

    ```yaml
    version: '3.3'

    services:
    db:
        image: docker.pkg.github.com/naiba/dockerfiles/qiongbi:latest
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

3. 启动 `docker-compose up -d`，可以开始接受捐赠啦
