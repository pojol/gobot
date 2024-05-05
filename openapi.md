# OpenAPI


> The driver provides a series of APIs to assist users in using the robot. For example, in the CI/CD process, if logical validation of servers is required, it can be achieved by calling APIs (blocking).

---
|Method|Path|Content-Type|
|-|-|-|
|POST|`/bot.run`|application/json|

### bot.run
> Run a bot that blocks execution

* request
```shell
# Name - bot name
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
> Run a batch of non-blocking bots, typically for stress testing
* request
```shell
# Name - bot name
# Num - The number of robots that need to be run (Note: In cluster mode, the number of robots will be automatically divided evenly among each node.
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
> Get running robot list information
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
            ID : "",    // The unique id of the robot
            Name : "",  // name
            Cur : 0,    // How many bots are currently running
            Max : 1024, // The total number of runs required
            Errors : 0, // Number of errors encountered during operation
        }
    ]
}
```

---

|errcode|desc|
|-|-|
|200| Succ |
|1002| Incorrect input, invalid request parameters |
|1007| The operation encountered an abnormal interruption |
|1008| The specified robot cannot be found |
