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

	if behaviorBuffer.Len() == 0 {
		code = -1
		msg = "not behavior data, need upload!"
		goto EXT
	}

	if !mbot.RunStep() {
		mbot = nil
		behaviorBuffer.Reset()
		code = -4
		msg = "执行到末尾，重新上传行为文件开始！"
		goto EXT
	}
	body.Blackboard, err = mbot.GetMetadata()
	if err != nil {
		msg = err.Error()
		code = -3
	}
	body.Cur = mbot.GetCurNodeID()
	body.Prev = mbot.GetPrevNodeID()

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

type infoRes struct {
	Diamond int32
	Gold    int32
}

var srv *httptest.Server
var mbot *bot.Bot
var behaviorBuffer bytes.Buffer

func main() {

	behaviorBuffer.Reset()

	// mock server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		ioutil.ReadAll(req.Body)
		var byt []byte

		fmt.Println("http server recv ", req.RequestURI)

		if req.RequestURI == "/login/guest" {
			byt, _ = json.Marshal(guestRes{
				Token: "abcd",
			})
		} else if req.RequestURI == "/base/acc.info" {
			byt, _ = json.Marshal(infoRes{
				Diamond: 100,
				Gold:    100,
			})
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
