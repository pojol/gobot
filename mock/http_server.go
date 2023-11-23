package mock

import (
	"math/rand"
	"time"

	"github.com/labstack/echo/v4"
)

func NewHttpServer() *echo.Echo {

	rand.Seed(time.Now().UnixNano())

	mock := echo.New()
	mock.HideBanner = true
	mock.POST("/login/guest", routeGuest)
	mock.POST("/base/acc.info", routeAccInfo)
	mock.POST("/base/hero.info", routeHeroInfo)
	mock.POST("/base/hero.lvup", routeHeroLvup)

	return mock
}
