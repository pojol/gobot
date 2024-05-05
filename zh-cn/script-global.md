# Global


* Json
* Merge
* table.print

> Bot的全局数据结构（存放当前bot的各种状态信息

```lua
--[[
    meta table
        这是一个空的 table 用于存储 bot 的全局状态
        访问 - 在任何节点中使用 meta. 既可获取 meta 中的数据
        赋值 - 在默认的 http 模板中，我们会将 http res 返回的 body 序列化成 table 后 merge 到 meta 结构上
            - 在用户需要精确控制的情况下，我们也可以自定义赋值逻辑
]]--
bot = {
    Meta = {            -- meta data (Gobot system occupancy structure
        Err = "",       -- debug log [err]
        ID = "",        -- bot index
        Batch = "",     -- batch index(Only assign value when running in batches
        Name = "",      -- bot name
    },
}

-- 脚本节点返回状态
state = {
    Succ    = "Succ",   -- 脚本节点返回成功状态
    Error   = "Error",  -- 脚本节点返回错误状态（正常执行，但携带错误
    Break   = "Break",  -- 返回中断状态（中断执行，且携带错误
    Exit    = "Exit",   -- =返回退出状态（中断执行，正常退出
}
```

### 脚本样例
```lua
--[[
    每个 script 节点的代码将被预编译后存放在池中，供每次执行到的时候调用
]]--
local parm = {
    body = {},    -- request body
    timeout = "10s",
    headers = {},
}

-- REMOTE 可以存放在 global 脚本中（editor/config/global 便于统一修改
local url = REMOTE .. "/group/methon"
-- 载入预设的模块
local http = require("http")

--[[
    execute 每次执行到 script 或 condition 节点时会调用一次这个函数
        script 节点时返回值用作于 editor/response 面板的展示（仅调试阶段
        condition 节点时返回值用于 判定节点执行结果（true or false
]]--
function execute()
    res, errmsg = http.post(url, parm)
  	if errmsg ~= nil then
		meta.Err = errmsg
    	return
  	end
  	
  	if res["status_code"] ~= 200 then
		meta.Err = "post " .. url .. " http status code err " .. res["status_code"]
  		return
  	end
  
  	body = json.decode(res["body"])
  	merge(meta, body.Body)  -- 将 res 数据合并到 meta 结构中（覆盖

    -- 将http response传递给 editor 中的 response 栏
    return state.Succ, body.Body 
end
```

### 模板自动赋值
### 自定义赋值

### Merge
```
-- merge table overwrite t2 to t1
merge(t1, t2)
```

### table.print
```lua
--[[
    table: 0xc00005fe00 {
        [Token] => ""
    }
]]--
table.print(table)
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