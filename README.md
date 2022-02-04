# gobot
Gobot is a stateful api testing tool that supports graph editing, api calling, and binding script execution.

[![](https://img.shields.io/badge/%E6%96%87%E6%A1%A3-Doc-2ca5e0?style=flat&logo=github)](https://docs.gobot.fun/)
[![](https://img.shields.io/badge/Trello-Todo-2ca5e0?style=flat&logo=trello)](https://trello.com/b/8eDZ6h7n/)

[中文](https://github.com/pojol/gobot/blob/master/README_CN.md)

# What is the goal of the tool ?
1. Use bots for complex logic (stateful) testing
    * For example, in the game business, create a character → send an email → use items → fight ...
    * For example, create multiple roles in social business → discover each other | add friends → like | comment ...
2. Keep it simple

# Feature
* Use behavior tree to arrange bot's running logic
* Use lua script to control the execution logic of bot
* Each bot has a meta data structure to store the context of the entire test process
* Use tag + filter to manage bot behavior files
* With an intuitive debugging window and environment, you can view the execution of the node logic in a single step

# Try it out
Try the editor out [on website](http://1.117.168.37:7777/)

## Preview
[![image.png](https://i.postimg.cc/mrbSNKmS/image.png)](https://postimg.cc/CRQDwrTZ)


## Script interface
* [http](https://docs.gobot.fun/#/zh-cn/advance/script_http)
* [proto](https://docs.gobot.fun/#/zh-cn/advance/script_protobuf)
* [utils](https://docs.gobot.fun/#/zh-cn/advance/script_utils)
* [base64](https://docs.gobot.fun/#/zh-cn/advance/script_base64)
* [json](https://docs.gobot.fun/#/zh-cn/advance/script_utils)


## Http request sample
```lua
-- lua script
local http = require("http")

reqTable = {
    body = {},       -- post body
    timeout = "10s", -- http timeout
    headers = {},    -- http headers
}

-- .post .put .get
res, err = http.post("url", reqTable)

--[[
    res                 -- userdata
    res["body"]         -- http response body
    res["body_size"]    -- body size
    res["headers"]      -- http headers
    res["cookies"]      -- http cookies
    res["status_code"]  -- http status code
    res["url"]          -- request url

    err                 -- error message
]]--
```

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