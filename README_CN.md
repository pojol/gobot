# gobot
Gobot是一个有状态的api测试工具，支持图形编辑、api调用、绑定脚本执行。

[![Go Report Card](https://goreportcard.com/badge/github.com/pojol/gobot)](https://goreportcard.com/report/github.com/pojol/gobot)
[![](https://img.shields.io/badge/%E6%96%87%E6%A1%A3-Doc-2ca5e0?style=flat&logo=github)](https://pojol.gitee.io/gobot/#/)
[![](https://img.shields.io/badge/Trello-Todo-2ca5e0?style=flat&logo=trello)](https://trello.com/b/8eDZ6h7n/)
[![CI](https://github.com/pojol/gobot/actions/workflows/dockerimage.yml/badge.svg?branch=develop)](https://github.com/pojol/gobot/actions/workflows/dockerimage.yml)


### 特性
* 使用`行为树`控制机器人的运行逻辑，使用`脚本`控制节点的具体行为（比如发起一次http请求
* 提供图形化的编辑，调试能力
* 可以在配置页中`预制`模版节点，在编辑器中进行节点的`复用`
* 可以通过 http `api` 进行驱动（`post /bot.run -d '{"Name":"某个机器人"}'` 可以方便的集成进CI
* 进行`压力测试`（可以在配置页设置并发数
* 提供报告查看页（一些简略的统计信息


### 脚本层模块
|||||||
|-|-|-|-|-|-|
|`base64`|`http`|`protobuf`|`mongoDB`|`json`|
|`md5`|`uuid`|`random`|...|

### [在线试用](http://123.60.17.61:7777) <--
### [文档](https://pojol.gitee.io/gobot/#/) <--


### 编辑器预览
[![botgif2.gif](https://pojol.oss-cn-shanghai.aliyuncs.com/gobot/botgif2.gif
)](https://www.bilibili.com/video/BV1sS4y1z7Dg?share_source=copy_web)



### Report
[![image.png](https://pojol.oss-cn-shanghai.aliyuncs.com/gobot/report.png
)](https://www.bilibili.com/video/BV1sS4y1z7Dg?share_source=copy_web)
[![image.png](https://pojol.oss-cn-shanghai.aliyuncs.com/gobot/gobot_qps.png)]