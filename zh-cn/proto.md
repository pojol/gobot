# 使用 protobuf 进行协议的编解码

## 参考样例 --> [Sample](https://github.com/pojol/gobot-sample) <--


### 工程目录
* proto - 协议文件
    > 这里需要注意 proto 文件中的 go_package 需要和目录一致

* script - 脚本目录（gobot会读取这个目录的所有脚本载入到每个bot的luavm中



### 协议解析 [gogoprotobuf](https://github.com/gogo/protobuf)
> 这里并未使用官方的 protobuf 解析器（因为官方实现污染了 struct，会导致无法反射； 所以这里采用的是 gogoprotobuf （比官方实现解析速度更快，并且更原生 没有 XXX_* 的定义


### 安装
> 1. 安装 [protoc](https://github.com/protocolbuffers/protobuf/releases) 推荐使用 v2.3.1+


> 2. 安装 gogo protobuf 需要 go1.9+
``` shell
go get github.com/gogo/protobuf/gogofaster_out
```


> 当完成安装后，你就可将proto文件放到自行定义的 文件夹中在 bot 中自行使用啦

```lua
local proto = require("proto")

local person = {
    name = "joy",
    id = 111,
    email = "joy@outlook.com",
    phones = { {number = "555", type = 2} }
}

local book = {
    people = { person }
}

ebody, errmsg = proto.marshal("AddressBook", json.encode(book))
if errmsg ~= "" then
    print(errmsg)
end


-- dbody, err = proto.unmarshal("AddressBook", ebody)
-- print("lua table", json.encode(dbody))

```