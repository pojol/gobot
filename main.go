package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"runtime"
	"strconv"

	_ "net/http/pprof"

	"github.com/pojol/gobot/driver/factory"
	"github.com/pojol/gobot/driver/mock"
	"github.com/pojol/gobot/driver/openapi"
	"github.com/pojol/gobot/driver/utils"
	lua "github.com/yuin/gopher-lua"
)

const (
	// Version of gobot driver
	Version = "v0.4.4"

	banner = `
              __              __      
             /\ \            /\ \__   
   __     ___\ \ \____    ___\ \ ,_\  
 /'_ '\  / __'\ \ '__'\  / __'\ \ \  
/\ \L\ \/\ \L\ \ \ \L\ \/\ \L\ \ \ \_ 
\ \____ \ \____/\ \_,__/\ \____/\ \__\
 \/___L\ \/___/  \/___/  \/___/  \/__/
   /\____/                            
   \_/__/             %s                

`
)

func main() {

	defer func() {
		if err := recover(); err != nil {
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			fmt.Println("panic:", string(buf[:n]))
		}
	}()

	f := utils.InitFlag()
	flag.Parse()
	if utils.ShowUseage() {
		return
	}

	botFactory, err := factory.Create(
		factory.WithDatabase(f.DBType),
		factory.WithClusterMode(f.Cluster),
	)
	if err != nil {
		panic(err)
	}
	defer botFactory.Close()

	L := lua.NewState()
	defer L.Close()
	L.DoFile(f.ScriptPath + "/" + "message.lua")
	byteOrder := L.GetGlobal("ByteOrder").String()

	if f.OpenHttpMock != 0 {
		ms := mock.NewHttpServer()
		go ms.Start(":" + strconv.Itoa(f.OpenHttpMock))
		defer ms.Close()
	}

	if f.OpenTcpMock != 0 {
		tcpls := mock.StarTCPServer(byteOrder, ":"+strconv.Itoa(f.OpenTcpMock))
		defer tcpls.Close()
	}

	if f.OpenWSMock != 0 {
		ws := mock.StartWebsocketServe(byteOrder, ":"+strconv.Itoa(f.OpenWSMock))
		defer ws.Close()
	}

	go func() {
		http.ListenAndServe(":6060", nil)
	}()

	// 查看有没有未完成的队列
	factory.Global.CheckTaskHistory()

	fmt.Printf(banner, Version)

	openApiPort := 8888
	if f.OpenAPIPort != 0 {
		openApiPort = f.OpenAPIPort
	}

	e := openapi.Start(openApiPort)
	defer e.Close()

	// Stop the service gracefully.
	if err := e.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}
