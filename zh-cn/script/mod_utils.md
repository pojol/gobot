# Utils 模块
```lua
local utils = require("utils")
```

* UUID
* Random

### UUID
```lua
-- 00d460f0-ec1a-4a0f-a452-1afb4b5d1686
utils.uuid()
```

### Random
```lua
-- random [0 ~ 100]  seed = time.now().unixnano()
utils.random(100)
```