package mock

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

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

var accmap map[string]*MockAcc

var strlst = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k"}

func createAcc() *MockAcc {

	token := ""
	for i := 0; i < 6; i++ {
		token += strlst[rand.Intn(len(strlst))]
	}

	acc := &MockAcc{
		Token:   token,
		Diamond: 50,
		Gold:    100,
		Heros: []heroInfo{
			{ID: "joy", Lv: 1},
			{ID: "pojoy", Lv: 2},
		},
	}

	accmap[token] = acc

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

	if _, ok := accmap[accinfo.Token]; ok {
		body = accInfoRes{
			Token:   accinfo.Token,
			Diamond: accmap[accinfo.Token].Diamond,
			Gold:    accmap[accinfo.Token].Gold,
		}
		res.Body = body
		res.Code = 200
	} else {
		res.Code = 400
		res.Msg = "can't find acc" + accinfo.Token
	}

ext:
	ctx.JSON(http.StatusOK, res)
	return nil
}

func routeHeroInfo(ctx echo.Context) error {
	var heroinfo heroInfoReq
	var res MockResponse

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

	if _, ok := accmap[heroinfo.Token]; ok {
		res.Body = heroInfoRes{
			Token: heroinfo.Token,
			Heros: accmap[heroinfo.Token].Heros,
		}
		res.Code = 200
	} else {
		res.Code = 400
		res.Msg = "can't find acc : " + heroinfo.Token
	}
ext:
	ctx.JSON(http.StatusOK, res)
	return nil
}

func routeHeroLvup(ctx echo.Context) error {
	lvup := lvupReq{}
	var res MockResponse

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

	if _, ok := accmap[lvup.Token]; ok {

		flag := false

		for k := range accmap[lvup.Token].Heros {
			if accmap[lvup.Token].Heros[k].ID == lvup.HeroID {
				accmap[lvup.Token].Heros[k].Lv++
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
			Heros: accmap[lvup.Token].Heros,
		}

	} else {
		res.Code = 400
		res.Msg = "can't find acc : " + lvup.Token
	}

ext:
	ctx.JSON(http.StatusOK, res)
	return nil
}

func NewServer() *echo.Echo {

	accmap = make(map[string]*MockAcc)
	rand.Seed(time.Now().UnixNano())

	mock := echo.New()
	mock.POST("/login/guest", routeGuest)
	mock.POST("/base/acc.info", routeAccInfo)
	mock.POST("/base/hero.info", routeHeroInfo)
	mock.POST("/base/hero.lvup", routeHeroLvup)

	return mock
}
