# gobot editor
基于行为树的机器人编辑器框架，节点支持绑定脚本执行。


[![](https://img.shields.io/badge/%E6%96%87%E6%A1%A3-Doc-2ca5e0?style=flat&logo=github)](https://docs.gobot.fun/)
[![](https://img.shields.io/badge/Trello-Todo-2ca5e0?style=flat&logo=trello)](https://trello.com/b/8eDZ6h7n/)


# 在线试用
尝试在[网站](http://1.117.168.37:7777/)上进行编辑 

# 特性
* 方便的机器人行为编辑
* 方便的机器人管理
* 支持 Lua 脚本节点
* 方便的调试机器人行为

## 编辑器预览
[![image.png](https://i.postimg.cc/FFcBHwg5/image.png)](https://postimg.cc/5j43P7Rn)

## 脚本接口
||||||||
|-|-|-|-|-|-|-|
|[http](https://docs.gobot.fun/#/zh-cn/advance/script).get|http.post|http.put|[proto]().marshal|proto.unmarshal|[utils]().random|utils.uuid|
|[base64]().encode|base64.decode|[meta]()|[merge]()|[json]().encode|json.decode|more ...|

# 安装
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
