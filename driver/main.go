package main

import (
	"context"
	"flag"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pojol/gobot/driver/factory"
	"github.com/pojol/gobot/driver/mock"
	"github.com/pojol/gobot/driver/server"
)

var (
	help bool

	dbHost string
	dbPwd  string
	dbName string
	dbUser string
)

func initFlag() {
	flag.BoolVar(&help, "h", false, "this help")

	flag.StringVar(&dbHost, "db_host", "127.0.0.1:3306", "set consul address")
	flag.StringVar(&dbPwd, "db_pwd", "gobot", "set mysql password")
	flag.StringVar(&dbName, "db_name", "gobot", "set mysql database name")
	flag.StringVar(&dbUser, "db_user", "gobot", "set mysql user")

}

func main() {
	initFlag()
	flag.Parse()

	_, err := factory.Create()
	if err != nil {
		panic(err)
	}

	ms := mock.NewServer()
	go ms.Start(":6666")

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		Skipper:           middleware.DefaultSkipper,
		StackSize:         4 << 10, // 4 KB
		DisableStackAll:   true,
		DisablePrintStack: true,
		LogLevel:          1,
	}))
	server.Route(e)
	e.Start(":8888")

	// Stop the service gracefully.
	if err := e.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}
