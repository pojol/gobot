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

func TestWebsocket(t *testing.T) {

	L := lua.NewState()
	defer L.Close()

	WebsocketMod := NewWebsocketModule()
	protomod := ProtoModule{}
	path := "../../script"
	rand.Seed(time.Now().UnixMicro())

	// 这里要对script目录进行一次检查，将lua脚本都载入进来
	preScripts := utils.GetDirectoryFiels(path, ".lua")
	for _, v := range preScripts {
		L.DoFile(path + "/" + v)
	}

	L.PreloadModule("proto", protomod.Loader)
	L.PreloadModule("websocket", WebsocketMod.Loader)
	RegisterMessageType(L)

	ln := mock.StartWebsocketServe(L.GetGlobal("ByteOrder").String())
	go ln.Start(":6669")
	defer ln.Close()

	go func() {
		err := L.DoString(`
		local conn = require("websocket")
		local proto = require("proto")

		-- 建立连接
		local ret = conn.dail("ws", "127.0.0.1", "6669")
		print("connect websocket " .. ret)

		os.execute("sleep " .. 0.5)

		body, errmsg = proto.marshal("LoginGuestReq", json.encode({}))
		if errmsg ~= nil then
			meta.Err = "proto.marshal" .. errmsg
		end
		ret = conn.write(WSPackMsg(1001, body))
		print("create guest account " .. ret)

		for i = 0, 5, 1 do

			msgid, msgbody = WSUnpackMsg(conn.read())
			if msgid == 1001 then
				body = proto.unmarshal("LoginGuestRes", msgbody)
			elseif msgid == 1002 then
				body = proto.unmarshal("HelloRes", msgbody)
			end

			if msgid ~= 0 then
				print("recv=>", msgid, body,  err)
			end 

			local reqbody, errmsg = proto.marshal("HelloReq", json.encode({
				Message = "hello",
			}))
			conn.write(WSPackMsg(1002, reqbody))
			print(i, "send say hello " .. ret)

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
		<-time.After(time.Second * 5)
		break
	}
}
