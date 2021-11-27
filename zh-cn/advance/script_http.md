# HTTP 模块


### POST
```lua
local http = require("http")

reqTable = {
    body = {},       -- post body
    timeout = "10s", -- http timeout  
    headers = {},    -- http headers
}

-- post
res, err = http.post("url", reqTable)
```

### GET

### PUT


