package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"runtime"

	_ "net/http/pprof"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pojol/gobot/factory"
	"github.com/pojol/gobot/mock"
	"github.com/pojol/gobot/server"
)

var (
	help bool

	dbmode       bool
	scriptPath   string
	openHttpMock bool
	openTcpMock  bool
)

const (
	// Version of gobot driver
	Version = "v0.2.1"

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

func initFlag() {
	flag.BoolVar(&help, "h", false, "this help")

	flag.BoolVar(&dbmode, "no_database", false, "Run in local mode")
	flag.BoolVar(&openHttpMock, "httpmock", false, "open http mock server")
	flag.BoolVar(&openTcpMock, "httpmock", false, "open tcp mock server")
	flag.StringVar(&scriptPath, "script_path", "script/", "Path to bot script")
}

func main() {

	defer func() {
		if err := recover(); err != nil {
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			fmt.Println("panic:", string(buf[:n]))
		}
	}()

	initFlag()
	flag.Parse()
	if help {
		flag.Usage()
		return
	}

	_, err := factory.Create(
		factory.WithNoDatabase(dbmode),
	)
	if err != nil {
		panic(err)
	}

	if openHttpMock {
		ms := mock.NewHttpServer()
		go ms.Start(":6666")
		defer ms.Close()
	}

	go func() {
		http.ListenAndServe(":6060", nil)
	}()

	fmt.Printf(banner, Version)

	// 查看有没有未完成的队列
	factory.Global.CheckTaskHistory()

	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())

	server.Route(e)
	e.Start(":8888")
	// Stop the service gracefully.
	if err := e.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}
