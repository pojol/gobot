package mock

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"sync"
	"time"

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

type heroInfo struct {
	ID string
	Lv int32
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

type MockAcc struct {
	Token   string
	Heros   []heroInfo
	Diamond int32
	Gold    int32
}

type MockResponse struct {
	Code int32
	Msg  string
	Body interface{}
}

var accmap = sync.Map{}

func createAcc() *MockAcc {

	acc := &MockAcc{
		Token:   uuid.New().String(),
		Diamond: 50,
		Gold:    100,
		Heros: []heroInfo{
			{ID: "joy", Lv: 1},
			{ID: "pojoy", Lv: 2},
		},
	}

	accmap.Store(acc.Token, acc)
	return acc
}

func routeGuest(ctx echo.Context) error {
	var res MockResponse
	acc := createAcc()

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
	var mapval interface{}
	var ok bool

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

	mapval, ok = accmap.Load(accinfo.Token)
	if !ok {
		res.Code = 400
		res.Msg = "can't find acc" + accinfo.Token
		goto ext
	}

	accPtr = mapval.(*MockAcc)
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
	var mapval interface{}
	var ok bool

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

	mapval, ok = accmap.Load(heroinfo.Token)
	if !ok {
		res.Code = 400
		res.Msg = "can't find acc" + heroinfo.Token
		goto ext
	}

	accPtr = mapval.(*MockAcc)
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
	var mapval interface{}
	var ok bool
	flag := false

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

	mapval, ok = accmap.Load(lvup.Token)
	if !ok {
		res.Code = 400
		res.Msg = "can't find acc" + lvup.Token
		goto ext
	}
	accPtr = mapval.(*MockAcc)

	for k := range accPtr.Heros {
		if accPtr.Heros[k].ID == lvup.HeroID {
			accPtr.Heros[k].Lv++
			flag = true
			break
		}
	}

	if !flag {
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
