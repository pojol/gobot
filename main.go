package main

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pojol/gobot-driver/factory"
	"github.com/pojol/gobot-driver/mock"
	"github.com/pojol/gobot-driver/server"
)

func main() {

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
