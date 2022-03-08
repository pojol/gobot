# API 列表

* file.upload
* file.list
* file.get
* file.remove

---

|接口|描述|
|-|-|
|`file.upload`|将bot描述文件上传|
```go

```

### bot.create
* 请求
    ```json
        [
            {
                "Name" string
                "Num" int
            }
        ]
    ```

### bot.list
* 请求
    - 空
* 结果
    ```json
        [
            {
                "ID" string     // bot 的唯一id
                "Name" string   // bot 名称
                "Cur" int       // 当前完成执行的 bot 数量
                "Max" int       // 总共需要执行的 bot 数量
                "Errors" int    // 当前遇到的 error 数量
            }
        ]
    ```

### debug.create
* 请求
    ```json
        {
            "Name" string // bot 名称
        }
    ```
* 结果
    ```json
        {
            "BotID" string // bot uuid
            "Code" int // error code
        }
    ```

### debug.step
* 请求
    ```json
        {
            "BotID" string // bot uuid
        }
    ```
* 结果
    ```json
        {
            "Prev" string       // prev node id
            "Cur" string        // cur node id
            "Blackboard" string // meta data
            "Code" int          // error code
        }
    ```

### file.txtUpload

### file.remove
* 请求
    ```json
        {
            "Name" string // bot 名称
        }
    ```
* 结果
    ```
        // 返回新的行为树列表
        [
            {
                Name string,    // 名称
                Update int64,   // 更新时间
                Status string   // 状态
            },
        ]
    ```

### file.list
* 请求 `空`
* 结果
    ```
    // 返回新的行为树列表
    [
        {
            Name string,    // 名称
            Update int64,   // 更新时间
            Status string   // 状态
        },
    ]
    ```

### file.get
* 请求
    ```json
        {
            "Name" string // bot 名称
        }
    ```
* 结果


---

|错误码|描述|
|-|-|
|200| 成功 |
 No newline at end of file
