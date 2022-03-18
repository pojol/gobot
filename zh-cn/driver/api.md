# API 列表
> 进阶内容， api列表描述了驱动端提供的调用接口；通常在使用工具过程中不需要了解

> 在某些情况，例如 CI/CD 流程中使用机器人对服务器进行逻辑验证；则可以通过非界面的方式，直接对驱动端进行 http 调用来完成验证。

* bot.run - 运行一个阻塞执行的bot
* bot.batch - 运行一批非阻塞的bot
* bot.list - 列举运行中的bot

* debug.create - 创建一个调试用的bot
* debug.step - 单步运行调试用的bot

* report.info - 获取 bot.batch 产生的运行报告

* file.uploadTxt - 通过txt形式将行为树文件上传到后台
* file.uploadBlob - 通过blob形式上传行为树文件数据
* file.remove - 通过名字移除一个后台的行为树文件
* file.list - 列举后台的行为树文件
* file.get - 获取后台的行为树文件数据
* file.setTags - 为后台的行为树文件添加标签

---
|方法|路径|报文类型|
|-|-|-|
|POST|`/bot.run`|application/json|

### Sample Request
```shell
$ curl -H 'Content-Type:application/json' -X POST \
    http://127.0.0.1:8888/bot.run -d '{"Name":"api_test_http.xml"}'
```
### Sample Response
```json
{
    "Code":200,
    "Msg":"",
    "Body":null
}
```

---

|错误码|描述|
|-|-|
|200| 成功 |
