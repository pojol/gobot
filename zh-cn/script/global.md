# Global

* Meta
* Json
* Merge
* table.print

### Meta
> Bot的元数据结构，所有协议的返回默认都会被merge到meta中，用于方便用户调用历史的状态数据

> 用户也可以按自己的需求将数据填写到meta中，在edit面板会实时显示meta数据

```lua
-- meta table
meta = {
}
```

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
