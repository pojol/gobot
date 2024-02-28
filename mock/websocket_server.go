package mock

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{}
)

func routes(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	go func() {

		for {
			// Read
			id, msg, err := ws.ReadMessage()
			if err != nil {
				break
			}

			if id == LoginGuest {
				acc := createAcc(uuid.NewString())
				setSession(acc.SessionID, ws)
			} else if id == Hello {
				err = wsHelloHandle(ws)
			} else if id == HeroInfo {
				err = wsHeroInfoHandle(ws, msg)
			} else if id == HeroLvup {
				err = wsHeroLvupHandle(ws, msg)
			}

			if err != nil {
				break
			}

		}
	}()

	return nil
}

func StartWebsocketServe(port string) {

	s := echo.New()
	s.HideBanner = true

	s.GET("/ws", routes)

	s.Start(":" + port)
}
