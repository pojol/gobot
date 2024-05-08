# HTTP


```lua
local http = require("http")
```

### POST
```lua

-- Request body
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
    res["body_size"] : 
    res["headers"] : http headers
    res["cookies"] : http cookies
    res["status_code"] : http status code
    res["url"] : request url

    err : error message
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

