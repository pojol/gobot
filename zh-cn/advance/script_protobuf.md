# Protobuf



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