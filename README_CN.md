# gobot
Gobot是一个有状态的api测试工具，支持图形编辑、api调用、绑定脚本执行。

[![Go Report Card](https://goreportcard.com/badge/github.com/pojol/gobot)](https://goreportcard.com/report/github.com/pojol/gobot)
[![](https://img.shields.io/badge/%E6%96%87%E6%A1%A3-Doc-2ca5e0?style=flat&logo=github)](https://pojol.gitee.io/gobot/#/)
[![](https://img.shields.io/badge/Trello-Todo-2ca5e0?style=flat&logo=trello)](https://trello.com/b/8eDZ6h7n/)
[![CI](https://github.com/pojol/gobot/actions/workflows/dockerimage.yml/badge.svg?branch=develop)](https://github.com/pojol/gobot/actions/workflows/dockerimage.yml)


## 特性
* 使用`行为树`控制机器人的运行顺序，使用`脚本`控制节点的具体行为（比如发起一次http请求
* 提供图形化的编辑，调试能力
* 可以在配置页中`预制`模版节点，在编辑器中进行节点的`复用`
* 可以通过 http `api` 进行驱动（`post /bot.run -d '{"Name":"某个机器人"}'` 可以方便的集成进CI
* 进行`压力测试`（可以在配置页设置并发数
* 提供报告查看页（一些简略的统计信息

## 节点脚本
> 通过内置的模块+脚本可以使我们拥有丰富的逻辑表达能力，也可以使用全局的（单个bot）meta 结构来维护 bot 的各种状态变更
```lua
--[[
    每个节点都拥有一个独立属于自己的 .lua 脚本，当节点被执行到的时候会调用 dostring 加载并运行这个脚本
]]--

-- 用户可以在脚本中加载自己想要使用的 “模块”
-- document https://pojol.gitee.io/gobot/#/zh-cn/script/meta
local http = require("http")

-- 定义一些逻辑所需的结构
req = {
    body = {},       -- post body
    timeout = "10s", -- http timeout  
    headers = {},    -- http headers
}

-- 当脚本成功加载后，会调用这个 execute 函数
function execute()

    -- 用户可以在这里自行定义节点的执行逻辑（例如发送一次http请求
    res, err = http.post("url", req)

end
```

## 脚本层模块
|||||||
|-|-|-|-|-|-|
|`base64`|`http`|`protobuf`|`mongoDB`|`json`|
|`md5`|`uuid`|`random`|...|

## [在线试用](http://123.60.17.61:7777) <--
## [文档](https://pojol.gitee.io/gobot/#/) <--


## 编辑器预览
[![botgif2.gif](https://i.postimg.cc/NGmddV78/botgif3.gif
)](https://www.bilibili.com/video/BV1sS4y1z7Dg?share_source=copy_web)

[![image.png](https://pojol.oss-cn-shanghai.aliyuncs.com/gobot/gobot_qps.png)