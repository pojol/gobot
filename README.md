# gobot
Gobot is a stateful api testing tool that supports graph editing, api calling, and binding script execution.

[![Go Report Card](https://goreportcard.com/badge/github.com/pojol/gobot)](https://goreportcard.com/report/github.com/pojol/gobot)
[![](https://img.shields.io/badge/%E6%96%87%E6%A1%A3-Doc-2ca5e0?style=flat&logo=github)](https://pojol.gitee.io/gobot/#/)
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
* Provide `prefab` template code function in config

### Script Module
|||||||
|-|-|-|-|-|-|
|`base64`|`http`|`protobuf`|`mongoDB`|`json`|
|`md5`|`uuid`|`random`|...|

# Try it out
Try the editor out [on website](http://123.60.17.61:7777)

### Preview
[![botgif2.gif](https://i.postimg.cc/SNKQG50m/botgif2.gif)](https://www.bilibili.com/video/BV1sS4y1z7Dg?share_source=copy_web)




### Report
[![image.png](https://i.postimg.cc/4d3TTrvf/image.png)](https://postimg.cc/yJ2Gmprt)