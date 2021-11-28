# Protobuf

* 怎样引入自己的protobuf文件
* 使用protobuf编解码

### 怎样引入自己的protobuf文件
* 有golang环境
    ```shell
    go install github.com/gogo/protobuf/protoc-gen-gogofaster
    ```
* 没有golang环境

### 使用protobuf编解码
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

-- marshal
byt, err = proto.marshal("AddressBook", json.encode(book))

-- ununmarshal
table, err = proto.ununmarshal("AddressBook", byt)
```