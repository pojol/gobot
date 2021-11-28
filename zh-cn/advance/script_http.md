# HTTP 模块


### POST
```lua
local http = require("http")

-- 请求结构
reqTable = {
    body = {},       -- post body
    timeout = "10s", -- http timeout  
    headers = {},    -- http headers
}

-- post
res, err = http.post("url", reqTable)
--[[
    res : userdata
    res["body"] : http response body
    res["body_size"] : 消息体长度
    res["headers"] : http headers
    res["cookies"] : http cookies
    res["status_code"] : http status code
    res["url"] : 请求地址

    err : 错误信息，如果没有则为 nil
]]--

```

### GET
```lua
res, err = http.get("url", {})
```

### PUT
```lua
res, err = http.get("put", {})
```

