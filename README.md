# 简介

> Gobot 是一个支持`可视化`的机器人编辑平台，它依据编辑出来的`行为树`文件进行驱动，同时支持挂接 `lua` 脚本节点 ，用于执行具体的逻辑；



## 特性
* 使用`行为树`控制机器人的运行顺序，使用`脚本`控制节点的具体行为（比如发起一次http请求
* 提供图形化的编辑，调试能力
* 可以在配置页中`预制`模版节点，在编辑器中进行节点的`复用`
* 可以通过 http `api` 进行驱动（`post /bot.run -d '{"Name":"某个机器人"}'` 可以方便的集成进CI
* 进行`压力测试`（可以在配置页设置并发数
* 提供报告查看页（一些简略的统计信息



### 界面预览
![img](/res/preview.png)


### 客户端支持
||Web|Windows|Mac|Android|IOS|
|-|-|-|-|-|-|
|i32|✅|✅|❌|❌|❌|
|x64|✅|✅|✅|❌|❌|
|arm|✅|❌|✅|❌|❌|