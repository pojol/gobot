package mock

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type guestRes struct {
	SessionID string
}

type accInfoReq struct {
	SessionID string
}

type accInfoRes struct {
	SessionID string
	Diamond   int32
	Gold      int32
}

type heroInfoReq struct {
	SessionID string
}

type heroInfoRes struct {
	SessionID string
	Heros     []heroInfo
}

type lvupReq struct {
	SessionID string
	HeroID    string
}

type lvupRes struct {
	SessionID string
	Heros     []heroInfo
}

type MockResponse struct {
	Code int32
	Msg  string
	Body interface{}
}

func routeGuest(ctx echo.Context) error {
	var res MockResponse
	acc := createAcc(uuid.New().String())

	res = MockResponse{
		Code: 200,
		Msg:  "",
		Body: guestRes{
			SessionID: acc.SessionID,
		},
	}

	ctx.JSON(http.StatusOK, res)
	return nil
}

func routeAccInfo(ctx echo.Context) error {
	accinfo := accInfoReq{}
	var res MockResponse
	var body accInfoRes
	var accPtr *MockAcc

	bts, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		res.Code = 400
		res.Msg = err.Error()
		goto ext
	}

	err = json.Unmarshal(bts, &accinfo)
	if err != nil {
		res.Code = 400
		res.Msg = err.Error()
		goto ext
	}

	accPtr, err = getAccInfo(accinfo.SessionID)
	if err != nil {
		res.Msg = err.Error()
		res.Code = 400
		goto ext
	}

	body = accInfoRes{
		Diamond: accPtr.Diamond,
		Gold:    accPtr.Gold,
	}
	res.Body = body
	res.Code = 200

ext:
	ctx.JSON(http.StatusOK, res)
	return nil
}

func routeHeroInfo(ctx echo.Context) error {
	var heroinfo heroInfoReq
	var res MockResponse
	var accPtr *MockAcc

	bts, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		res.Code = 400
		res.Msg = err.Error()
		goto ext
	}

	err = json.Unmarshal(bts, &heroinfo)
	if err != nil {
		res.Code = 400
		res.Msg = err.Error()
		goto ext
	}

	accPtr, err = getAccInfo(heroinfo.SessionID)
	if err != nil {
		res.Msg = err.Error()
		res.Code = 400
		goto ext
	}

	res.Body = heroInfoRes{
		Heros: accPtr.Heros,
	}
	res.Code = 200

ext:
	ctx.JSON(http.StatusOK, res)
	return nil
}

func routeHeroLvup(ctx echo.Context) error {
	lvup := lvupReq{}
	var res MockResponse
	var accPtr *MockAcc

	bts, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		res.Code = 400
		res.Msg = err.Error()
		goto ext
	}

	err = json.Unmarshal(bts, &lvup)
	if err != nil {
		res.Code = 400
		res.Msg = err.Error()
		goto ext
	}

	accPtr, err = setHeroLv(lvup.SessionID, lvup.HeroID)
	if err != nil {
		res.Code = 400
		res.Msg = "can't find hero sessionid " + lvup.SessionID + " hero " + lvup.HeroID
		goto ext
	}

	res.Code = 200
	res.Body = lvupRes{
		Heros: accPtr.Heros,
	}

ext:
	ctx.JSON(http.StatusOK, res)
	return nil
}
