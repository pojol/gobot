# 安装

### 本地运行
1. 在 release 页下载对应的客户端版本
2. 在 release 页下载对应的服务器版本
3. 启动服务器

### 通过 docker-compose
1. 安装 docker-compose
    ```shell
    # for CentOS
    yum install docker-compose -y

    # for Ubuntu
    apt-get install docker-compose -y
    ```

2. 下载并编辑 [docker-compose.yml](https://github.com/pojol/gobot-driver/blob/develop/docker-compose.yml) 在启动前，请确保对 `MYSQL_ROOT_PASSWORD` 和 `MYSQL_PASSWORD` 参数进行赋值

    ```yaml
    version: "3.7"

    volumes:
    db:

    services:
    db:
        image: mariadb:10.5
        restart: always
        networks:
        - gnet
        volumes:
        - db:/var/lib/mysql
        environment:
        - MYSQL_ROOT_PASSWORD=
        - MYSQL_PASSWORD=
        - MYSQL_DATABASE=gobot
        - MYSQL_USER=gobot

    gobot_driver:
        image: braidgo/gobot-driver:latest
        restart: always
        networks:
        - gnet
        depends_on:
        - db
        ports:
        - 8888:8888
        deploy:
        resources:
            limits:
            cpus: "0.3"
        environment:
        - MYSQL_PASSWORD=
        - MYSQL_DATABASE=gobot
        - MYSQL_USER=gobot
        - MYSQL_HOST=db

    gobot_editor:
        image: braidgo/gobot-editor:latest
        restart: always
        depends_on:
        - gobot_driver
        ports:
        - 7777:7777

    networks:
    gnet:
        driver: bridge
    ```
3. 运行命令 `docker-compose up -d` 运行成功后，访问 http://localhost:7777/ 即可进行 gobot 的编辑