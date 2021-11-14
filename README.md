# gobot
Robot editor framework based on behavior tree, node supports binding script execution.


[![](https://img.shields.io/badge/%E6%96%87%E6%A1%A3-Doc-2ca5e0?style=flat&logo=github)](https://docs.gobot.fun/)
[![](https://img.shields.io/badge/Trello-Todo-2ca5e0?style=flat&logo=trello)](https://trello.com/b/8eDZ6h7n/)


# Try it out
Try the editor out [on website](http://1.117.168.37:7777/)

# Feature
* Easily edit robot behavior
    * Write stateful test cases
    * Write timed tasks
    * ...
* Lua script node
    * Support http module
    * Support protobuf
* Easily manage robot runtime status information (Metadata

## Preview
[![image.png](https://i.postimg.cc/RCdYzPHh/image.png)](https://postimg.cc/DmW6rP53)

# Install
1. Install docker-compose
    ```shell
    # for CentOS
    yum install docker-compose -y

    # for Ubuntu
    apt-get install docker-compose -y
    ```

2. Down load and modify [docker-compose.yml](https://github.com/pojol/gobot-driver/blob/develop/docker-compose.yml) file make sure to pass in values for `MYSQL_ROOT_PASSWORD` and `MYSQL_PASSWORD` variables before you run this setup.

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
3. Run `docker-compose up -d`, now you can access gobot at http://localhost:7777/ from your host system.