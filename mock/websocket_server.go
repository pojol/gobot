package mock

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	upgrader  = websocket.Upgrader{}
	byteOrder = "LittleEndian"
)

func getByteOrder() binary.ByteOrder {
	if byteOrder == "BigEndian" {
		return binary.BigEndian
	}
	return binary.LittleEndian
}

func routes(c echo.Context) error {

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			break
		}

		br := bytes.NewReader(msg)

		var msgId uint16
		binary.Read(br, getByteOrder(), &msgId)

		if msgId == LoginGuest {
			err = wsGuestHandle(ws)
		} else if msgId == Hello {
			err = wsHelloHandle(ws)
		} else if msgId == HeroInfo {
			err = wsHeroInfoHandle(ws, msg[2:])
		} else if msgId == HeroLvup {
			err = wsHeroLvupHandle(ws, msg[2:])
		}

		if err != nil {
			fmt.Println("recv msg", msgId, "err", err.Error())
			break
		}

	}

	return nil
}

func StartWebsocketServe(btorder string) *echo.Echo {

	byteOrder = btorder

	s := echo.New()
	s.HideBanner = true
	s.Use(middleware.Recover())

	s.GET("/ws", routes)
	//s.Static("/", "../public")

	return s
}
