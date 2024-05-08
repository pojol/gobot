# MongoDB

```lua
local mgo = require("mgo")
```

### conn
```lua
local mgo = require("mgo")

-- "bot" 
-- "mongodb://" connect url
    -- https://www.docs4dev.com/docs/zh/mongodb/v3.6/reference/reference-connection-string.html#connections-connection-examples

ret = mgo.conn("bot", "mongodb://127.0.0.1:27017")
assert(ret == "succ", "mgo connect err " .. ret)
```

### disconn

```lua
mgo.disconn()
```

### insert_one

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

--[[
[
    {
        "_id":{
            "$oid":"6288e5aa2a534719e4103d72"
        },
        "b":"a",
        "c":false,
        "a":1
    },
    {
        "_id":{
            "$oid":"6288e5aa2a534719e4103d73"
        },
        "a":2,
        "b":"b",
        "c":true
    },
    {
        "_id":{
            "$oid":"6288e5aa2a534719e4103d74"
        },
        "a":{
            "b":"c"
        }
    },
    {
        "_id":{
            "$oid":"6288e5aa2a534719e4103d75"
        },
        "a":[ "a", "b", "c", "d" ]
    },
    {
        "_id":{
            "$oid":"6288e5aa2a534719e4103d76"
        },
        "a":[ 1, 2, 3, 4 ]
    }
]
]]--
```

### find

```lua
val, ret = mgo.find(testdb, {})
assert(ret == "succ", "mgo find err " .. ret)
```

### find_one

```lua
ret = mgo.insert_one(testdb, {a = 1, b = "aa"})
assert(ret == "succ", "mgo insert one err " .. ret)

val, ret = mgo.find_one(testdb, {a = 1})
assert(ret == "succ", "mgo find err " .. ret)
```

### update_one

```lua
ret = mgo.insert_one(testdb, {a = 1, b = "aa"})
assert(ret == "succ", "mgo insert one err " .. ret)

settable = {}
settable["$set"] = {b = "bb"}
mgo.update_one(testdb, {a = 1}, settable)
```

### update_many

### delete_one

```lua
ret = mgo.insert_one(testdb, {_id = test1id, msg = "aa"})
assert(ret == "succ", "mgo insert_one err " .. ret)

assert(lt._id == test1id, "find err " .. test1id)
mgo.delete_one(testdb, {_id = test1id})
```

### delete_many
