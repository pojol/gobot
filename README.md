# gobot
Gobot is a stateful api testing tool that supports graph editing, api calling, and binding script execution.

[![](https://img.shields.io/badge/%E6%96%87%E6%A1%A3-Doc-2ca5e0?style=flat&logo=github)](https://docs.gobot.fun/)
[![](https://img.shields.io/badge/Trello-Todo-2ca5e0?style=flat&logo=trello)](https://trello.com/b/8eDZ6h7n/)
[![CI](https://github.com/pojol/gobot/actions/workflows/dockerimage.yml/badge.svg?branch=develop)](https://github.com/pojol/gobot/actions/workflows/dockerimage.yml)

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

[Document](https://docs.gobot.fun)

## Preview
[![2022-04-20-9-51-09.png](https://i.postimg.cc/xCW3KnxD/2022-04-20-9-51-09.png)](https://postimg.cc/bD9nPcf3)

## Benchmark
[![gobot-qps.png](https://i.postimg.cc/5y72F2Nb/gobot-qps.png)](https://postimg.cc/WqZvBjkH)

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