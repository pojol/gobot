package script

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/pojol/gobot/mock"
	"github.com/pojol/gobot/utils"
	lua "github.com/yuin/gopher-lua"
)

func TestCustomMsgPack(t *testing.T) {

	L := lua.NewState()
	defer L.Close()

	ln := mock.StarTCPServer(":20008")
	defer ln.Close()

	tcpMod := TCPModule{}
	protomod := ProtoModule{}
	path := "../../script"
	rand.Seed(time.Now().UnixMicro())

	// 这里要对script目录进行一次检查，将lua脚本都载入进来
	preScripts := utils.GetDirectoryFiels(path, ".lua")
	for _, v := range preScripts {
		L.DoFile(path + "/" + v)
	}
	L.PreloadModule("proto", protomod.Loader)
	L.PreloadModule("tcpconn", tcpMod.Loader)
	registerMessageType(L)

	go func() {
		err := L.DoString(`
		local conn = require("tcpconn")
		local proto = require("proto")

		local ret = conn.dail("127.0.0.1", "20008")
		print("conn dail " .. ret)

		os.execute("sleep " .. 0.5)

		body, errmsg = proto.marshal("LoginGuestReq", json.encode({}))
		if errmsg ~= nil then
			meta.Err = "proto.marshal" .. errmsg
		end
		ret = conn.write(TCPPackMsg(1001, body))
		print("write msg 1001 " .. ret)

		for i = 0, 5, 1 do
			--[[
				| 2 byte   , 1 byte,    2 byte      , 2byte		  |                        |
				|包长度 len, 协议格式 ty, 预留2自定义字节, 协议号 msgid |                        |
				|                  消息头                          |         消息体          |
			]]--

			msgid, msgbody = TCPUnpackMsg(conn.read(2))

			if msgid == 1001 then
				body = proto.unmarshal("LoginGuestRes", msgbody)
			elseif msgid == 1002 then
				body = proto.unmarshal("HelloRes", msgbody)
			end

			local reqbody, errmsg = proto.marshal("HelloReq", json.encode({
				Message = "hello",
			}))
			ret = conn.write(TCPPackMsg(1002, reqbody))
			print("write msg 1002 " .. ret)

			print("read==>", msgid, body, err)
			os.execute("sleep " .. 0.5)
		end

		print("client connect closed !")
		conn.close()
		`)

		if err != nil {
			fmt.Println("dostring", err.Error())
		}
	}()

	for {
		<-time.After(time.Second * 3)
		break
	}
}
