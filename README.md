# 简介

> Gobot 是一个支持`可视化`的机器人编辑平台，它依据编辑出来的`行为树`文件进行驱动，同时支持挂接 `lua` 脚本节点 ，用于执行具体的逻辑；

---
# v0.1.8

### 用来解决那些问题

* **状态管理**
> 当前市面上绝大部分的测试工具都是基于单条API进行测试的，但是我们同样也有一些API会依赖各种状态，例如社交中想要测试一些行为必须先是好友关系，例如游戏中想要进入战斗必须要有编队信息等等，gobot 拥有完整的生命周期我们可以在不同的节点（测试）上，保留和使用各种状态信息。

* **降低使用门槛** 
> 在定制AI逻辑的时候，首先我们想到的是使用代码（或者脚本）进行控制，这会让使用的门槛相对有些高；在 gobot 中我们可以使用界面化的工具`行为树编辑器`对机器人进行流程上的编辑，避免了手写相关控制代码。

* **提供API接口**
> 由于现代的服务器体系中引入了大量的 CI/CD 流程，在使用 gobot 的时候，我们也可以方便的在流程中插入API调用`post /bot.run -d '{"Name":"某个机器人"}'`来进行集成测试。



### 拥有那些特性
1. 使用行为树 + 脚本控制机器人的执行逻辑
2. 图形化的编辑，调试功能
3. 使用 tag + filter 管理和查找仓库中的机器人
4. 可以进行压力测试（并发驱动机器人
5. 驱动端提供api调用
6. 提供报告查看页（机器人调用的api统计信息



### 客户端支持
||Web|Windows|Mac|Android|IOS|
|-|-|-|-|-|-|
|i32|✅|✅|❌|❌|❌|
|x64|✅|✅|✅|❌|❌|
|arm|✅|❌|✅|❌|❌|



### 界面预览
![img](/res/preview.png)



### 性能预览

> 单位机器人并发数量下的 QPS

![img](/res/gobot_qps.png)

> 报告概览

![img](/res/report.png)