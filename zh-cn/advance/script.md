# 脚本层接口

* Global
* Utils
* Http
* Json
* Report

### Global
```lua
-- merge table overwrite t2 to t1
merge(t1, t2)

--[[
    table: 0xc00005fe00 {
        [Token] => ""
    }
]]--
table.print(table)

-- meta table
meta = {
    Token = "",
    LogInfo = "",
    LogErr = "",
}
```

### Utils
```lua
local utils = require("utils")

-- 00d460f0-ec1a-4a0f-a452-1afb4b5d1686
utils.uuid()

-- random [0 ~ 100]  seed = time.now().unixnano()
utils.random(100)
```

### Http
```lua
local http = require("http")

reqtable = {
    body = {},    -- post body
    timeout = "10s",
    headers = {},
}

-- post
res, err = http.post("url", reqtable)

```

### Json
```lua
-- 
jstr = json.encode({
    Name = "joy",
    Age = 3,
})

--[[

]]--
json.decode(jstr)
```

### Protobuffer
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

### Report
```lua
```