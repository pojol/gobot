# gobot
Gobot is a powerful stateful API testing robot. It provides a graphical interface for building test scenarios, allows for easy test script writing, step-by-step debugging and pressure testing, and can share and store states between each stage of the testing process. 

[![Go Report Card](https://goreportcard.com/badge/github.com/pojol/gobot)](https://goreportcard.com/report/github.com/pojol/gobot)
[![](https://img.shields.io/badge/%E6%96%87%E6%A1%A3-Doc-2ca5e0?style=flat&logo=github)](https://pojol.gitee.io/gobot/#/)
[![](https://img.shields.io/badge/Trello-Todo-2ca5e0?style=flat&logo=trello)](https://trello.com/b/8eDZ6h7n/)
[![CI](https://github.com/pojol/gobot/actions/workflows/dockerimage.yml/badge.svg?branch=develop)](https://github.com/pojol/gobot/actions/workflows/dockerimage.yml)

[中文](https://github.com/pojol/gobot/blob/master/README_CN.md)


## Quick Installation
> Note: In local running mode, all changes are recorded in memory (not actually saved). To save, please download the files locally; or use the official deployment method.

1. Download the specified version of the editor and driver on the release page. 
2. Execute the driver end in memory mode on the command line `./gobot-driver-win32-v0.3.x.exe --no_database --mock`
3. Start gobot_editor_win_x64_v0.3.x and fill in the driver address http://127.0.0.1:8888
4.  If using for the first time, you can find sample robots in the /demos directory and load them on the bots page.


## Feature

* Use the `behavior tree` to control the running order of the robot, and use the `script` to control the specific behavior of the node (such as making an http request
* SuProvides graphical editing and debugging capabilities
* You can `prefab` template nodes in the configuration page, and `reuse` the nodes in the editor
* It can be driven by http `api` (`post /bot.run -d '{"Name":"a robot"}'` can be easily integrated into CI
* Support a `stress test` (you can set the number of concurrency on the configuration page


## NodeScript
> Through built-in modules and scripts, we can have rich logical expression capabilities. We can also use global (single bot) meta structures to maintain various state changes of the bot.
```lua
--[[
    Each node has its own independent .lua script. When the node is executed, dostring will be called to load and run this script.
]]--

-- Users can load "modules" they want to use in the script.
-- document https://pojol.gitee.io/gobot/#/zh-cn/script/meta
local http = require("http")

-- request body
req = {
    body = {},       -- post body
    timeout = "10s", -- http timeout  
    headers = {},    -- http headers
}

-- When the robot runs to a node, the execute function will be executed.
function execute()

    -- Here, users can define the execution logic of nodes themselves (for example, sending an HTTP request)
    res, err = http.post("url", req)

end
```

## Script Module
|||||||
|-|-|-|-|-|-|
|`base64`|`http`|`protobuf`|`mongoDB`|`json`|
|`md5`|`uuid`|`random`|...|

## Try it out
Try the editor out [on website](http://123.60.17.61:7777)

## Preview
[![image.png](https://i.postimg.cc/t4jMVjp1/image.png)](https://postimg.cc/PPS4B0Lh)