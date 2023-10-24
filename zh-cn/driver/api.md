# API 列表
> driver 提供了一系列的 api ，可以辅助用户使用机器人                           
> 比如在 CI/CD 流程中需要对服务器进行逻辑验证，就可以通过调用 api（阻塞）来进行

---
|方法|路径|报文类型|
|-|-|-|
|POST|`/bot.run`|application/json|

### bot.run
> 运行一个阻塞执行的bot

* request
```shell
# Name - 机器人名称
$ curl -H 'Content-Type:application/json' -X POST \
    http://127.0.0.1:8888/bot.run -d '{"Name":"api_test_http.xml"}'
```
* response
```json
{
    "Code":200,
    "Msg":"",
    "Body":null
}
```

### bot.batch
> 运行一批非阻塞的bot, 通常用于压力测试或批量运行
* request
```shell
# Name - 机器人名称
# Num - 需要运行的机器人人数量（注：在集群模式下运行数量会被自动均分在各个节点中
$ curl -H 'Content-Type:application/json' -X POST \
    http://127.0.0.1:8888/bot.batch -d '{"Name":"api_test_http.xml", "Num": 1024}'
```
* response
```json
{
    "Code":200,
    "Msg":"",
    "Body":null
}
```

### bot.list
> 获取运行中的机器人列表信息
* request
```shell
$ curl -H 'Content-Type:application/json' -X POST \
    http://127.0.0.1:8888/bot.list
```
* response
```json
{
    "Code":200,
    "Msg":"",
    "Body": [
        {
            ID : "",    // 机器人的唯一id
            Name : "",  // 机器人名称
            Cur : 0,    // 当前运行了多少机器人
            Max : 1024, // 总共需要运行的数量
            Errors : 0, // 运行中遇到的错误数
        }
    ]
}
```

---

|错误码|描述|
|-|-|
|200| 成功 |
|1002| 错误的输入，请求参数不合法 |
|1007| 运行遇到异常中断 |
|1008| 找不到指定的机器人 |
