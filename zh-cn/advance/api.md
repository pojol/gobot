# 驱动层 API

* `/bot.create` - 运行一个库中的 bot
    * 请求
        ```json
            [
                {
                    "Name" string
                    "Num" int
                }
            ]
        ```
* `/bot.list` - 获取正在运行的 bot
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

* `/debug.create` - 创建一个临时调试使用的bot
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
* `/debug.step` - 将临时的bot向下执行一步
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

* `/file.txtUpload` - 上传一个新的 bot 到库中

* `/file.remove` - 移除一个库中的 bot
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

* `/file.list` - 获取库中的 bot 列表
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

* `/file.get` - 拉取库中 bot 的描述文件
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