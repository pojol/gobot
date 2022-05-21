# gobot
Gobot is a stateful api testing tool that supports graph editing, api calling, and binding script execution.

[![](https://img.shields.io/badge/%E6%96%87%E6%A1%A3-Doc-2ca5e0?style=flat&logo=github)](https://docs.gobot.fun/)
[![](https://img.shields.io/badge/Trello-Todo-2ca5e0?style=flat&logo=trello)](https://trello.com/b/8eDZ6h7n/)
[![CI](https://github.com/pojol/gobot/actions/workflows/dockerimage.yml/badge.svg?branch=develop)](https://github.com/pojol/gobot/actions/workflows/dockerimage.yml)

[中文](https://github.com/pojol/gobot/blob/master/README_CN.md)

# What is the goal of the tool ?
1. With a complete life cycle, we can retain and use various state information on different nodes (tests for different purposes can be customized by processing the context
2. No need to write control code logic, only need to do logic arrangement and write node scripts (you can use prefab templates
3. Through the API provided by the driver, we can use the test robot in the ci/cd pipeline

# Feature
* Use `behavior tree` + `script` to control the execution logic of the robot
* Support interface editing and debugging
* Manage and find bots in your warehouse with `tag` + `filter`
* Can perform `stress tests` (concurrent driving robots
* The driver provides api calls
* Provides a report viewing page (api statistics called by the robot

# Try it out
Try the editor out [on website](http://1.117.168.37:7777/)


### Preview
[![image.png](https://i.postimg.cc/LXCt5Zcd/image.png)](https://postimg.cc/ZBNBD0Yj)

### Http request sample
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

### Script Module
|||||||
|-|-|-|-|-|-|
|`base64`|`http`|`protobuf`|`utils`|`mongoDB`|`json`|

### Report
[![image.png](https://i.postimg.cc/4d3TTrvf/image.png)](https://postimg.cc/yJ2Gmprt)