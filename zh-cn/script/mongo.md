# Mongo 模块


### conn
> 连接到数据库
```lua
-- "bot" 需要连接到的数据库
-- "mongodb://" 连接字符串
    -- https://www.docs4dev.com/docs/zh/mongodb/v3.6/reference/reference-connection-string.html#connections-connection-examples

ret = mgo.conn("bot", "mongodb://127.0.0.1:27017")
assert(ret == "succ", "mgo connect err " .. ret)
```

### disconn
> 断开数据库连接

```lua
mgo.disconn()
```

### insert_one
> 插入一个数据

```lua
ret = mgo.insert_one(testdb, {_id = "001", msg = "aa"})
assert(ret == "succ", "mgo insert_one err " .. ret)

--[[
    {
        _id : "001",
        msg : "aa"
    }
]]--
```

### insert_many
> 插入一组数据

```lua
doc = {
    {a = 1, b = "a", c = false},
    {a = 2, b = "b", c = true},
    {a = { b = "c" }},
    {a = {"a", "b", "c", "d"}},
    {a = {1, 2, 3, 4}}
}

ret = mgo.insert_many(testdb, doc)
assert(ret == "succ", "mgo insert many err " .. ret)
```

### find
> 查找数据

```lua
val, ret = mgo.find(testdb, {})
assert(ret == "succ", "mgo find err " .. ret)
```

### find_one
> 查找一份指定的数据

```lua
ret = mgo.insert_one(testdb, {a = 1, b = "aa"})
assert(ret == "succ", "mgo insert one err " .. ret)

val, ret = mgo.find_one(testdb, {a = 1})
assert(ret == "succ", "mgo find err " .. ret)
```

### update_one
> 更新一份指定的数据

```lua
ret = mgo.insert_one(testdb, {a = 1, b = "aa"})
assert(ret == "succ", "mgo insert one err " .. ret)

settable = {}
settable["$set"] = {b = "bb"}
mgo.update_one(testdb, {a = 1}, settable)
```

### update_many
> 更新一组指定的数据

### delete_one
> 删除一份指定的数据

```lua
ret = mgo.insert_one(testdb, {_id = test1id, msg = "aa"})
assert(ret == "succ", "mgo insert_one err " .. ret)

assert(lt._id == test1id, "find err " .. test1id)
mgo.delete_one(testdb, {_id = test1id})
```

### delete_many
> 删除一组指定的数据