# gobot
Gobot是一个功能强大的有状态API测试机器人。它提供图形界面进行测试场景的搭建,可以方便的进行测试脚本编写、单步调试和压力测试,并可以在测试过程的每个阶段之间共享和存储状态。

[![Go Report Card](https://goreportcard.com/badge/github.com/pojol/gobot)](https://goreportcard.com/report/github.com/pojol/gobot)
[![](https://img.shields.io/badge/%E6%96%87%E6%A1%A3-Doc-2ca5e0?style=flat&logo=github)](https://pojol.gitee.io/gobot/#/)
[![](https://img.shields.io/badge/Trello-Todo-2ca5e0?style=flat&logo=trello)](https://trello.com/b/8eDZ6h7n/)
[![CI](https://github.com/pojol/gobot/actions/workflows/dockerimage.yml/badge.svg?branch=develop)](https://github.com/pojol/gobot/actions/workflows/dockerimage.yml)

## 快速安装
> 注：开启本地运行模式，所有的改动被记录在内存（并不会真正保存）如果需要保存请将文件下载到本地；或采用正式的部署方式
1. 进入最新的 [release页面](https://github.com/pojol/gobot/releases/tag/v0.3.8) 下载可执行程序
2. 执行 gobot_driver_win_x64_v0.3.8 目录中的 run.bat 文件， 运行服务器
    * 执行 gobot_editor_win_x64_v0.3.8 目录中的 gobot.ext， 运行编辑器程序
3. 在弹出的地址输入窗口 或 config 页的地址栏中填入 http://127.0.0.1:8888 本地服务器地址
4. 切换到编辑器的 bots 面板，将 http_sample.txt 和 tcp_sample.txt 两个用例拖入
5. 选择一个用例，点击 load 将机器人加载到编辑界面
    * 点击下方的 debug （爬虫）按钮进行调试（创建一个新的调试机器人
    * 点击旁边的 运行 按钮，单步执行（运行行为树节点
    * 点击编辑器中的任意一个节点 可以查看这个节点的设置
    * Meta 面板 可以查看机器人的所有数据
    * Response 显示的是每个节点中的返回值
    * RuntimeErr 显示的是执行节点可能遇到的错误信息（会自动切换过去

## 特性
* 使用`行为树`控制机器人的运行逻辑，使用`脚本`控制节点的具体行为（比如发起一次http请求
* 提供图形化的编辑，调试能力
* 可以`预制`模版节点，在编辑器中直接使用预制过的节点（可通过标签筛选
* 可以通过 http api `'curl post /bot.run -d '{"Name":"某个机器人"}'` 驱动一个阻塞式的机器人，通过这种方式可以方便的集成进`CI`中的测试流程
* 支持多种协议格式（HTTP, TCP ...
* 可以进行`压力测试`（可以在配置页设置不同的并发策略
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
| Module | interface |Description |
|-------------|-------------|-------------|
| base64 | `encode` `decode` |Provides base64 encoding/decoding functionality.|
| http | `post` `get` `put` | Support HTTP connection. |
| tcp | `dail` `close` `write` `read` | Support TCP connection. |
| protobuf | `marshal` `unmarshal` | Provides Protobuf operations. |
| mongoDB | `insert` `find` `update` `delete` ... | Provides MongoDB operations. |
| json | `encode` `decode` | Offers JSON functionalities. |
| md5 | `sum` | Calculates MD5 hashes. |
| utils | `uuid` `random` | Generates random values, UUIDs. |
| ... | More modules available. |

## [在线试用](http://178.128.113.58:31293)
> 服务器部署在新加坡（ping有些高

## [文档](https://pojol.gitee.io/gobot/#/)

## [视频演示](https://www.bilibili.com/video/BV1sS4y1z7Dg/?vd_source=7c2dfd750914fd5f8a9811b19f0bf447)

## 编辑器预览
[![image.png](https://i.postimg.cc/t4jMVjp1/image.png)](https://postimg.cc/PPS4B0Lh)