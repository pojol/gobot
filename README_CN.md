# gobot editor
Gobot是一个有状态的api测试工具，支持图形编辑、api调用、绑定脚本执行。

[![](https://img.shields.io/badge/%E6%96%87%E6%A1%A3-Doc-2ca5e0?style=flat&logo=github)](https://pojol.gitee.io/gobot/#/)
[![](https://img.shields.io/badge/Trello-Todo-2ca5e0?style=flat&logo=trello)](https://trello.com/b/8eDZ6h7n/)
[![CI](https://github.com/pojol/gobot/actions/workflows/dockerimage.yml/badge.svg?branch=develop)](https://github.com/pojol/gobot/actions/workflows/dockerimage.yml)

### 用来解决那些问题

* **状态管理**
> 当前市面上绝大部分的测试工具都是基于单条API进行测试的，但是我们同样也有一些API会依赖各种状态，例如社交中想要测试一些行为必须先是好友关系，例如游戏中想要进入战斗必须要有编队信息等等，gobot 拥有完整的生命周期我们可以在不同的节点（测试）上，保留和使用各种状态信息。

* **降低使用门槛** 
> 在定制AI逻辑的时候，首先我们想到的是使用代码（或者脚本）进行控制，这会让使用的门槛相对有些高；在 gobot 中我们可以使用界面化的工具`行为树编辑器`对机器人进行流程上的编辑，避免了手写相关控制代码。

* **提供API接口**
> 由于现代的服务器体系中引入了大量的 CI/CD 流程，在使用 gobot 的时候，我们也可以方便的在流程中插入API调用`post /bot.run -d '{"Name":"某个机器人"}'`来进行集成测试。


### 拥有那些特性
1. 使用`行为树`+`脚本`控制机器人的执行逻辑
2. 图形化的编辑，调试功能
3. 使用 tag + filter 管理和查找仓库中的机器人
4. 可以进行`压力测试`（并发驱动机器人
5. 驱动端提供`http调用接口`
6. 提供报告查看页（机器人调用的api统计信息
7. 在配置页中可以进行模版代码`预制`，为自己添加功能节点或模版节点


### [在线试用](http://123.60.17.61:7777) <--
### [文档](https://pojol.gitee.io/gobot/#/) <--


### 编辑器预览
[![botgif2.gif](https://i.postimg.cc/SNKQG50m/botgif2.gif)](https://www.bilibili.com/video/BV1sS4y1z7Dg?share_source=copy_web)

### 一个http节点的例子
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

### 脚本层支持的模块
|||||||
|-|-|-|-|-|-|
|`base64`|`http`|`protobuf`|`utils`|`mongoDB`|`json`|

### Report
[![image.png](https://i.postimg.cc/4d3TTrvf/image.png)](https://www.bilibili.com/video/BV1sS4y1z7Dg?share_source=copy_web)