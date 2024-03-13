
-- 字节序
-- 大端：BigEndian
-- 小端：LittleEndian
ByteOrder = "BigEndian"


-- 从报文中解析出消息ID和消息体
-- 示例，用户可以参照实际报文格式进行解析
function WSUnpackMsg(buf, errmsg)

    if errmsg ~= "nil" then
        return 0, ""
    end

    local msg = message.new(buf, ByteOrder, 0)

    local msgId = msg:readi2()
    local msgbody = msg:readBytes(2, -1)

    return msgId, msgbody
    
end

function WSPackMsg(msgid, msgbody)

    local msg = message.new("", ByteOrder, 2+#msgbody)
    msg:writei2(msgid)
    msg:writeBytes(msgbody)

    return msg:pack()

end

------------------------------------------------------------------------

function TCPUnpackMsg(msglen, buf, errmsg)
    if errmsg ~= "nil" then
        return 0, ""
    end

    local msg = message.new(buf, ByteOrder, 0)

    local msgTy = msg:readi1()
    local msgCustom = msg:readi2()
    local msgId = msg:readi2()
    local msgbody = msg:readBytes(msglen-(2+1+2+2), -1)

    return msgId, msgbody

end

function TCPPackMsg(msgid, msgbody)
    local msglen = #msgbody+2+1+2+2

    local msg = message.new("", ByteOrder, msglen)
    msg:writei2(msglen)
    msg:writei1(1)
    msg:writei2(0)
    msg:writei2(msgid)
    msg:writeBytes(msgbody)

    return msg:pack()

end