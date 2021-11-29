# Protobuf

* 怎样引入自己的protobuf文件
* 使用protobuf编解码

### 怎样引入自己的protobuf文件
1. 引入 protobuf 需要构建自己的 gobot
```shell
$ git clone https://github.com/pojol/gobot-sample-protobuf
```
2. 安装生成器
```shell
$ go install github.com/gogo/protobuf/protoc-gen-gogofaster
```
3. 创建一个目录(和protobuf包同名
```shell
$ mkdir book

# 将proto定义文件拷贝到目录
$ cp proto.* book/ 

# 通过proto文件生成.go协议解析文件
$ protoc --gogofaster_out=. *.proto 

# 在 main.go 中引用这个包
import (
    _ "gobot/book"
)

# 通过 docker 进行镜像的构建
$ docker build
```

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