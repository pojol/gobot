package main

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pojol/apibot/factory"
	"github.com/pojol/apibot/plugins"
	"github.com/pojol/apibot/server"
)

func main() {
	err := plugins.Load("./json.so")
	if err != nil {
		panic(err)
	}

	_, err = factory.Create(factory.WithMock())
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Use(middleware.CORS())

	server.Route(e)

	e.Start(":8888")

	// Stop the service gracefully.
	if err := e.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}
