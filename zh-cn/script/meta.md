### Meta
> Bot的全局数据结构（行为树中的 blackboard


```lua
--[[
    meta table
        这是一个空的 table 用于存储 bot 的全局状态
        访问 - 在任何节点中使用 meta. 既可获取 meta 中的数据
        赋值 - 在默认的 http 模板中，我们会将 http res 返回的 body 序列化成 table 后 merge 到 meta 结构上
            - 在用户需要精确控制的情况下，我们也可以自定义赋值逻辑
]]--
meta = {
}
```

### 模板自动赋值
### 自定义赋值