# Global

* Json
* Merge
* table.print


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
