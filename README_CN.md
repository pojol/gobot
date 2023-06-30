# gobot
Gobot是一个功能强大的有状态API测试机器人。它提供图形界面进行测试场景的搭建,可以方便的进行测试脚本编写、单步调试和压力测试,并可以在测试过程的每个阶段之间共享和存储状态。

[![Go Report Card](https://goreportcard.com/badge/github.com/pojol/gobot)](https://goreportcard.com/report/github.com/pojol/gobot)
[![](https://img.shields.io/badge/%E6%96%87%E6%A1%A3-Doc-2ca5e0?style=flat&logo=github)](https://pojol.gitee.io/gobot/#/)
[![](https://img.shields.io/badge/Trello-Todo-2ca5e0?style=flat&logo=trello)](https://trello.com/b/8eDZ6h7n/)
[![CI](https://github.com/pojol/gobot/actions/workflows/dockerimage.yml/badge.svg?branch=develop)](https://github.com/pojol/gobot/actions/workflows/dockerimage.yml)

## 快速安装
> 注：开启本地运行模式，所有的改动被记录在内存（并不会真正保存）如果需要保存请将文件下载到本地；或采用正式的部署方式
1. 在 release 页面下载指定版本的 编辑端(editor 以及 驱动端(driver
2. 在命令行以内存模式执行驱动端 `./gobot-driver-win32-v0.3.x.exe --no_database --mock`
3. 启动 gobot_editor_win_x64_v0.3.x ,并将 driver 地址填入 http://127.0.0.1:8888
4. 如果是初次使用，可以在 /sample 目录中找到示例机器人，在bots页面中载入使用

## 特性
* 使用`行为树`控制机器人的运行逻辑，使用`脚本`控制节点的具体行为（比如发起一次http请求
* 提供图形化的编辑，调试能力
* 可以`预制`模版节点，在编辑器中直接使用预制过的节点（可通过标签筛选
* 可以通过 http api `'curl post /bot.run -d '{"Name":"某个机器人"}'` 驱动一个阻塞式的机器人，通过这种方式可以方便的集成进`CI`中的测试流程
* 可以进行`压力测试`（可以在配置页设置并发数
* 提供压力测试后的API/协议`报告`查看

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
    
    -- todo
    
    -- 返回值
    --  state 状态码
    --  res 显示在 Response 面板的信息
    return state.Succ, res
end
```

## 脚本层模块
|||||||
|-|-|-|-|-|-|
|`base64`|`http`|`protobuf`|`mongoDB`|`json`|
|`md5`|`uuid`|`random`|...|

## [在线试用](http://123.60.17.61:7777)
## [文档](https://pojol.gitee.io/gobot/#/)

## [视频演示](https://www.bilibili.com/video/BV1sS4y1z7Dg/?vd_source=7c2dfd750914fd5f8a9811b19f0bf447)

## 编辑器预览
[![image.png](https://i.postimg.cc/t4jMVjp1/image.png)](https://postimg.cc/PPS4B0Lh)