# gobot editor
Gobot是一个有状态的api测试工具，支持图形编辑、api调用、绑定脚本执行。

[![](https://img.shields.io/badge/Trello-Todo-2ca5e0?style=flat&logo=trello)](https://trello.com/b/8eDZ6h7n/)
[![CI](https://github.com/pojol/gobot/actions/workflows/dockerimage.yml/badge.svg?branch=develop)](https://github.com/pojol/gobot/actions/workflows/dockerimage.yml)

# 工具的目标是什么？
1. 使用 bot 进行复杂逻辑（有状态）的测试
    * 如游戏业务中 创建角色→发送邮件→使用道具→进行战斗 ...
    * 如社交业务中 创建多个角色 → 互相发现|添加好友 → 点赞|评论 ...
2. 尽可能的简单，只需要拖拽行为树节点 + 少量的脚本编辑

# 特性
* 使用 behavior tree 编排 bot 的运行逻辑 
* 使用 lua script 控制 bot 的执行逻辑
* 每个 bot 都拥有一个 meta 数据结构，用于存储整个测试流程的上下文
* 使用 tag + filter 管理 bot 行为文件
* 拥有直观的调试窗口和环境，可以单步查看节点逻辑的执行情况

# [在线试用](http://1.117.168.37:7777) <--
# [文档](https://pojol.gitee.io/gobot/#/) <--


## 编辑器预览
[![image.png](https://i.postimg.cc/mrbSNKmS/image.png)](https://postimg.cc/CRQDwrTZ)

## 脚本接口
* [http](https://docs.gobot.fun/#/zh-cn/advance/script_http)
* [proto](https://docs.gobot.fun/#/zh-cn/advance/script_protobuf)
* [utils](https://docs.gobot.fun/#/zh-cn/advance/script_utils)
* [base64](https://docs.gobot.fun/#/zh-cn/advance/script_base64)
* [json](https://docs.gobot.fun/#/zh-cn/advance/script_utils)


## Http请求例子
```lua
-- lua script
local http = require("http")

reqTable = {
    body = {},       -- 消息内容
    timeout = "10s", -- http 请求超时时间
    headers = {},    -- http 消息头
}

-- .post .put .get
res, err = http.post("url", reqTable)

--[[
    res                 -- userdata
    res["body"]         -- http 回复内容
    res["body_size"]    -- 回复内容大小
    res["headers"]      -- http 消息头
    res["cookies"]      -- http cookies
    res["status_code"]  -- http 状态码
    res["url"]          -- 请求地址

    err                 -- 错误信息
]]--
```

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
