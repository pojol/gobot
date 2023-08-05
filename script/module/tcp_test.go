package script

import (
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	"github.com/pojol/gobot/utils"
	lua "github.com/yuin/gopher-lua"
)

// 处理函数,在一个新的goroutine中处理每个连接的请求
func tcphandle(conn net.Conn) {
	//go keepalive(conn)
	defer func() {
		if err := recover(); err != nil {
			log.Printf("recover from err: %v", err)
		}
	}()

	buf := make([]byte, 128)
	// 使用 for 循环读取请求,如果遇到错误则break跳出
	for {

		n, err := conn.Read(buf)
		if n == 0 || err != nil {
			if err != nil {
				fmt.Println("read ", n, err)
			}
			continue
		}

		fmt.Printf("recv: %s\n", buf[:n])
		buf = make([]byte, 128) // 重置 buf

		// 处理请求......
		_, err = conn.Write([]byte("world !"))
		if err != nil {
			log.Printf("write to client error: %v", err)
			break
		}
	}

	// 请求循环结束,关闭连接
	fmt.Println("server conn closed")
	conn.Close()
}

func starMockServer() net.Listener {
	ln, err := net.Listen("tcp", ":20008")
	if err != nil {
		panic(err)
	}

	fmt.Println("Server listening on port 20008")
	go func() {
		for {
			// 接收新连接
			conn, err := ln.Accept()
			if err != nil {
				//fmt.Println("accept err", err)
				continue
			}

			fmt.Println("new client conn =>", conn.RemoteAddr())

			// 为每个连接启动一个goroutine进行处理
			go tcphandle(conn)
		}
	}()

	return ln
}

func TestTcpConnect(t *testing.T) {

	L := lua.NewState()
	defer L.Close()

	ln := starMockServer()
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
