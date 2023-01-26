# gobot
Gobot is a stateful api testing tool that supports graph editing, api calling, and binding script execution.

[![Go Report Card](https://goreportcard.com/badge/github.com/pojol/gobot)](https://goreportcard.com/report/github.com/pojol/gobot)
[![](https://img.shields.io/badge/%E6%96%87%E6%A1%A3-Doc-2ca5e0?style=flat&logo=github)](https://pojol.gitee.io/gobot/#/)
[![](https://img.shields.io/badge/Trello-Todo-2ca5e0?style=flat&logo=trello)](https://trello.com/b/8eDZ6h7n/)
[![CI](https://github.com/pojol/gobot/actions/workflows/dockerimage.yml/badge.svg?branch=develop)](https://github.com/pojol/gobot/actions/workflows/dockerimage.yml)

[中文](https://github.com/pojol/gobot/blob/master/README_CN.md)


## Feature

* Use the `behavior tree` to control the running order of the robot, and use the `script` to control the specific behavior of the node (such as making an http request
* SuProvides graphical editing and debugging capabilities
* You can `prefab` template nodes in the configuration page, and `reuse` the nodes in the editor
* It can be driven by http `api` (`post /bot.run -d '{"Name":"a robot"}'` can be easily integrated into CI
* Support a `stress test` (you can set the number of concurrency on the configuration page


## Script Module
|||||||
|-|-|-|-|-|-|
|`base64`|`http`|`protobuf`|`mongoDB`|`json`|
|`md5`|`uuid`|`random`|...|

## Try it out
Try the editor out [on website](http://123.60.17.61:7777)

## Preview
[![image.png](https://i.postimg.cc/x1cFp5vF/image.png)](https://postimg.cc/dhcBL8S8)