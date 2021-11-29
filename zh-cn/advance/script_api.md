# 捕获 API 结构

> 用户可以通过脚本简单的获取到协议的上下文，并将结构记录下来，用于展示给不同的用户阅览；

```lua
api_info = {
    -- lst.url 请求地址
    -- lst.req 请求结构
    -- lst.res 回复结构
    -- lst.desc 协议描述
    lst = {}
}

function api_info.record_req(url, req, desc)
    if type(req) ~= "table" then
        error("expected argument of type table got " .. type(req))
    end

    local exist = false
    for k,v in pairs(api_info.lst) do
        if v.url == url then
            exist = true
            merge(api.lst[k].req, req)
        end
    end

    if not exist then
        table.insert(api_info.lst, {
            url = url,
            req = req,
            res = {},
            desc = desc,
        })
    end

end

function api_info.record_res(url, res)
    if type(res) ~= "table" then
        error("expected argument of type table got" .. type(res))
    end

    for k,v in pairs(api_info.lst) do
        if v.url == url then
            merge(api_info.lst[k].res, res)
        end
    end
end


function api_info.report()
    local http = require("http")

    local api_list = json.encode(api_info.lst)

end
```