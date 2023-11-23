package mock

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type guestRes struct {
	Token string
}

type accInfoReq struct {
	Token string
}

type accInfoRes struct {
	Token   string
	Diamond int32
	Gold    int32
}

type heroInfoReq struct {
	Token string
}

type heroInfoRes struct {
	Token string
	Heros []heroInfo
}

type lvupReq struct {
	Token  string
	HeroID string
}

type lvupRes struct {
	Token string
	Heros []heroInfo
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
			Token: acc.Token,
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

	accPtr, err = getAccInfo(accinfo.Token)
	if err != nil {
		res.Msg = err.Error()
		res.Code = 400
		goto ext
	}

	body = accInfoRes{
		Token:   accinfo.Token,
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

	accPtr, err = getAccInfo(heroinfo.Token)
	if err != nil {
		res.Msg = err.Error()
		res.Code = 400
		goto ext
	}

	res.Body = heroInfoRes{
		Token: accPtr.Token,
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

	accPtr, err = setHeroLv(lvup.Token, lvup.HeroID)
	if err != nil {
		res.Code = 400
		res.Msg = "can't find hero token " + lvup.Token + " hero " + lvup.HeroID
		goto ext
	}

	res.Code = 200
	res.Body = lvupRes{
		Token: lvup.Token,
		Heros: accPtr.Heros,
	}

ext:
	ctx.JSON(http.StatusOK, res)
	return nil
}
