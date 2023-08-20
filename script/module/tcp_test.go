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

func TestTcpConnect(t *testing.T) {

	L := lua.NewState()
	defer L.Close()

	ln := mock.StarTCPServer("20008")
	defer ln.Close()

	tcpMod := TCPModule{}
	path := "../../script"

	// 这里要对script目录进行一次检查，将lua脚本都载入进来
	preScripts := utils.GetDirectoryFiels(path, ".lua")
	for _, v := range preScripts {
		L.DoFile(path + "/" + v)
	}
	L.PreloadModule("tcpconn", tcpMod.Loader)

	go func() {
		err := L.DoString(`
		local conn = require("tcpconn")
		
		local ret = conn.dail("127.0.0.1", "20008")
		print("conn dail " .. ret)

		ret = conn.write("hello")
		print("write " .. ret)

		for i = 0, 5, 1 do
			print("read==>")
			local state, msg = conn.read()
			if state == "succ" then
				print("read succ << " .. msg)

				ret = conn.write("hello")
				print("write " .. ret)
			end

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

func TestTcpMsgPack(t *testing.T) {

	L := lua.NewState()
	defer L.Close()

	ln := mock.StarTCPServer("20008")
	defer ln.Close()

	tcpMod := TCPModule{}
	path := "../../script"
	rand.Seed(time.Now().UnixMicro())

	// 这里要对script目录进行一次检查，将lua脚本都载入进来
	preScripts := utils.GetDirectoryFiels(path, ".lua")
	for _, v := range preScripts {
		L.DoFile(path + "/" + v)
	}
	L.PreloadModule("tcpconn", tcpMod.Loader)

	go func() {
		err := L.DoString(`
		local conn = require("tcpconn")

		local ret = conn.dail("127.0.0.1", "20008")
		print("conn dail " .. ret)

		conn.write_msg(7+5, 1, 0, 1000, "hello")

		for i = 0, 5, 1 do
			--[[
				| 2 byte   , 1 byte,    2 byte      , 2byte		  |                        |
				|包长度 len, 协议格式 ty, 预留2自定义字节, 协议号 msgid |                        |
				|                  消息头                          |         消息体          |
			]]--

			ty, _, msgid, msgbody, err = conn.read_msg(2,1,2,2)
			print("read==>", ty, msgid, msgbody, err)

			if msgid == 1001 then
				conn.write_msg(7+3, 1, 0, 1001, "joy")
			elseif msgid == 1002 then
				conn.write_msg(7+3, 1, 0, 1002, "ppp")
			end

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
