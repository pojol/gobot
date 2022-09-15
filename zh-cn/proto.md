# 使用 protobuf 进行协议的编解码

参考样例 --> [Sample](https://github.com/pojol/gobot-sample) <--


1. 安装 [protoc](https://github.com/protocolbuffers/protobuf/releases) 推荐使用 v2.3.1+
2. 安装 [gogoprotobuf](https://github.com/gogo/protobuf) 需要 go1.9+
    > 这里并未使用官方的 protobuf 解析器（因为官方实现污染了 struct，会导致无法反射； 所以这里采用的是 gogoprotobuf （比官方实现解析速度更快，并且更原生 没有 XXX_* 的定义
    ``` shell
    $go get github.com/gogo/protobuf/gogofaster_out
    ```
3. 自定义 proto 文件
    ```pb
    syntax = "proto3";
    option csharp_namespace = "Google.Protobuf";
    option go_package = "book";


    // [START messages]
    message Person {
        string name = 1;
        int32 id = 2;  // Unique ID number for this person.
        string email = 3;
    
        enum PhoneType {
        MOBILE = 0;
        HOME = 1;
        WORK = 2;
        }
    
        message PhoneNumber {
        string number = 1;
        PhoneType type = 2;
        }
    
        repeated PhoneNumber phones = 4;
    }

    // Our address book file is just one of these.
    message AddressBook {
    repeated Person people = 1;
    }
    ```
4. 生成 .pb.go 文件
    ```shell
    $protoc --gogofaster_out=. book.proto
    ```
5. 将 .pb.go 移动到自己的项目中, 并通过import引入到工程中
    ```go
    import (
    _ "gobot/book"
    )
    ```
6. 在脚本中使用
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

    ebody, errmsg = proto.marshal("book.AddressBook", json.encode(book))
    if errmsg ~= "" then
        print(errmsg)
    end

    ```