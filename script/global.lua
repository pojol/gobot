function print_table ( t )  
    local print_r_cache={}
    local function sub_print_r(t,indent)
        if (print_r_cache[tostring(t)]) then
            print(indent.."*"..tostring(t))
        else
            print_r_cache[tostring(t)]=true
            if (type(t)=="table") then
                for pos,val in pairs(t) do
                    if (type(val)=="table") then
                        print(indent.."["..pos.."] => "..tostring(t).." {")
                        sub_print_r(val,indent..string.rep(" ",string.len(pos)+8))
                        print(indent..string.rep(" ",string.len(pos)+6).."}")
                    elseif (type(val)=="string") then
                        print(indent.."["..pos..'] => "'..val..'"')
                    else
                        print(indent.."["..pos.."] => "..tostring(val))
                    end
                end
            else
                print(indent..tostring(t))
            end
        end
    end
    if (type(t)=="table") then
        print(tostring(t).." {")
        sub_print_r(t,"  ")
        print("}")
    else
        sub_print_r(t,"  ")
    end
    print()
end

-- initialize the meta structure
meta = {
    Token = "",
    Err = "",       -- debug log [err]
}

state = {
    Succ    = "Succ",   -- 脚本节点返回成功状态
    Exit    = "Exit",   -- 脚本节点返回退出状态（中断执行
    Error   = "Error"   -- 脚本节点返回错误状态（中断执行，且带有错误信息
}

local function _merge(t1, t2)
    for k,v in pairs(t2) do
        if type(v) == "table" then
            if type(t1[k] or false) == "table" then
                _merge(t1[k] or {}, t2[k] or {})
            else
                t1[k] = v
            end
        else
            t1[k] = v
        end
    end

    return t1
end

--[[
    merge table
    overwrite t2 to t1
]]--
function merge(t1, t2)
    _merge(t1,t2)
end

--[[
    print table like:
    table.print(meta)

    table: 0xc00005fe00 {
        [Token] => ""
        [Err] => ""
    }
]]--
table.print = print_table