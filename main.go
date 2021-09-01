package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pojol/apibot/bot"
	"github.com/pojol/apibot/plugins"
)

type Response struct {
	Code int
	Msg  string
	Body interface{}
}

func Upload(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	res := &Response{}
	code := 200
	msg := ""

	bts, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		code = -1 // tmp
		msg = err.Error()
		goto EXT
	}

	if !json.Valid(bts) {
		code = -2
		msg = "json invalid"
		goto EXT
	}

	err = plugins.Load("./json.so")
	if err != nil {
		msg = err.Error()
		code = -3
		goto EXT
	}

	heros = heros[:0]
	heros = append(heros, heroInfo{ID: "pojol", Lv: 1})
	heros = append(heros, heroInfo{ID: "joy", Lv: 1})

	behaviorBuffer.Reset()
	behaviorBuffer.Write(bts)

	mbot, _ = bot.NewWithBehaviorFile(behaviorBuffer.Bytes(), srv.URL)

EXT:
	res.Code = code
	res.Msg = msg

	ctx.JSON(http.StatusOK, res)
	return nil
}

type RunResponse struct {
	Blackboard string
}

func Run(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	res := &Response{}
	body := &RunResponse{}
	code := 200
	msg := ""

	var err error

	if behaviorBuffer.Len() == 0 {
		code = -1
		msg = "not behavior data, need upload!"
		goto EXT
	}

	mbot.Run()

	body.Blackboard, err = mbot.GetMetadata()
	if err != nil {
		msg = err.Error()
		code = -3
	}

EXT:
	res.Code = code
	res.Msg = msg
	res.Body = body

	ctx.JSON(http.StatusOK, res)
	return nil
}

type StepResponse struct {
	Prev       string
	Cur        string
	Blackboard string
}

func Step(ctx echo.Context) error {
	ctx.Response().Header().Set("Access-Control-Allow-Origin", "*")
	res := &Response{}
	body := &StepResponse{}
	code := 200
	msg := ""

	var err error
	var s bot.State

	if behaviorBuffer.Len() == 0 {
		code = -1
		msg = "not behavior data, need upload!"
		goto EXT
	}

	s = mbot.RunStep()

	body.Blackboard, err = mbot.GetMetadata()
	if err != nil {
		msg = err.Error()
		code = -3
	}
	body.Cur = mbot.GetCurNodeID()
	body.Prev = mbot.GetPrevNodeID()

	if s == bot.SEnd {
		mbot = nil
		behaviorBuffer.Reset()
		code = -4
		msg = "执行到末尾，重新上传行为文件开始！"
		goto EXT
	} else if s == bot.SBreak {
		body.Blackboard, _ = mbot.GetMetadata()
		code = -5
		msg = "执行遇到请求错误!"
		mbot = nil
		behaviorBuffer.Reset()
		goto EXT
	}

EXT:
	res.Code = code
	res.Msg = msg
	res.Body = body

	ctx.JSON(http.StatusOK, res)
	return nil
}

type guestRes struct {
	Token string
}

type accInfoReq struct {
	Token string
}

type accInfoRes struct {
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

var heros []heroInfo

type heroInfoRes struct {
	Heros []heroInfo
}

type lvupReq struct {
	Token  string
	HeroID string
}

type lvupRes struct {
	Heros []heroInfo
}

var srv *httptest.Server
var mbot *bot.Bot
var behaviorBuffer bytes.Buffer

const ACCTOKEN = "abcd"

func main() {

	behaviorBuffer.Reset()

	// mock server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		reqbyt, _ := ioutil.ReadAll(req.Body)
		var byt []byte

		fmt.Println("http server recv ", req.RequestURI)

		if req.RequestURI == "/login/guest" {
			byt, _ = json.Marshal(guestRes{
				Token: ACCTOKEN,
			})
		} else if req.RequestURI == "/base/acc.info" {
			accinfo := accInfoReq{}
			err := json.Unmarshal(reqbyt, &accinfo)
			if err == nil && accinfo.Token == ACCTOKEN {
				byt, _ = json.Marshal(accInfoRes{
					Diamond: 100,
					Gold:    100,
				})
			} else {
				byt = []byte("error")
			}

		} else if req.RequestURI == "/base/hero.info" {

			heroinfo := heroInfoReq{}
			err := json.Unmarshal(reqbyt, &heroinfo)
			if err == nil && heroinfo.Token == ACCTOKEN {
				byt, _ = json.Marshal(heroInfoRes{
					Heros: heros,
				})
			} else {
				byt = []byte("error")
			}

		} else if req.RequestURI == "/base/hero.lvup" {
			lvup := lvupReq{}
			err := json.Unmarshal(reqbyt, &lvup)
			if err == nil && lvup.Token == ACCTOKEN {
				for k := range heros {
					if heros[k].ID == lvup.HeroID {
						heros[k].Lv++
						byt, _ = json.Marshal(lvupRes{
							Heros: heros,
						})
						break
					}
				}
			} else {
				byt = []byte("error")
			}
		}

		w.Write(byt)
	}))
	defer srv.Close()

	e := echo.New()
	e.Use(middleware.CORS())

	e.POST("/upload", Upload)
	e.POST("/run", Run)
	e.POST("/step", Step)

	e.Start(":8888")

	// Stop the service gracefully.
	if err := e.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}
