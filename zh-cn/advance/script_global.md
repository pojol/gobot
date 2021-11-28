# Global

* Json

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
