--[[
    每个 script 节点的代码将被预编译后存放在池中，供每次执行到的时候调用
]]
--
local parm = {
    body = {}, -- request body
    timeout = "10s",
    headers = {}
}

-- REMOTE 可以存放在 global 脚本中（editor/config/global 便于统一修改
local url = REMOTE .. "/group/methon"
-- 载入预设的模块
local http = require("http")
--

--[[
    execute 每次执行到 script 或 condition 节点时会调用一次这个函数
        script 节点时返回值用作于 editor/response 面板的展示（仅调试阶段
        condition 节点时返回值用于 判定节点执行结果（true or false
]]
function execute()
    res, errmsg = http.post(url, parm)
    if errmsg ~= nil then
        bot.Meta.Err = errmsg
        return
    end

    if res["status_code"] ~= 200 then
        bot.Meta.Err = "post " .. url .. " http status code err " .. res["status_code"]
        return
    end

    body = json.decode(res["body"])
    merge(bot, body.Body)  -- 将 res 数据合并到 meta 结构中（覆盖

    return state.Succ, body -- 将http response传递给 editor 中的 response 栏
end
